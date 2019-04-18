PROJECT_NAME=swimming-goggles
BUILD_OUTPUT=app

default: fmt test build

build:
	GOOS=linux go build -o ${BUILD_OUTPUT}

fmt:
	go fmt

test:
	go test ./...

serverless: default
	sam local invoke "${BUILD_OUTPUT}" --event event.json --profile logtest --region ${AWS_REGION}

debug: default
	dlv debug . --listen=:2345 --headless --log --api-version=2 -- server

zip:
ifeq ($(wildcard ./$(BUILD_OUTPUT)),)
	$(error Build output file not found, please run make again.)
else
	zip handler.zip ./${BUILD_OUTPUT}
endif

create: zip
	aws lambda create-function \
		--region ${AWS_REGION} \
		--function-name ${PROJECT_NAME} \
		--memory 128 \
		--role ${AWS_LAMBDA_S3_BUCKET_POLICY_ROLE_ARN} \
		--runtime go1.x \
		--zip-file fileb://./handler.zip \
		--handler ${BUILD_OUTPUT}

deploy: zip
	aws lambda update-function-code \
		--function-name ${PROJECT_NAME} \
		--zip-file fileb://./handler.zip \
		--publish

clean:
	rm ${BUILD_OUTPUT}
	rm handler.zip
