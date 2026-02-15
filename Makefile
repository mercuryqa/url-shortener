MOCKERY_VERSION := v2.53.3
GOFUMPT_VERSION := v0.8.0
GCI_VERSION := v0.13.4
GOLANGCI_LINT_VERSION := v2.1.5


BIN_DIR := $(HOME)/go/bin
GOFUMPT := $(BIN_DIR)/gofumpt
GCI := $(BIN_DIR)/gci
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint


.PHONY: mockery-install mocks
mockery-install:
	go install github.com/vektra/mockery/v2@$(MOCKERY_VERSION)

mocks:
	$(BIN_DIR)/mockery --all --config .mockery.yml

.PHONY: install-formatters
install-formatters:
	@mkdir -p $(BIN_DIR)
	@if [ ! -f $(GOFUMPT) ]; then \
		echo "Устанавливаем gofumpt $(GOFUMPT_VERSION)..."; \
		GOBIN=$(BIN_DIR) go install mvdan.cc/gofumpt@$(GOFUMPT_VERSION); \
	fi
	@if [ ! -f $(GCI) ]; then \
		echo "Устанавливаем gci $(GCI_VERSION)..."; \
		GOBIN=$(BIN_DIR) go install github.com/daixiang0/gci@$(GCI_VERSION); \
	fi

.PHONY: format
format: install-formatters
	@echo "Форматируем через gofumpt ..."
	@find . -type f -name '*.go' ! -path '*/mocks/*' -exec $(GOFUMPT) -extra -w {} +

	@echo "Сортируем импорты через gci ..."
	@find . -type f -name '*.go' ! -path '*/mocks/*' -exec $(GCI) write \
		-s standard -s default -s "prefix(github.com/mercuryqa/rocket-lab)" {} +


.PHONY: install-golangci-lint
install-golangci-lint:
	@if [ ! -x "$(GOLANGCI_LINT)" ]; then \
		mkdir -p "$(BIN_DIR)"; \
		echo "Устанавливаем golangci-lint $(GOLANGCI_LINT_VERSION)..."; \
		GOBIN="$(BIN_DIR)" go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION); \
	else \
		echo "golangci-lint уже установлен"; \
	fi

.PHONY: lint
lint: install-golangci-lint
	@echo "Запуск линтера..."; \
	$(GOLANGCI_LINT) run ./... --config=.golangci.yml