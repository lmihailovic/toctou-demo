CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(70) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    balance NUMERIC(12,2) NOT NULL DEFAULT 0.00
);

CREATE TABLE IF NOT EXISTS transfers (
    id SERIAL PRIMARY KEY,
    sender_id INTEGER NOT NULL REFERENCES users(id),
    recipient_id INTEGER NOT NULL REFERENCES users(id),
    amount NUMERIC(12,2) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO users (id, email, password, first_name, last_name, is_admin, balance)
VALUES
    (1, 'lobradovic@mail.com', 'password1', 'Luka', 'Obradovic', true, 17000.00),
    (2, 'lmihailovic@mail.com', 'password2', 'Luka', 'Mihailovic', false, 17000.00),
    (3, 'vlazarevic@mail.com', 'password3', 'Vukadin', 'Lazarevic', false, 17000.00)
ON CONFLICT (email) DO NOTHING;

SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

INSERT INTO transfers (sender_id, recipient_id, amount)
SELECT 2, 3, 1000.00
WHERE EXISTS (SELECT 1 FROM users WHERE id = 2)
  AND EXISTS (SELECT 1 FROM users WHERE id = 3);

INSERT INTO transfers (sender_id, recipient_id, amount)
SELECT 3, 2, 1000.00
WHERE EXISTS (SELECT 1 FROM users WHERE id = 3)
  AND EXISTS (SELECT 1 FROM users WHERE id = 2);
