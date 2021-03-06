CREATE SEQUENCE IF NOT EXISTS contact_info_sequence;

CREATE TABLE IF NOT EXISTS contact_info (
    id BIGINT NOT NULL DEFAULT nextval('contact_info_sequence') PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    phoneNum VARCHAR(100) NOT NULL,
    address VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE INDEX contact_index ON contact_info(id, name);