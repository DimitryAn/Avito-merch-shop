CREATE TABLE IF NOT EXISTS shop.operations(
    id SERIAL PRIMARY KEY,
    fk_id_sender INT NOT NULL REFERENCES shop.users(id), -- кто отправил
    fk_id_getter INT NOT NULL REFERENCES shop.users(id), -- кто получил
    amount INT CHECK (amount > 0),
    CHECK (fk_id_sender != fk_id_getter),
    UNIQUE (fk_id_sender, fk_id_getter)
);