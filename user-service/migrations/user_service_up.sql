CREATE TABLE IF NOT EXISTS users
(
    id           serial primary key,
    name         varchar NOT NULL,
    surname      varchar NOT NULL DEFAULT '********',
    phone_number varchar NOT NULL UNIQUE,
    chat_id      float NOT NULL UNIQUE,
    role         varchar check ( role in ('ADMIN', 'USER') ),
    created_at timestamp DEFAULT now()
);