
.PHONY: node-daemon
node-daemon: generate
	KO_DOCKER_REPO=ko.local ko build -B ./cmd/node-daemon/

.PHONY: containerdbg-entrypoint
containerdbg-entrypoint: generate
	KO_DOCKER_REPO=ko.local ko build -B ./cmd/entrypoint/

.PHONY: test-binary
test-binary: generate
	KO_DOCKER_REPO=ko.local ko build -B ./cmd/test-binary/

.PHONY: test-openfile
test-openfile: generate
	KO_DOCKER_REPO=ko.local ko build -B ./test/test-images/test-openfile/
