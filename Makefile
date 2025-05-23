build: generate
	go build -o cmd/app-server/main cmd/app-server/main.go

clean:
	git ls-files --others --ignored --exclude-standard ./ | \
	grep -v '\.env*' | \
	grep -v '\.idea*' | \
	grep -v 'build-errors.log' | \
	xargs -I {} rm -f {}

generate: clean
	go generate ./...

test: generate
	go tool gotestsum --hide-summary=skipped -- ./... -v

coverage: generate
	go test -coverpkg=./... -coverprofile=coverage/coverage.o ./... > /dev/null
	go tool cover -html=coverage/coverage.o -o coverage/coverage.html