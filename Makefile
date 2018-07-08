.PHONY: test build
SHELL   := bash
NAME    := navikt/vks
LATEST  := ${NAME}:latest

push-dockerhub: docker-build
docker-build: build
build: test

clean:
	rm vks

test:
	go test ./...
build:
	go build -o vks

docker-build:
	docker image build -t ${NAME}:$(CIRCLE_BUILD_NUM) -t ${LATEST} -f Dockerfile .

push-dockerhub:
	docker image push ${NAME}:$(CIRCLE_BUILD_NUM) ${LATEST}
