SHELL := /bin/bash

# ================================================================
# TOOLS

# expvarmon -ports=":3000" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"


# ================================================================
# GO

go-run:
	go run app/services/service-api/main.go | go run app/tooling/logfmt/main.go

go-build:
	go build -ldflags "-X main.build=local"

go-tidy:
	go mod tidy
	go mod vendor


# ================================================================
# DOCKER

IMAGE = service-amd64
VERSION := 1.0

docker-build: docker-build-service

docker-build-service:
	docker build \
	-f zarf/docker/service-api.dockerfile \
	-t $(IMAGE):$(VERSION) \
	--build-arg BUILD_REF=$(VERSION) \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	.

docker-run-service:
	docker run $(IMAGE):$(VERSION)

docker-sh-service:
	docker run -it $(IMAGE):$(VERSION) sh

# ================================================================
# KIND

KIND_CLUSTER = localhost-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.24.0@sha256:0866296e693efe1fed79d5e6c7af8df71fc73ae45e3679af05342239cdc5bc8e \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current=true --namespace=service-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	cd zarf/k8s/kind/service-pod; kustomize edit set image service-api-image=$(IMAGE):$(VERSION)
	kind load docker-image service-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build zarf/k8s/kind/service-pod | kubectl apply -f -

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-service:
	kubectl get pods -o wide --watch

kind-logs:
	kubectl logs -l app=service --all-containers=true -f --tail=100

kind-restart:
	kubectl rollout restart deployment service-pod

kind-update-restart: docker-build kind-load kind-restart

kind-update-apply: docker-build kind-load kind-apply

kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=service
