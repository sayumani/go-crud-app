-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE item
(
    item_id serial PRIMARY KEY,
    name VARCHAR ( 50 ) NOT NULL,
    rating INT NOT NULL,
    category VARCHAR ( 50 ) NOT NULL,
    image TEXT NOT NULL,
    reputation INT NOT NULL,
    price INT NOT NULL,
    availability INT NOT NULL
);

CREATE TABLE item_location
(
    id_location serial PRIMARY KEY,
    item_id INT NOT NULL,
    city VARCHAR ( 50 ) NOT NULL,
    state VARCHAR ( 50 ) NOT NULL,
    country VARCHAR ( 50 ) NOT NULL,
    zip_code INT NOT NULL,
    address VARCHAR ( 250 ) NOT NULL,
    CONSTRAINT fk_item
        FOREIGN KEY(item_id) 
	    REFERENCES item(item_id)
        ON DELETE CASCADE
);

CREATE TABLE item_booking
(
    id_booking serial PRIMARY KEY,
    item_id INT NOT NULL,
    person_name VARCHAR ( 50 ) NOT NULL,
    no_of_rooms INT NOT NULL,
    CONSTRAINT fk_item
        FOREIGN KEY(item_id) 
	    REFERENCES item(item_id)
        ON DELETE CASCADE
);


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE item;

DROP TABLE item_location;

DROP TABLE TABLE item_booking;