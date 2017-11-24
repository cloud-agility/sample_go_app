NAME		= go-sample

PRODUCTION = production

ifndef TAGS
	TAGS	= local
else
	TAGS :=$(subst /,-,$(TAGS))
endif
DOCKER_IMAGE	= $(NAME):$(TAGS)

ifndef REGISTRY
# use minikube by default
	REGISTRY	= 192.168.99.100:32767/default
	REGISTRY_SECRET = $(shell kubectl get secret | grep default-token | awk '{print $$1}')
endif

ifndef RELEASE
	NAMESPACE = staging-$(TAGS)
	NAMESPACE :=$(subst .,-,$(NAMESPACE))
else
	NAMESPACE = production
	REGISTRY = cloudagility
endif

# COMMAND DEFINITIONS
BUILD		= docker build -t
TEST		= docker run --rm
TEST_CMD	= go test
TEST_DIR	= /go/src/sample_go_app
VOLUME		= -v$(CURDIR):/$(TEST_DIR)
DEPLOY		= helm
LOGIN		= docker login
PUSH		= docker push
TAG		= docker tag

.PHONY: all
all: build unittest

.PHONY: build
build: Dockerfile
	echo ">> building app as $(DOCKER_IMAGE)"
	$(BUILD) $(DOCKER_IMAGE) .
	echo ">> packaging the $(DEPLOY) charts"
	$(DEPLOY) lint $(NAME)-chart
	$(DEPLOY) package $(NAME)-chart

.PHONY: unittest
unittest:
	echo ">> running tests on $(DOCKER_IMAGE)"
	$(TEST) $(VOLUME) $(DOCKER_IMAGE) $(TEST_CMD)

.PHONY: push
push:
ifeq ($(TAGS),$(RELEASE))
	echo ">> pushing release $(RELEASE) image to docker hub as $(REGISTRY)/$(DOCKER_IMAGE)"
	$(LOGIN) -u="$(DOCKER_USERNAME)" -p="$(DOCKER_PASSWORD)"
else
	echo ">> using $(REGISTRY) registry"
endif
	$(TAG) $(DOCKER_IMAGE) $(REGISTRY)/$(DOCKER_IMAGE)
	$(PUSH) $(REGISTRY)/$(DOCKER_IMAGE)

.PHONY: namespace
namespace:
	echo "Creating namespace $(NAMESPACE) if it doesn't already exist"
	$(DEPLOY) upgrade $(NAMESPACE) namespace-chart --install
ifneq ($(NAMESPACE),$(PRODUCTION))
	echo "Patching registry secret $(REGISTRY_SECRET) for namespace $(NAMESPACE)"
	kubectl patch sa default -p '{"imagePullSecrets": [{"name": "$(REGISTRY_SECRET)"}]}' --namespace $(NAMESPACE)
endif

.PHONY: deploy
deploy: push namespace
	echo ">> Use $(DEPLOY) to install $(NAME)-chart"
	## Override the values.yaml with the target
	$(DEPLOY) upgrade $(NAME)-$(NAMESPACE) $(NAME)-chart --install --set image.repository=$(REGISTRY),image.name=$(NAME) --namespace $(NAMESPACE)  --wait

.PHONY: cleankube
cleankube:
	echo ">> cleaning kube cluster for namespace $(NAMESPACE)"
	$(DEPLOY) delete $(NAME) --purge
	kubectl delete namespace $(NAMESPACE)
