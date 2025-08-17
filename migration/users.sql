CREATE TABLE IF NOT EXISTS users(
  id SERIAL PRIMARY KEY,
  name VARCHAR(64),
  surname VARCHAR(64),
  email VARCHAR(64),
  password VARCHAR(64)
)
