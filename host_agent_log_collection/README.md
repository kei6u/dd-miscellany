# Host Agent Log collection

This guide is based on [Host Agent Log collection](https://docs.datadoghq.com/agent/logs/?tab=tailfiles).
You can configure Datadog Agent to collect logs from the specified log files.

## Agent on your machine

### Run Agent as a container

Replace the value of `DD_API_KEY` with your api key.

```shell
docker run -d --name dd-agent \
-v /var/run/docker.sock:/var/run/docker.sock:ro \
-v /proc/:/host/proc/:ro \
-v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
-v $(pwd)/test.log:/test/test.log \
-v $(pwd)/conf.yaml:/conf.d/test.d/conf.yaml \
-v $(pwd)/datadog.yaml:/etc/datadog-agent/datadog.yaml \
-e DD_API_KEY=***** \
-e DD_ENV=mac \
gcr.io/datadoghq/agent:7
```

### Write messages to the log file

```shell
go run main.go --duration-ms=100 --silent --messages="hello,こんにちは,Hola"
```
