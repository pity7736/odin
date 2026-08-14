.PHONY: lint test mocks coverage coverage-check check install-hooks

lint:
	golangci-lint run ./...

test:
	go test ./...

mocks:
	go run github.com/vektra/mockery/v3

COVERAGE_THRESHOLD := 92

coverage:
	go test -coverprofile=coverage.out -coverpkg=./src/... ./...
	go tool cover -func=coverage.out | tail -1

coverage-check: coverage
	@total=$$(go tool cover -func=coverage.out | tail -1 | awk '{print $$NF}' | tr -d '%'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "Coverage $$total% is below threshold $(COVERAGE_THRESHOLD)%"; \
		exit 1; \
	else \
		echo "Coverage $$total% meets threshold $(COVERAGE_THRESHOLD)%"; \
	fi

check: lint test coverage-check

install-hooks:
	@ln -sf $(shell pwd)/scripts/hooks/pre-commit .git/hooks/pre-commit
	@echo "hooks installed"
