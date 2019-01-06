PROJECT_NAME=lambda-s3-cloudflare
BUILD_OUTPUT=app
REGION=us-west-2

default: fmt test build

build:
	GOOS=linux go build -o ${BUILD_OUTPUT}

fmt:
	go fmt

test:
	go test ./...

serverless: default
	sam local invoke "${BUILD_OUTPUT}" --event event.json --profile logtest --region ${REGION}

debug: default
	dlv debug . --listen=:2345 --headless --log --api-version=2 -- server

zip: default
	zip handler.zip ./${BUILD_OUTPUT}

create: zip
	aws lambda create-function \
		--region ${REGION} \
		--function-name ${PROJECT_NAME} \
		--memory 128 \
		--role arn:aws:iam::515609462839:role/lambda-s3-bucket-policy \
		--runtime go1.x \
		--zip-file fileb://./handler.zip \
		--handler ${BUILD_OUTPUT} \
		--profile lambda-devops

deploy: zip
	aws lambda update-function-code \
		--function-name ${PROJECT_NAME} \
		--zip-file fileb://./handler.zip \
		--profile lambda-devops

clean:
	rm ${BUILD_OUTPUT}
	rm handler.zip
