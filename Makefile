.PHONY: build localstack logs sam-local build-Runner call purge

queue_name = Events
lambda_name = sheepdog-runner
event_queue = sheepdog-dispatcher
endpoint = http://localhost:4566
region = us-east-1
archi = $(shell arch)

statuslocalhost = $(shell curl --write-out %{http_code} --silent --output /dev/null ${endpoint})
time = $(shell date -Iseconds)

build:
	go env -w GOPRIVATE='github.com/e-berger/*'
	goreleaser release --snapshot --clean

init:
	@if [ "$(statuslocalhost)" != "200" ]; then\
		docker-compose up -d;\
	fi

deploystart: deploy deployasia

deploy: build
	@aws --region $(region) --endpoint-url $(endpoint) lambda delete-function --function-name $(lambda_name)  2>/dev/null 1>/dev/null || true
	@aws --region $(region) --endpoint-url=$(endpoint) events create-event-bus --name ${event_queue} --tags "Key"="test","Value"="test" 2>/dev/null 1>/dev/null || true
	@aws --region $(region) --endpoint-url=$(endpoint) events put-rule --name ${event_queue} --event-bus-name $(event_queue) \
	--event-pattern "{\"detail\":{\"location\":[\"europe\"]},\"source\":[\"sheepdog-dispatcher\"]}" 2>/dev/null 1>/dev/null || true
	@aws --region $(region) --endpoint-url=$(endpoint) lambda create-function --function-name $(lambda_name) \
	--zip-file fileb://dist/$(lambda_name)_Linux_$(archi).zip \
	--architecture $(archi) \
	--handler bootstrap --runtime go1.x \
	--role arn:aws:iam::000000000000:role/lambda-role \
	--timeout 900 \
	--description "$(time)" \
	--environment Variables="{SQS_QUEUE_NAME=Events,LOGLEVEL=debug,AWS_REGION_CENTRAL=us-east-1,CLOUDWATCHPREFIX=/probe}" | jq
	@sleep 1
	@aws --region $(region) --endpoint-url=$(endpoint) events put-targets --rule ${event_queue} --event-bus-name $(event_queue) \
	--targets "Id"="1","Arn"="arn:aws:lambda:us-east-1:000000000000:function:$(lambda_name)" 2>/dev/null | jq || true

deployasia: build
	@aws --region $(region) --endpoint-url $(endpoint) lambda delete-function --function-name $(lambda_name)-asia  2>/dev/null 1>/dev/null || true
	@aws --region $(region) --endpoint-url=$(endpoint) events create-event-bus --name ${event_queue} --tags "Key"="test","Value"="test" 2>/dev/null 1>/dev/null || true
	@aws --region $(region) --endpoint-url=$(endpoint) events put-rule --name ${event_queue}-asia --event-bus-name $(event_queue) \
	--event-pattern "{\"detail\":{\"location\":[\"asia\"]},\"source\":[\"sheepdog-dispatcher\"]}" 2>/dev/null 1>/dev/null || true
	@aws --region $(region) --endpoint-url=$(endpoint) lambda create-function --function-name $(lambda_name)-asia \
	--zip-file fileb://dist/$(lambda_name)_Linux_$(archi).zip \
	--architecture $(archi) \
	--handler bootstrap --runtime go1.x \
	--role arn:aws:iam::000000000000:role/lambda-role \
	--timeout 900 \
	--description "$(time)" \
	--environment Variables="{SQS_QUEUE_NAME=Events,LOGLEVEL=debug,AWS_REGION_CENTRAL=us-east-1,CLOUDWATCHPREFIX=/probe}" | jq
	@sleep 1
	@aws --region $(region) --endpoint-url=$(endpoint) events put-targets --rule ${event_queue}-asia --event-bus-name $(event_queue) \
	--targets "Id"="1","Arn"="arn:aws:lambda:us-east-1:000000000000:function:$(lambda_name)-asia" 2>/dev/null | jq || true

localstack: build init deploy

logs:
	aws --region $(region) --endpoint-url=http://localhost:4566 logs tail "/aws/lambda/$(lambda_name)" --follow

sam-local:
	sam build
	sam local start-api --env-vars .env.local.json

build-Runner:
	goreleaser release --snapshot --clean
	cp ./dist/sheepdog-runner_linux_amd64_v1/bootstrap $(ARTIFACTS_DIR)/.

call:
	aws --region $(region) --endpoint-url=$(endpoint) lambda invoke --function-name $(lambda_name) --cli-binary-format raw-in-base64-out --payload file://inputs.txt /dev/stdout

callhttp:
	aws --region $(region) --endpoint-url=$(endpoint) lambda invoke --function-name $(lambda_name) --cli-binary-format raw-in-base64-out --payload file://inputs_http.txt /dev/stdout

purge:
	aws --region $(region) --endpoint-url=$(endpoint) sqs purge-queue --queue-url http://localhost:4566/000000000000/${queue_name}

test:
	go test -v ./...
