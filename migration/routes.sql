CREATE TYPE requestMethod AS ENUM ('POST', 'GET', 'PUT', 'PATCH', 'DELETE');

CREATE TABLE IF NOT EXISTS routes(
  id SERIAL PRIMARY KEY,
  route VARCHAR(256) NOT NULL,
  method requestMethod NOT NULL DEFAULT 'GET',
  queryData JSON,
  paramData JSON,
  bodyData JSON,
  predictedStatusCode INT NOT NULL DEFAULT 200,
  predictedBodyData JSON,
  appId INT REFERENCES apps(id)
)
