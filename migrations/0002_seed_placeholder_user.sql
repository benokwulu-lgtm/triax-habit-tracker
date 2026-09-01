INSERT INTO users (id, name, email, password_hash)
VALUES (1, 'Placeholder User', 'placeholder@example.local', 'placeholder')
ON CONFLICT (id) DO NOTHING;