-- +goose Up
-- +goose StatementBegin
CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    fio VARCHAR(100),
    password_hash VARCHAR(255) NOT NULL,
    icon_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE private_contact (
    id SERIAL PRIMARY KEY,
    user_id INT,
    contact_user_id INT,
    status INT,
    
    CONSTRAINT check_user_order CHECK (user_id < contact_user_id),
    CONSTRAINT unique_user_contact UNIQUE (user_id, contact_user_id),
    CONSTRAINT fk_user_id 
        FOREIGN KEY (user_id) REFERENCES "user"(id),
    CONSTRAINT fk_contact_user_id 
        FOREIGN KEY (contact_user_id) REFERENCES "user"(id)
);

CREATE TABLE private_message (
    id SERIAL PRIMARY KEY,
    contact_id INT,
    user_id INT,
    data VARCHAR(1023),
    additionals TEXT[],
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_private_message_contact 
        FOREIGN KEY (contact_id) REFERENCES private_contact(id),
    CONSTRAINT fk_private_message_user 
        FOREIGN KEY (user_id) REFERENCES "user"(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS private_message;
DROP TABLE IF EXISTS private_contact;
DROP TABLE IF EXISTS "user";
-- +goose StatementEnd
