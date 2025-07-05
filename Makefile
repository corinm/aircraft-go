SHELL = /bin/bash

unit-tests:
	cd enricher && go test -v ./...
