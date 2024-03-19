# sheepdog-runner

Lambda runner that execute monitoring and alerting tasks

## Environment variable

```bash
PUSHGATEWAY: xxxx # pushgateway url
LOGLEVEL: DEBUG | INFO | WARN | ERROR  # log level default INFO
SQS_QUEUE_NAME: xxx # sqs queue name default : test
```

Another way is to setup a file `.env.local.json` with this format

```yaml
{
    "Runner": {
        "PUSHGATEWAY": "xxxx"
    }
}
```


## local release

```bash
goreleaser release --snapshot --clean
```

## launch local sam

lambda runner will be available at `http://localhost:3000/runner`

```bash
make sam-local
```
## launch localstack & send event with content of inputs.txt

```bash
make localstack
make call
make log # to follow lambda logs
```





