
.PHONY: node-daemon
node-daemon: generate
	KO_DOCKER_REPO=${TARGET_REPO} ko build -B ./cmd/node-daemon/ -t ${TAG}

.PHONY: containerdbg-entrypoint
containerdbg-entrypoint: generate
	KO_DOCKER_REPO=${TARGET_REPO} ko build -B ./cmd/entrypoint/ -t ${TAG}

.PHONY: dnsproxy
dnsproxy: generate
	KO_DOCKER_REPO=${TARGET_REPO} ko build -B ./cmd/dnsproxy/ -t ${TAG}

.PHONY: test-binary
test-binary: generate
	KO_DOCKER_REPO=${TARGET_REPO} ko build -B ./cmd/test-binary/

.PHONY: test-openfile
test-openfile: generate
	KO_DOCKER_REPO=${TARGET_REPO} ko build -B ./test/test-images/test-openfile/

btf-install-image:
	docker build btf-install/ -t ${TARGET_REPO}/btf-installer:latest
	docker push ${TARGET_REPO}/btf-installer:latest
