
lambda_name = sheepdog-runner
endpoint = http://localhost:4566

statuslocalhost = $(shell curl --write-out %{http_code} --silent --output /dev/null ${endpoint})
time = $(shell date --iso=seconds)

build:
	goreleaser release --snapshot --clean

localstack: build
	@if [ "$(statuslocalhost)" != "200" ]; then\
		docker-compose up -d;\
	fi
	aws --endpoint-url $(endpoint) lambda delete-function --function-name $(lambda_name) || true
	aws --endpoint-url=$(endpoint) \
	lambda create-function --function-name $(lambda_name) \
	--zip-file fileb://dist/$(lambda_name)_Linux_x86_64.zip \
	--handler bootstrap --runtime go1.x \
	--role arn:aws:iam::000000000000:role/lambda-role \
	--timeout 900 \
	--description "$(time)" \
	--environment Variables="{SQS_QUEUE_NAME=Events}" | jq

logs:
	aws --endpoint-url=http://localhost:4566 logs tail "/aws/lambda/$(lambda_name)" --follow

sam-local:
	sam build
	sam local start-api --env-vars .env.local.json

build-Runner:
	goreleaser release --snapshot --clean
	cp ./dist/sheepdog-runner_linux_amd64_v1/bootstrap $(ARTIFACTS_DIR)/.

call:
	awslocal --endpoint-url=$(endpoint) lambda invoke --function-name $(lambda_name) --cli-binary-format raw-in-base64-out --payload file://inputs.txt /dev/stdout

