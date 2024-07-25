CREATE TABLE messages
(
    id            serial PRIMARY KEY,
    message       varchar(255) NOT NULL,
    processed     boolean NOT NULL
);

CREATE TABLE outboxmessages
(
    id            serial PRIMARY KEY,
    idmessage     integer NOT NULL
);

CREATE TABLE agregate_processed_messages
(
    date timestamp NOT NULL,
    id integer NOT NULL,
    UNIQUE (id),
    PRIMARY KEY (date, id)
);
