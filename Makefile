BINARY := provisium
VERSION :=`cat VERSION`
.DEFAULT_GOAL := linux

linux:
	cd cmd/$(BINARY) ; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o $(BINARY)

docker:
	docker build --tag="nsfearthcube/thprov:$(VERSION)"  --file=./build/Dockerfile . ; \
	docker tag nsfearthcube/thprov:$(VERSION) earthcube/thprov:latest

removeimage:
	docker rmi --force nsfearthcube/thprov:$(VERSION)
	docker rmi --force nsfearthcube/thprov:latest

publish: docker
	docker push nsfearthcube/thprov:$(DOCKERVER)
