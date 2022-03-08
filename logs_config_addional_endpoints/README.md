# NOT WORKING PROPERLY

```bash
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout go-server.key -out go-server.crt
```

```bash
docker compose up
```

TODO: Fix warning message.

```
==========
Logs Agent
==========

    Reliable: Sending uncompressed logs in SSL encrypted TCP to agent-intake.logs.datadoghq.com on port 10516
    Unreliable: Sending compressed logs in SSL encrypted TCP to logger on port 50000

    You are currently sending Logs to Datadog through TCP (either because logs_config.use_tcp or logs_config.socks5_proxy_address is set or the HTTP connectivity test has failed). To benefit from increased reliability and better network performances, we strongly encourage switching over to compressed HTTPS which is now the default protocol.

    BytesSent: 479772
    EncodedBytesSent: 479772
    LogsProcessed: 1231
    LogsSent: 1231

  Warnings
  ========
    Connection to the log intake cannot be established: x509: certificate is not valid for any names, but wanted to match logger
```
