# Events with a Custom Agent Check

This guide is based on [Events with a Custom Agent Check](https://docs.datadoghq.com/events/guides/agent/#submission).

## Docker

### Run Agent as a container

Replace the value of `DD_API_KEY` with your api key.

```bash
export DD_API_KEY=***
```

```bash
docker run -d --name dd-agent \
-v /var/run/docker.sock:/var/run/docker.sock:ro \
-v /proc/:/host/proc/:ro \
-v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
-v $(pwd)/event_example.yaml:/etc/datadog-agent/conf.d/event_example.d/event_example.yaml \
-v $(pwd)/event_example.py:/etc/datadog-agent/checks.d/event_example.py \
-e DD_API_KEY=${DD_API_KEY} \
-e DD_LOGS_ENABLED=true \
gcr.io/datadoghq/agent:7
```
