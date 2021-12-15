CREATE USER 'datadog'@'%' IDENTIFIED BY 'password';
GRANT REPLICATION CLIENT ON *.* TO 'datadog'@'%';
ALTER USER 'datadog'@'%'
WITH MAX_USER_CONNECTIONS 5;
GRANT PROCESS ON * . * TO 'datadog'@'%';
