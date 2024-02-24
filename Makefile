
sam-local:
	sam build
	sam local start-api --env-vars .env.local.json

build-Runner:
	goreleaser release --snapshot --clean
	cp ./dist/sheepdog-runner_linux_amd64_v1/bootstrap $(ARTIFACTS_DIR)/.
