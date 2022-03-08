# Collecting Live Processes

This guide is based on [Live Processes](https://docs.datadoghq.com/infrastructure/process/?tab=docker).

## Run

Replace the value of `DD_API_KEY` with your api key.

```bash
export DD_API_KEY=***
```

```bash
docker run -d --name dd-agent \
-v /var/run/docker.sock:/var/run/docker.sock:ro \
-v /proc/:/host/proc/:ro \
-v /etc/passwd:/etc/passwd:ro \
-v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
-e DD_PROCESS_AGENT_ENABLED=true \
-e DD_API_KEY=${DD_API_KEY} \
gcr.io/datadoghq/agent:7
```
