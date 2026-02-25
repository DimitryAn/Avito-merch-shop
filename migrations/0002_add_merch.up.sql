CREATE TABLE IF NOT EXISTS shop.merch(
    id SERIAL PRIMARY KEY,
    name_type TEXT NOT NULL UNIQUE,
    price INT NOT NULL
);

INSERT INTO shop.merch (name_type,price)
VALUES ('t-shirt', 80);

INSERT INTO shop.merch (name_type,price)
VALUES ('cup', 20);

INSERT INTO shop.merch (name_type,price)
VALUES ('book', 50);

INSERT INTO shop.merch (name_type,price)
VALUES ('pen', 10);

INSERT INTO shop.merch (name_type,price)
VALUES ('powerbank', 200);

INSERT INTO shop.merch (name_type,price)
VALUES ('hoody', 300);

INSERT INTO shop.merch (name_type,price)
VALUES ('umbrella', 200);

INSERT INTO shop.merch (name_type,price)
VALUES ('socks', 10);

INSERT INTO shop.merch (name_type,price)
VALUES ('wallet', 50);

INSERT INTO shop.merch (name_type,price)
VALUES ('pink-hoody', 500);


