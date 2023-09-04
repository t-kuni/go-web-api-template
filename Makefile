clean:
	git ls-files --others --ignored --exclude-standard ./ | \
	grep -v '.env*' | \
	grep -v '.idea*' | \
	grep -v 'build-errors.log' | \
	xargs -I {} rm -f {}

generate: clean
	alias swagger='docker run --rm -it  --user $(id -u):$(id -g) -e GOPATH=$(go env GOPATH):/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger' && \
	swagger version && \
	go generate ./...

test: generate
	go test -tags feature -v ./...

coverage: generate
	go test -tags feature -coverpkg=./... -coverprofile=coverage/coverage.o ./... > /dev/null
	go tool cover -html=coverage/coverage.o -o coverage/coverage.html