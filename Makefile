MOCKERY_VERSION := v2.53.3
BIN_DIR := $(HOME)/go/bin

.PHONY: mockery-install mocks

mockery-install:
	go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)

mocks:
	$(BIN_DIR)/mockery --all --config .mockery.yml