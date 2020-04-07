# Go parameters
VERSION=0.1.0
build:
	@echo "Building binary.."
	@CGO_ENABLED=0 go build -v

docker:
	@echo "Building docker image.."
	docker build . -t yulrizka/redir:${VERSION}
	docker tag yulrizka/redir:${VERSION} yulrizka/redir:latest
