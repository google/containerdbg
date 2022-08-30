.PHONY: ALL
all: images

include Makefile.btf.mk
include Makefile.proto.mk
include Makefile.go.mk
include Makefile.image.mk

test-image:
	docker build . -t containerdbg-test -f test/image/Dockerfile

test-in-container:
	export TEMPDIR=$(shell mktemp -d) && \
	cp -r $(PWD)/. $${TEMPDIR} && \
	docker run --privileged -e DOCKER_IN_DOCKER_ENABLED=true \
		-v ${HOME}/go/pkg:/go/pkg \
		-v /tmp/docker-graph:/docker-graph -v $${TEMPDIR}:/build -w /build \
		-v /sys/kernel/debug/:/sys/kernel/debug/ \
		-v /sys/fs/bpf/:/sys/fs/bpf/ \
		-v ${PWD}/artifacts:/artifacts \
		eu.gcr.io/modernize-prow/containerdbg-test:latest \
		./test/image/runner.sh make test

#--mount type=bind,source=/sys/fs/bpf/,target=/sys/fs/bpf/,bind-propagation=shared \

.PHONY: images
images: node-daemon containerdbg-entrypoint

.PHONY: test-images
test-images: test-binary test-openfile

.PHONY: prepare-kind-cluster
prepare-kind-cluster: images test-images
	-kind load docker-image ko.local/entrypoint
	-kind load docker-image ko.local/node-daemon
	-kind load docker-image ko.local/test-binary
	-kind load docker-image ko.local/test-openfile

.PHONY: example
example: prepare-kind-cluster
	kubectl apply -f ./deploy/node-daemon.yaml
	kubectl apply -f ./examples/modified_pod.yaml

.PHONY: test
test: generate build-go-linux images test-images
	go test ./...
