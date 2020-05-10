APP = imagenie

IMAGE = quay.io/otaviof/$(APP)
IMAGE_TAG = $(IMAGE):latest
IMAGE_DEV_TAG = $(IMAGE)-dev:latest

OUTPUT_DIR ?= _output
OUTPUT_BIN = $(OUTPUT_DIR)/$(APP)
GO_FLAGS ?= -v -mod=vendor
GO_TEST_FLAGS ?= -failfast

UNIT_TEST_TARGET ?= ./cmd/... ./pkg/...
E2E_TEST_TARGET ?= ./test/e2e/...

DEVCONTAINER_ARGS ?=
RUN_ARGS ?=
TEST_ARGS ?=

default: build

.PHONY: vendor
vendor:
	@go mod vendor

.PHONY: clean
clean:
	@rm -rf $(OUTPUT_DIR)

.PHONY: $(OUTPUT_BIN)
$(OUTPUT_BIN):
	@if [[ ! -d "$(OUTPUT_DIR)" ]] ; then mkdir -v $(OUTPUT_DIR) ; fi
	go build $(GO_FLAGS) -o $(OUTPUT_BIN) ./cmd/$(APP)/.

build: vendor $(OUTPUT_BIN)

run:
	go run $(GO_FLAGS) ./cmd/$(APP)/* $(RUN_ARGS)

test: test-unit test-e2e

.PHONY: test-unit
test-unit:
	go test $(GO_FLAGS) $(GO_TEST_FLAGS) $(TEST_ARGS) $(UNIT_TEST_TARGET)

.PHONY: test-e2e
test-e2e:
	echo "## TODO: write me! ##"

image:
	docker build --tag="$(IMAGE_TAG)" .

devcontainer-image:
	docker build --tag="$(IMAGE_DEV_TAG)" --file="Dockerfile.dev" .

devcontainer-run:
	docker run \
		--rm \
		--privileged \
		--volume="${PWD}:/src/$(APP)" \
		--workdir="/src/$(APP)" \
		$(IMAGE_DEV_TAG) $(DEVCONTAINER_ARGS)

devcontainer-exec:
	@docker exec --interactive --tty --workdir="/workspaces/$(APP)" $(APP) bash
