-- корзина пользователя
CREATE TABLE IF NOT EXISTS shop.bucket( 
    id SERIAL PRIMARY KEY,
    fk_user_id INT NOT NULL REFERENCES shop.users(id), -- кто купил
    fk_merch_id INT NOT NULL REFERENCES shop.merch(id), --тип мерча 
    cnt INT NOT NULL,
    UNIQUE (fk_user_id, fk_merch_id)
);