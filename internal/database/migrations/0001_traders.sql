-- migrate:up
CREATE TABLE traders (
    id INTEGER PRIMARY KEY,
    balance DECIMAL NOT NULL DEFAULT 0.00
);

-- migrate:down
DROP TABLE traders
