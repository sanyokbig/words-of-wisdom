download-quotes:
	curl https://raw.githubusercontent.com/JamesFT/Database-Quotes-JSON/master/quotes.json > ./data/quotes.json

gen-json:
	cd internal/quotes-dispenser && easyjson quotes-dispenser.go
	cd internal/message && easyjson message.go

lint:
	golangci-lint run

build-server:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./words-of-wisdom-server ./cmd/server
	docker build -t words-of-wisdom-server -f ./build/server.Dockerfile .

build-client:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./words-of-wisdom-client ./cmd/client
	docker build -t words-of-wisdom-client -f ./build/client.Dockerfile .

run-server:
	docker rm -f words-of-wisdom-server | true
	docker run \
		--env-file server.env \
		--name words-of-wisdom-server \
		-p 8000:8000 \
		words-of-wisdom-server:latest

run-client:
	docker rm -f words-of-wisdom-client | true
	docker run \
		--env-file client.env \
		--name words-of-wisdom-client \
		words-of-wisdom-client:latest

