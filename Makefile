PROJECT_NAME=lambda-s3-cloudflare

default:
	go build

fmt:
	go fmt

test:
	go test ./...

debug: default
	dlv debug . --listen=0.0.0.0:2345 --headless --log --api-version=2 -- server

clean:
	rm ${PROJECT_NAME}