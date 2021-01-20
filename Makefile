#################################################################################
# BUILD COMMANDS
#################################################################################
build: 
	go build -o ./bin/bible ./cmd/bible
gen:
	go generate ./...

#################################################################################
# TEST COMMANDS
#################################################################################
test: 
	go test -cover ./internal/... 
test-coverage:
	go test -coverpkg ./internal/... -coverprofile coverage.out ./... && go tool cover -html=coverage.out
