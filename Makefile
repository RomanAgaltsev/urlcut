.PHONY: gen-mock
gen-mock:	# Generate mocks
	mockgen -destination=internal/mocks/mock_repository.go -package=mocks github.com/RomanAgaltsev/urlcut/internal/interfaces Repository

.PHONY: buf-dep-upd
buf-dep-upd:	
	buf dep update

.PHONY: buf-lint
buf-lint:
	buf lint

.PHONY: buf-gen
buf-gen:
	buf generate

.PHONY: lint
lint:
	golangci-lint run ./...
	buf lint

.PHONY: tidy
tidy:	# Cleanup go.mod
	go mod tidy

.PHONY: fmt
fmt:
	gofumpt -l -w .
	goimports -local "github.com/RomanAgaltsev/urlcut" -w .

.PHONY: test
test:	# Execute the unit tests
	go test -count=1 -race -v -timeout 30s -coverprofile cover.out ./...

.PHONY: cover
cover:	# Show the cover report
	go tool cover -html cover.out

.PHONY: update
update:	# Update dependencies as recorded in the go.mod and go.sum files
	go list -m -u all
	go get -u ./...
	go mod tidy

.PHONY: doc
doc:	# godoc
	godoc -http=:6060

.PHONY: run
run:	# Run the application
	./cmd/shortener/shortener -d="postgres://postgres:12345@localhost:5432/praktikum?sslmode=disable" -k=secret