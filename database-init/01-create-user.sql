CREATE DATABASE gringotts;
CREATE USER 'gringotts'@'%' IDENTIFIED BY 'gringotts';
GRANT ALL PRIVILEGES ON gringotts . * TO 'gringotts'@'%';
FLUSH PRIVILEGES;