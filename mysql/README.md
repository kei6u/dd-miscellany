# Agent based MySQL integration

This instruction is based on [Integration/MySQL](https://docs.datadoghq.com/integrations/mysql/?tab=host) to get `mysql.*` metrics.

Run Agent and MySQL containers.

```shell
docker compose up -d
```

You can see the custom metric named `test.cols.col1` from the query, `select * from test.cols`.
