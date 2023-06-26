CREATE TABLE todos (
    todos_id BIGINT NOT NULL AUTO_INCREMENT,
    title VARCHAR(99) NOT NULL,
    description VARCHAR(999) NOT NULL,
    status TINYINT DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- user_id BIGINT NOT NULL,
    PRIMARY KEY (todos_id)
    -- FOREIGN KEY (user_id) REFERENCES users(user_id)
);