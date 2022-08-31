.PHONY: ALL
all: images build-go-linux

TARGET_REPO ?= ko.local
export TARGET_REPO

IMAGE_PULL_POLICY ?= IfNotPresent


ifeq ($(TARGET_REPO),ko.local)
IMAGE_PULL_POLICY = Never
endif

include Makefile.btf.mk
include Makefile.proto.mk
include Makefile.deploy.mk
include Makefile.go.mk
include Makefile.image.mk

export TEST_FLAGS

test-image:
	docker build . -t containerdbg-test -f test/image/Dockerfile

test-in-container:
	export TEMPDIR=$(shell mktemp -d) && \
	echo "$$TEMPDIR" && \
	cp -r $(PWD)/. $${TEMPDIR} && \
	docker run --privileged -e DOCKER_IN_DOCKER_ENABLED=true \
		-v ${HOME}/go/pkg:/go/pkg \
		-v /tmp/docker-graph:/docker-graph -v $${TEMPDIR}:/build -w /build \
		-v /sys/kernel/debug/:/sys/kernel/debug/ \
		-v /sys/fs/bpf/:/sys/fs/bpf/ \
		-v ${PWD}/artifacts:/artifacts \
		-e TEST_FLAGS=$(TEST_FLAGS) \
		-e COVERPROFILE=/artifacts/cover.prof \
		eu.gcr.io/modernize-prow/containerdbg-test:latest \
		./test/image/runner.sh make test

.PHONY: images
images: node-daemon containerdbg-entrypoint dnsproxy

.PHONY: test-tomcat
test-tomcat:
	#docker pull eu.gcr.io/modernize-prow/tomcat-petclinic

.PHONY: test-images
test-images: test-binary test-openfile test-tomcat

.PHONY: prepare-kind-cluster
prepare-kind-cluster: images test-images
	-kind load docker-image ko.local/entrypoint:latest
	-kind load docker-image ko.local/dnsproxy:latest
	-kind load docker-image ko.local/node-daemon:latest
	-kind load docker-image ko.local/test-binary:latest
	-kind load docker-image ko.local/test-openfile:latest

.PHONY: example
example: prepare-kind-cluster
	kubectl apply -f ./deploy/node-daemon.yaml
	kubectl apply -f ./examples/modified_pod.yaml

.PHONY: test
test: generate build-go-linux images test-images
	go test -coverprofile $(COVERPROFILE) -coverpkg ./... ./... $(TEST_FLAGS)

pre: deployment
