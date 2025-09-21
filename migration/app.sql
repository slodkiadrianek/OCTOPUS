CREATE TABLE IF NOT EXISTS apps(
  id SERIAL PRIMARY KEY,
  name VARCHAR(64) NOT NULL,
  description VARCHAR(256),
  dbLink VARCHAR(128) ,
  ownerId INT  REFERENCES  users(id),
  slackWebhook VARCHAR(256),
  discordWebhook VARCHAR(256),
)
