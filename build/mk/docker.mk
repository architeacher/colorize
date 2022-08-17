include build/mk/config/docker.mk

DOCKER := $(shell docker info > /dev/null 2>&1 || $(SUDO) 2>&1) docker
DOCKER_LOGO := "🐳"

images.amd64 = alpine
images.arm64 = alpine
images.arm = alpine
images.ppc64le = alpine

BASE_IMAGE ?= ${images.$(ARCH)}

## 🐳 Build image name.
DOCKER_BUILD_IMAGE ?= golang:$(GO_VERSION)-alpine
## 🐳 Build image tag.
DOCKER_IMAGE_TAG ?= $(BUILD_VERSION)
## 🐳 Target build image name.
DOCKER_TARGET_IMAGE ?= $(REGISTRY_REPO):$(DOCKER_IMAGE_TAG)
## 🐳 Container name that the action will be performed on.
DOCKER_CONTAINER_NAME ?= $(BINARY_PREFIX).${DOCKER_IMAGE_TAG}

.PHONY: docker
docker: build-dirs clean-containers clean-images clean-volumes deploy docker-clean docker-health image list-containers

.dockerfile-$(ARCH): .env
	$(call printMessage,"${DOCKER_LOGO} Preparing file",$(INFO_CLR))
	$(SUDO) BASE_IMAGE=${BASE_IMAGE} \
	    ARCH=${ARCH} \
	    PKG_NAMESPACE=${PKG_NAMESPACE} \
	    bash $(DOCKER_FILE_SCRIPT_PATH) $(args) 2>&1

.PHONY: .env
.env:
	$(call printMessage,"${DOCKER_LOGO} Preparing .env file",$(INFO_CLR))
	$(SUDO) ARCH=${ARCH} bash $(DOCKER_ENV_FILE_SCRIPT_PATH) $(args) 2>&1

.PHONY: build-dirs
build-dirs:
	$(call printMessage,"${DOCKER_LOGO} Building mapping directories for Go",$(INFO_CLR))
	mkdir -p ${GO_GENERATED_DIR}/bin ${GO_GENERATED_DIR}/pkg ${GO_GENERATED_DIR}/src/$(PKG_NAMESPACE) 2>&1

.PHONY: clean-containers
clean-containers: ## 🐳 to clean inactive containers data.
	$(call printMessage,"${DOCKER_LOGO} 🧹 Cleaning up containers data",$(WARN_CLR))
	$(DOCKER) container prune -f 2>&1
	$(eval EXITED_CONTAINERS := $(shell $(DOCKER) ps -aqf status=exited -f status=dead 2>&1))
	test -n "${EXITED_CONTAINERS}" && $(DOCKER) rm ${EXITED_CONTAINERS} || true 2>&1

.PHONY: clean-files
clean-files: ## 🐳 to clean deployment generated files.
	$(call printMessage,"${DOCKER_LOGO} 🧹 Cleaning up generated files and directories",$(WARN_CLR))
	rm -rf .docker-compose-*.yaml .dockerfile* .env 2>&1

.PHONY: clean-images
clean-images: ## 🐳 to clean inactive images data.
	$(call printMessage,"${DOCKER_LOGO} 🧹 Cleaning up images data",$(WARN_CLR))
	$(DOCKER) image prune -f 2>&1
	$(eval DANGLING_IMAGES:= $(shell $(DOCKER) images -aqf dangling=true 2>&1))
	test -n "${DANGLING_IMAGES}" && $(DOCKER) rmi ${DANGLING_IMAGES} || true 2>&1

.PHONY: clean-volumes
clean-volumes: ## 🐳 to clean inactive containers volumes.
	$(call printMessage,"${DOCKER_LOGO} 🧹 Cleaning up containers volumes",$(WARN_CLR))
	$(DOCKER) volume prune -f 2>&1
	$(eval DANGLING_VOLUMES := $(shell $(DOCKER) volume ls -qf dangling=true 2>&1))
	test -n "${DANGLING_VOLUMES}" && $(DOCKER) volume rm ${DANGLING_VOLUMES} || true 2>&1

.PHONY: deploy
deploy: build-dirs docker-prepare update-pkg-version ## 🐳 to deploy a docker container.
	$(call printMessage,"${DOCKER_LOGO} Deploying container",$(INFO_CLR))
	$(SUDO) ARCH=${ARCH} bash $(DEPLOY_SCRIPT_PATH) $(args) 2>&1

.PHONY: docker-clean
docker-clean: clean-images clean-containers clean-volumes ## 🐳 to clean inactive Docker data.

.PHONY: docker-exec
docker-exec: ## 🐳 to execute command inside the docker container.
	$(call printMessage,"${DOCKER_LOGO} Executing command inside the container",$(INFO_CLR))
	$(DOCKER) exec -it $(DOCKER_CONTAINER_NAME) $(CMD) 2>&1

.PHONY: docker-health
docker-health: ## 🐳 to get the health state docker container.
	$(call printMessage,"${DOCKER_LOGO} Getting health state of the container",$(INFO_CLR))
	$(DOCKER) inspect --format='{{json .State.Health}}' $(DOCKER_CONTAINER_NAME) 2>&1

.PHONY: docker-kill
docker-kill: ## 🐳 to send kill signal to the main process at the docker container.
	$(call printMessage,"${DOCKER_LOGO} 🥷️ Sending kill($(args)) signal to main process",$(INFO_CLR))
	$(MAKE) docker-exec CMD="pkill $(args) $(BINARY_PREFIX)" > /dev/null 2>&1

.PHONY: docker-logs
docker-logs: ## 🐳 to get logs from the docker container.
	$(call printMessage,"${DOCKER_LOGO} Following logs of the container",$(INFO_CLR))
	$(DOCKER) logs -f $(DOCKER_CONTAINER_NAME) 2>&1

.PHONY: docker-prepare
docker-prepare: ## 🐳 prepare docker files from the templates.
	$(call printMessage,"${DOCKER_LOGO} Preparing docker file",$(INFO_CLR))
	$(SUDO) BASE_IMAGE=${BASE_IMAGE} \
	    REGISTRY=${REGISTRY} \
	    IMAGE_NAME=${IMAGE_NAME} \
	    IMAGE_TAG=${DOCKER_IMAGE_TAG} \
	    REGISTRY_REPO=${REGISTRY_REPO} \
	    ARCH=${ARCH} \
	    SERVICE_NAME=${SERVICE_NAME} \
	    SERVICE_DESCRIPTION=${SERVICE_DESCRIPTION} \
	    GO_VERSION=${GO_VERSION} \
	    PKG_NAMESPACE=${PKG_NAMESPACE} \
	    bash $(PREPARE_SCRIPT_PATH) $(args) 2>&1

.PHONY: docker-shell
docker-shell: .env build-dirs ## 🐳 run shell command inside the docker container, Example: make docker-shell CMD="-c 'ls > .files'"
	$(call printMessage,"${DOCKER_LOGO} 🏃 Running an image \"$(DOCKER_BUILD_IMAGE)\" shell in the containerized build environment",$(INFO_CLR))
	$(DOCKER) run \
	    -it \
	    --rm \
	    -u $$(id -u):$$(id -g) \
	    -v $(CURDIR)/${GO_GENERATED_DIR}:/go \
	    -v $(CURDIR)/go/src/$(PKG_NAMESPACE) \
	    -v $(CURDIR)/${GO_GENERATED_DIR}/bin:/go/bin \
	    -v $(CURDIR)/${GO_GENERATED_DIR}/cache:/.cache \
	    -w /go/src/$(PKG_NAMESPACE) \
	    --env-file .env \
	    $(DOCKER_BUILD_IMAGE) \
	    /bin/sh $(CMD) 2>&1

.PHONY: image
image: .dockerfile-$(ARCH) ## 🐳 to build a docker image.
	$(call printMessage,"${DOCKER_LOGO} 🌄 Creating image \"${DOCKER_TARGET_IMAGE}\"",$(INFO_CLR))
	$(DOCKER) build \
	    ${DOCKER_BUILD_FLAGS} \
	    -t ${DOCKER_TARGET_IMAGE} \
	    -f $(DOCKER_FILE) \
	    $(args) . 2>&1

.PHONY: list-containers
list-containers: ## 🐳 to list all containers.
	$(call printMessage,"${DOCKER_LOGO} Listing containers",$(INFO_CLR))
	$(DOCKER) ps -a --format "table {{.ID}} {{.Names}} $(MSG_PRFX) {{json .Ports}}" ${args} 2>&1

.PHONY: list-images
list-images: ## 🐳 to list all images.
	$(call printMessage,"${DOCKER_LOGO} Listing images",$(INFO_CLR))
	$(DOCKER) images ${args} 2>&1

.PHONY: list-volumes
list-volumes: ## 🐳 to list all volumes.
	$(call printMessage,"${DOCKER_LOGO} Listing volumes",$(INFO_CLR))
	$(DOCKER) volume ls ${args} 2>&1

.PHONY: publish
publish: ## 🐳 to publish the docker image to dockerhub repository.
	$(call printMessage,"${DOCKER_LOGO} 🚀 Pushing image to \"${DOCKER_TARGET_IMAGE}\"",$(INFO_CLR))
	$(DOCKER) push ${DOCKER_TARGET_IMAGE} 2>&1
