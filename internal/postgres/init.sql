CREATE TABLE IF NOT EXISTS user_credentials (
user_id SERIAL PRIMARY KEY,
username VARCHAR(50) UNIQUE NOT NULL,
salt VARCHAR(100) NOT NULL,
hash VARCHAR(150) NOT NULL,
email VARCHAR(100) UNIQUE NOT NULL
);