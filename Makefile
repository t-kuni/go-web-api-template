build: generate
	go build -o cmd/app-server/main cmd/app-server/main.go

clean:
	git ls-files --others --ignored --exclude-standard ./ | \
	grep -v '.env*' | \
	grep -v '.idea*' | \
	grep -v 'build-errors.log' | \
	xargs -I {} rm -f {}

generate: clean
	go generate ./...

test: generate
	gotestsum --hide-summary=skipped -- -tags feature ./... -v

coverage: generate
	go test -tags feature -coverpkg=./... -coverprofile=coverage/coverage.o ./... > /dev/null
	go tool cover -html=coverage/coverage.o -o coverage/coverage.html