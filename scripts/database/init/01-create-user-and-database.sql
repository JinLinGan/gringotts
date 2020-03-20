CREATE DATABASE gringotts CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'gringotts'@'%' IDENTIFIED BY 'gringotts';
GRANT ALL PRIVILEGES ON gringotts . * TO 'gringotts'@'%';
FLUSH PRIVILEGES;