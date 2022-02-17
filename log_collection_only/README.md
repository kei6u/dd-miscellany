# Use the Datadog Agent for Log Collection Only

This guide is based on [Use the Datadog Agent for Log Collection Only](https://docs.datadoghq.com/logs/guide/how-to-set-up-only-logs/?tab=host).

## Run

Replace the value of `DD_API_KEY` with your api key.

```bash
export DD_API_KEY=***
```

```bash
docker run -d --name dd-agent \
-v /var/run/docker.sock:/var/run/docker.sock:ro \
-v /proc/:/host/proc/:ro \
-v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
-v $(pwd)/datadog.yaml:/etc/datadog-agent/datadog.yaml \
-e DD_API_KEY=${DD_API_KEY} \
gcr.io/datadoghq/agent:7
```
