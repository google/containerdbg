
deployment:
	rm -rf ./deploy/rendered
	kpt fn source deploy/ | kpt fn eval - --image gcr.io/kpt-fn/apply-setters:v0.2.0 -- repo=${TARGET_REPO} tag=latest policy=${IMAGE_PULL_POLICY} | kpt fn sink ./deploy/rendered

install-btf: btf-install-image
	kpt fn source btf-install/ | kpt fn eval - --image gcr.io/kpt-fn/apply-setters:v0.2.0 -- repo=${TARGET_REPO} tag=latest | kpt live apply -

