download-quotes:
	curl https://raw.githubusercontent.com/JamesFT/Database-Quotes-JSON/master/quotes.json > ./data/quotes.json

gen-json:
	cd internal/quotes-dispenser && easyjson quotes-dispenser.go
	cd internal/message && easyjson message.go

lint:
	golangci-lint run
