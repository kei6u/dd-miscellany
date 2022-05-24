CREATE USER 'datadog'@'%' IDENTIFIED BY 'password';
GRANT REPLICATION CLIENT ON *.* TO 'datadog'@'%';
ALTER USER 'datadog'@'%'
WITH MAX_USER_CONNECTIONS 5;
GRANT PROCESS ON * . * TO 'datadog'@'%';
GRANT SELECT ON *.* TO datadog@'%';

CREATE DATABASE IF NOT EXISTS test;
CREATE TABLE test.cols (col1 int, col2 int, col3 int);
INSERT INTO test.cols (col1, col2, col3) VALUES (10, 20, 30);
INSERT INTO test.cols (col1, col2, col3) VALUES (11, 21, 31);
INSERT INTO test.cols (col1, col2, col3) VALUES (12, 22, 32);
INSERT INTO test.cols (col1, col2, col3) VALUES (13, 23, 33);
