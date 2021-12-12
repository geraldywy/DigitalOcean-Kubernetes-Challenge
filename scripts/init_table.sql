CREATE TABLE IF NOT EXISTS example_db.messages (
    `author` VARCHAR(50) NOT NULL,
    `message` TEXT NULL DEFAULT NULL,
    PRIMARY KEY (`author`)
);