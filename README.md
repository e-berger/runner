# sheepdog-runner

Lambda runner that execute monitor

## Environment variable

```bash
TURSO_TOKEN: xxx
TURSO_DATABASE: xxxx  # only database name
PUSHGATEWAY: xxxx # pushgateway url
```

Another way is to setup a file `.env.local.json` with this format

```yaml
{
    "Runner": {
        "TURSO_TOKEN": "xxxx",
        "TURSO_DATABASE": "xxxx"
        "PUSHGATEWAY": "xxxx"
    }
}
```


## local release

```bash
goreleaser release --snapshot --clean
```

## lauch local sam

lambda runner will be available at `http://localhost:3000/runner`

```bash
make sam-local
```




