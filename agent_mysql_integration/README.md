# Agent based MySQL integration

This instruction is based on [Integration/MySQL](https://docs.datadoghq.com/integrations/mysql/?tab=host) to get `mysql.*` metrics.

## On your machine

Run Agent, MySQL containers.

```shell
docker compose up -d
```

Attach the MySQL container.

```shell
docker exec -it <CONTAINER_ID> /bin/bash
```

Login to MySQL on the container.

```shell
mysql -u root
```

Configuration for Datadog Agent.

Why do I use `192.168.80.3`? See [MySQL Localhost Error - Localhost VS 127.0.0.1](https://docs.datadoghq.com/integrations/faq/mysql-localhost-error-localhost-vs-127-0-0-1/)

```sql
CREATE USER 'datadog'@'192.168.80.3' IDENTIFIED BY 'password';
GRANT REPLICATION CLIENT ON *.* TO 'datadog'@'192.168.80.3' WITH MAX_USER_CONNECTIONS 5;
GRANT PROCESS ON *.* TO 'datadog'@'192.168.80.3';
ALTER USER 'datadog'@'192.168.80.3' WITH MAX_USER_CONNECTIONS 5;
```
