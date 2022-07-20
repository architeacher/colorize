# 🐳 Docker configurations

## 🐳 Docker build flags.
DOCKER_BUILD_FLAGS ?= --force-rm --rm --no-cache --pull

## 🐳 Docker file path.
DOCKER_FILE ?= .dockerfile-${ARCH}

## Build script directory.
BUILD_SCRIPTS_DIR ?= ./build/scripts

## 🐳 Docker .env file preparation script path.
DOCKER_ENV_FILE_SCRIPT_PATH ?= $(BUILD_SCRIPTS_DIR)/prepare-env-file.sh

## 🐳 Docker file preparation script path.
DOCKER_FILE_SCRIPT_PATH ?= $(BUILD_SCRIPTS_DIR)/prepare-docker-file.sh

## 🐳 Prepare script path.
PREPARE_SCRIPT_PATH ?= $(BUILD_SCRIPTS_DIR)/prepare.sh

## 🐳 Deploy script path.
DEPLOY_SCRIPT_PATH ?= $(BUILD_SCRIPTS_DIR)/deploy.sh

## 🐳 Docker image build tag.
DOCKER_IMAGE_TAG ?= latest

## 🐳 Registry name that image artifacts should be produced for.
REGISTRY ?= $(shell whoami | sed -e 's|\.||g' 2> /dev/null)

## 🐳 Image name.
IMAGE_NAME ?= $(shell basename $(CURDIR) 2> /dev/null)

## 🐳 Registry repository.
REGISTRY_REPO ?= $(REGISTRY)/$(IMAGE_NAME)
