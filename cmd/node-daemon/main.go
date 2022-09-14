// Copyright 2022 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"syscall"

	"github.com/go-logr/logr"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"github.com/google/containerdbg/pkg/consts"
	"github.com/google/containerdbg/pkg/events"
	"github.com/google/containerdbg/pkg/events/sources"
	"github.com/google/containerdbg/pkg/logger"
	"github.com/google/containerdbg/proto"
)

const MaxWrittenEventsPerFile = 10000
const EventsFilePath = "/var/run/containerdbg/data/events.pb"

var _ credentials.TransportCredentials = &UnixTransportCredentials{}
var _ credentials.AuthInfo = &UnixAuthInfo{}

type UnixAuthInfo struct {
	creds *syscall.Ucred
}

func (UnixAuthInfo) AuthType() string {
	return "unix"
}

type UnixTransportCredentials struct {
	log logr.Logger
}

func (*UnixTransportCredentials) getCredentials(conn *net.UnixConn) (*syscall.Ucred, error) {
	f, err := conn.File()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return syscall.GetsockoptUcred(int(f.Fd()), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
}

func (transport *UnixTransportCredentials) ClientHandshake(ctx context.Context, n string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	return conn, nil, nil
}

func (transport *UnixTransportCredentials) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	unixConn, ok := conn.(*net.UnixConn)
	if !ok {
		return nil, nil, status.Error(codes.Internal, "Bad socket")
	}

	creds, err := transport.getCredentials(unixConn)
	if err != nil {
		return nil, nil, status.Errorf(codes.Unauthenticated, "bad auth info: %s", err)
	}
	transport.log.Info("Got connection with", "creds", creds)

	return conn, &UnixAuthInfo{creds: creds}, nil
}

func (*UnixTransportCredentials) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{}
}

func (transport *UnixTransportCredentials) Clone() credentials.TransportCredentials {
	return transport
}

func (*UnixTransportCredentials) OverrideServerName(string) error {
	return nil
}

type DaemonServer struct {
}

func eventsFileName(index int) string {
	return fmt.Sprintf("%s.%d", EventsFilePath, index)
}

func SwitchEventsFile(fIndex int) (*events.EventWriter, error) {
	eventsFile, err := os.OpenFile(EventsFilePath, os.O_EXCL|os.O_CREATE|os.O_WRONLY, os.FileMode(0700))
	if err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
		if err := os.Rename(EventsFilePath, eventsFileName(fIndex)); err != nil {
			return nil, err
		}
		eventsFile, err = os.OpenFile(EventsFilePath, os.O_EXCL|os.O_CREATE|os.O_WRONLY, os.FileMode(0700))
		if err != nil {
			return nil, err
		}
	}
	eventsWriter := events.NewEventWriter(eventsFile)

	return eventsWriter, nil
}

func EventsPersister(ctx context.Context, log logr.Logger, e <-chan *proto.Event) error {
	log.Info("opening persistence file", "path", EventsFilePath)

	fIndex, err := findLastFileIndex()
	if err != nil {
		return err
	}

	eventsWriter, err := SwitchEventsFile(fIndex)
	if err != nil {
		return err
	}
	defer eventsWriter.Close()
	writtenEvents := 0

	for {
		select {
		case msg := <-e:
			if err := eventsWriter.Write(msg); err != nil {
				log.Error(err, "failed to write event to file", "event", msg)
			}
			if writtenEvents >= MaxWrittenEventsPerFile {
				eventsWriter.Close()
				fIndex++
				eventsWriter, err = SwitchEventsFile(fIndex)
				if err != nil {
					return err
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func StartNodeDaemon(ctx context.Context, log logr.Logger, errgr *errgroup.Group) error {
	sharedDir := os.Getenv(consts.SharedDirectoryEnv)

	if sharedDir == "" {
		return fmt.Errorf("%s env was not provided", consts.SharedDirectoryEnv)
	}

	mgr := sources.NewEventSourceManager(log)

	dynamicSource := events.NewDynamicSource()
	mgr.AddEventSource(dynamicSource)

	if err := mgr.Load(); err != nil {
		return fmt.Errorf("event sources failed to load: %w", err)
	}
	defer mgr.Unload()
	log.Info("event sources have been loaded")

	log.Info("starting event persister")
	errgr.Go(func() error {
		return EventsPersister(ctx, log, mgr.Events())
	})

	opts := []grpc.ServerOption{
		grpc.Creds(&UnixTransportCredentials{log: log}),
		grpc.ChainUnaryInterceptor(
			logger.LoggingInterceptor(log),
		),
	}

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterNodeDaemonServiceServer(grpcServer, &NodeDaemonServiceServer{
		Manager:       mgr,
		DynamicSource: *dynamicSource,
	})

	address := filepath.Join(sharedDir, consts.NodeDaemonSocketName)
	os.Remove(address)
	listener, err := net.Listen("unix", address)
	if err != nil {
		return fmt.Errorf("failed to listend on address %s: %w", address, err)
	}
	if err := os.Chmod(address, os.FileMode(0777)); err != nil {
		return fmt.Errorf("failed to set permissions on the unix socket %s: %w", address, err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}

	return nil
}

func handleCollect(w http.ResponseWriter, req *http.Request) {
	files, err := getEventFilesList()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			if os.IsNotExist(err) {
				break
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		io.Copy(w, f)

		f.Close()
	}

	// Send last file
	f, err := os.Open(EventsFilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(w, f)
}

func StartCollectionServer(ctx context.Context, log logr.Logger) error {
	srv := &http.Server{Addr: ":8080"}
	http.HandleFunc("/collect", handleCollect)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error(err, "failed to start server")
		}
	}()

	<-ctx.Done()
	return srv.Shutdown(context.Background())
}

func saveTerminationMessage(log logr.Logger, message string) {
	// We ignore error as there is nothing we can do about it
	if err := os.WriteFile("/dev/termination-log", []byte(message), os.FileMode(0777)); err != nil {
		log.Error(err, "error occured while saving termination message")
	}
}

func main() {
	log := logger.NewHeadlessLogger()
	errgr, ctx := errgroup.WithContext(context.Background())
	errgr.Go(func() error {
		return StartNodeDaemon(ctx, log, errgr)
	})
	errgr.Go(func() error {
		return StartCollectionServer(ctx, log)
	})

	if err := errgr.Wait(); err != nil {
		saveTerminationMessage(log, err.Error())
		panic(err)
	}
}
