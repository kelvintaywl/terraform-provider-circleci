default: testacc

TF_STACK_DIR := ./sandbox
OS ?= darwin
ARCH ?= arm64
OS_ARCH := $(OS)_$(ARCH)

.PHONY: docs
docs:
	go generate ./...

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Builds the go binary
.PHONY: binary
binary:
	go fmt ./...
	echo "Building Go binary"
	go build -o terraform-provider-circleci_v0.0.1

# Sets up your local workstation to "accept" this local provider binary
.PHONY: init
init: binary
	echo "Initializing..."
	echo "Setting up for local provider..."
	# assuming your workstation is on Mac M1
	mkdir -p ~/.terraform.d/plugins/example.com/kelvintaywl/circleci/0.0.1/$(OS_ARCH)
	ln -s $(CURDIR)/terraform-provider-circleci_v0.0.1 ~/.terraform.d/plugins/example.com/kelvintaywl/circleci/0.0.1/$(OS_ARCH)/terraform-provider-circleci_v0.0.1

# Builds the go binary, and cleans up Terraform lock file just in case
.PHONY: build
build: binary
	if [ -f "sandbox/.terraform.lock.hcl" ]; then \
	  rm sandbox/.terraform.lock.hcl; \
	fi

.PHONY: tf.init
tf.init:
	terraform -chdir=$(TF_STACK_DIR) fmt
	terraform -chdir=$(TF_STACK_DIR) init
	terraform -chdir=$(TF_STACK_DIR) validate

# Runs all critical Terraform commands, short of applying
.PHONY: tf.plan
tf.plan: tf.init
	terraform -chdir=$(TF_STACK_DIR) fmt
	terraform -chdir=$(TF_STACK_DIR) init
	terraform -chdir=$(TF_STACK_DIR) validate
	terraform -chdir=$(TF_STACK_DIR) plan

# Applies for sandbox
.PHONY: tf.apply
tf.apply: tf.plan
	terraform -chdir=$(TF_STACK_DIR) apply -auto-approve

# Deletes sandbox
.PHONY: tf.destroy
tf.destroy:
	terraform -chdir=$(TF_STACK_DIR) destroy -auto-approve
