# Agent check Directory

This document guides how to implement [Agent check Directory](https://docs.datadoghq.com/integrations/directory/).
There are options for conf.yaml [here](https://github.com/DataDog/integrations-core/blob/master/directory/datadog_checks/directory/data/conf.yaml.example).

## Docker

```bash
docker run -d --name dd-agent \
-v /var/run/docker.sock:/var/run/docker.sock:ro \
-v /proc/:/host/proc/:ro \
-v /sys/fs/cgroup/:/host/sys/fs/cgroup:ro \
-v $(pwd)/conf.yaml:/conf.d/directory.d/conf.yaml \
-v $(pwd)/alice/dir:/var/run/alice/dir \
-v $(pwd)/bob/dir:/var/run/bob/dir \
-v $(pwd)/john/dir:/var/run/john/dir \
-e DD_API_KEY=${DD_API_KEY} \
gcr.io/datadoghq/agent:7
```

Ignore changes in `test.log`.

```bash
git update-index --assume-unchanged alice
git update-index --assume-unchanged bob
git update-index --assume-unchanged john
```
