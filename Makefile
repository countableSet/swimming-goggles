PROJECT_NAME=lambda-s3-cloudflare

default: fmt test build

build:
	GOOS=linux go build -o app

fmt:
	go fmt

test:
	go test ./...

serverless: default
	sam local invoke "app" --no-event --profile logtest --region us-west-2

debug: default
	dlv debug . --listen=:2345 --headless --log --api-version=2 -- server

clean:
	rm app