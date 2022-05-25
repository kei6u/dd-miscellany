# Host Agent Log collection

This guide is based on [Host Agent Log collection](https://docs.datadoghq.com/agent/logs/?tab=tailfiles).
You can configure Datadog Agent to collect logs from the specified log files.

```bash
# Run the Datadog Agent as container
docker compose up -d

# Write logs in test.log.1, test.log.2
go run main.go
```
