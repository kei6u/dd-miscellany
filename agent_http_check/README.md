# Agent http check

https://docs.datadoghq.com/integrations/http_check/#pagetitle

## Docker

```bash
docker run -d --name dd-agent \
-v /var/run/docker.sock:/var/run/docker.sock:ro \
-v /proc/:/host/proc/:ro \
-v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
-v $(pwd)/conf.yaml:/conf.d/http_check.d/conf.yaml \
-e DD_API_KEY=${DD_API_KEY} \
gcr.io/datadoghq/agent:7
```
