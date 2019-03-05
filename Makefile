BINARY := provisium
VERSION :=`cat VERSION`
.DEFAULT_GOAL := linux

linux:
	cd cmd/$(BINARY) ; \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 env go build -o $(BINARY)

docker:
	docker build --tag="esip/$(BINARY):$(VERSION)"  --file=./build/Dockerfile . ; \
	docker tag esip/$(BINARY):$(VERSION) esip/thprov:latest

removeimage:
	docker rmi --force esip/$(BINARY):$(VERSION)
	docker rmi --force esip/$(BINARY):latest

publish: docker
	docker push esip/$(BINARY):$(VERSION)
