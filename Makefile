LOCAL_BIN=$(CURDIR)/bin

clean-deps:
	rm -rf ./bin

# ----------------------------------------
# Code, Docs and Vendor
# ----------------------------------------

source_file=find . -type f \( -name "*.go" ! -regex ".*/vendor.*" \)

.PHONY: deps
deps:
	GOPRIVATE=gitlab.com go mod tidy
	GOPRIVATE=gitlab.com go mod download
	GOPRIVATE=gitlab.com go mod vendor

# ----------------------------------------
# [optional] Binary Dependencies
# ----------------------------------------

bin-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/go-swagger/go-swagger/cmd/swagger@latest

migrate-up:
	migrate -path ./migrations -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./migrations -database 'postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable' down

# ----------------------------------------
# [optional] Documentation
# ----------------------------------------

genswagger:
	$(LOCAL_BIN)/swagger generate spec --output=docs/openapi.yml
