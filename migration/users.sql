CREATE TABLE IF NOT EXISTS users(
  id SERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  surname VARCHAR(64) NOT NULL,
  email VARCHAR(64) NOT NULL,
  password VARCHAR(64) NOT NULL,
  discordNotifications BOOLEAN DEFAULT FALSE,
  emailNotifications BOOLEAN DEFAULT FALSE,
  slackNotifications BOOLEAN DEFAULT FALSE
)
-- postgresql://myuser:mypassword@localhost:5432/mydb?sslmode=disable
