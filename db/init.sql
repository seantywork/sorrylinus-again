CREATE USER 'seantywork'@'%' IDENTIFIED BY 'letsshareitwiththewholeuniverse';

GRANT ALL PRIVILEGES ON *.* TO 'seantywork'@'%';

CREATE DATABASE sorrylinusdb;

USE sorrylinusdb;

CREATE TABLE users (
    id VARCHAR(128), 
    pw VARCHAR(256), 
    sess VARCHAR(128), 
    PRIMARY KEY(id)
    );

CREATE TABLE admin (
    id VARCHAR(128), 
    PRIMARY KEY(id)
    );


COMMIT;