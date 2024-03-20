-- +goose Up
-- +goose StatementBegin
CREATE TABLE racks (
    id serial primary key,
    name text
);

CREATE TABLE products (
   id serial primary key,
   name text
);

CREATE TABLE main_racks (
   product_id int references products(id),
   rack_id int references racks(id)
);

CREATE TABLE secondary_racks (
    product_id int references products(id),
    rack_id int references racks(id)
);

CREATE TABLE orders (
   id serial primary key
);

CREATE TABLE order_products (
    id int references orders(id),
    product_id int references products(id),
    amount int
);

INSERT INTO products VALUES
    (1, 'Ноутбук'),
    (2, 'Телевизор'),
    (3, 'Телефон'),
    (4, 'Системный блок'),
    (5, 'Часы'),
    (6, 'Микрофон');

INSERT INTO racks(name) VALUES
    ('А'),
    ('Б'),
    ('Ж'),
    ('З'),
    ('В');

INSERT INTO orders VALUES
    (10),
    (11),
    (14),
    (15);

INSERT INTO main_racks VALUES
   (1, 1),
   (2, 1),
   (3, 2),
   (4, 3),
   (5, 3),
   (6, 3);

INSERT INTO secondary_racks VALUES
    (3, 4),
    (3, 5),
    (5, 1);

INSERT INTO order_products VALUES
    (10, 1, 2),
    (10, 3, 1),
    (10, 6, 1),
    (11, 2, 3),
    (14, 1, 3),
    (14, 4, 4),
    (15, 5, 1);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE main_racks;
DROP TABLE secondary_racks;
DROP TABLE order_products;
DROP TABLE orders;
DROP TABLE racks;
DROP TABLE products;
-- +goose StatementEnd
