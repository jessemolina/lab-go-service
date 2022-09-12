SHELL := /bin/bash

# ================================================================
# GO

go-run:
	go run app/services/sales-api/main.go

go-build:
	go build -ldflags "-X main.build=local"


# ================================================================
# DOCKER

VERSION := 1.0

docker-build: docker-build-service

docker-build-service:
	docker build \
	-f zarf/docker/sales-api.dockerfile \
	-t service-amd64:$(VERSION) \
	--build-arg BUILD_REF=$(VERSION) \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	.

docker-sh-service:
	docker run -it service-amd64:$(VERSION) sh


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
	kind load docker-image service-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	cat zarf/k8s/base/service-pod/base-service.yaml | kubectl apply -f -

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

kind-update: docker-build kind-load kind-restart

kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=sales
