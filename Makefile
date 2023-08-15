build:
	@go build -o ./bin/web ./cmd/web

run: build
	@cp -r tls bin/
	@./bin/web -addr=":5000"

clean:
	@rm -rf bin/

test:
	@go test -v ./cmd/web
