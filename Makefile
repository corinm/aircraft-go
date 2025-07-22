SHELL = /bin/bash

unit-tests:
	cd enricher && go test -v ./...

generate:
	cd historian && rm -rf ./db/pg/* && sqlc generate
