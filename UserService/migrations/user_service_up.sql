CREATE TABLE users
(
    id           serial primary key,
    name         varchar NOT NULL,
    surname      varchar NOT NULL DEFAULT '********',
    phone_number varchar NOT NULL UNIQUE ,
    role varchar check ( role in ('ADMIN' , 'USER') )
);