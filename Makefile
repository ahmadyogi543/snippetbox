build:
	@go build -o ./bin/web ./cmd/web

run: build
	@cp -r tls bin/
	@./bin/web -addr=":5000"

clean:
	@rm -rf bin/
	@rm -rf tls/

test:
	@go test -v ./cmd/web

gen-tls:
	@mkdir tls
	@cd tls && go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
