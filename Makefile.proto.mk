#
# START protobuf list
#

node_daemon_proto = \
	proto/node_api.proto \
	$(NULL)

events_proto = \
	proto/event.proto \
	proto/summary.proto \
	$(NULL)

gen_protobuf_go = \
	$(node_daemon_proto) \
	$(events_proto) \
	$(NULL)

gen_grpc_go = \
	$(node_daemon_proto) \
	$(NULL)

#
# END protobuf list
#

proto_go_artifacts = $(gen_protobuf_go:.proto=.pb.go)
proto_go_grpc_artifacts = $(gen_grpc_go:.proto=_grpc.pb.go)
PROTOC := $(if $(PROTOC),$(PROTOC),$(shell which protoc))

TOOLS_DIR=$(abspath out/tools)
PROTOC_GEN_GO:=$(if $(PROTOC_GEN_GO),$(PROTOC_GEN_GO),$(TOOLS_DIR)/protoc-gen-go)
PROTOC_GEN_GO_GRPC:=$(if $(PROTOC_GEN_GO_GRPC),$(PROTOC_GEN_GO_GRPC),$(TOOLS_DIR)/protoc-gen-go-grpc)

$(PROTOC_GEN_GO):
	env "GOBIN=$(TOOLS_DIR)" go install -mod=mod google.golang.org/protobuf/cmd/protoc-gen-go

$(PROTOC_GEN_GO_GRPC):
	env "GOBIN=$(TOOLS_DIR)" go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0

$(proto_go_grpc_artifacts): %_grpc.pb.go: %.proto | $(PROTOC) $(PROTOC_GEN_GO_GRPC)
	$(PROTOC) \
			${PROTOBUF_INCLUDES} \
			-I $(dir $@) \
			--plugin "$(PROTOC_GEN_GO_GRPC)" \
			--go-grpc_out=$(dir $@) \
			--go-grpc_opt=paths=source_relative \
			--plugin=grpc \
			$< \

$(proto_go_artifacts): %.pb.go: %.proto | $(PROTOC) $(PROTOC_GEN_GO)
	$(PROTOC) \
			${PROTOBUF_INCLUDES} \
			-I $(dir $@) \
			--plugin "$(PROTOC_GEN_GO)" \
			--go_out=$(dir $@) \
			--go_opt=paths=source_relative \
			--plugin=go \
			$< \

all_proto_go = \
	$(proto_go_artifacts) \
	$(proto_go_grpc_artifacts) \
	$(NULL)
