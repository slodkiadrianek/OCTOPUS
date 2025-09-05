CREATE TYPE requestMethod AS ENUM ('POST', 'GET', 'PUT', 'PATCH', 'DELETE');

CREATE TABLE IF NOT EXISTS routes(
  id SERIAL PRIMARY KEY,
  path VARCHAR(256) NOT NULL,
  method requestMethod NOT NULL DEFAULT 'GET',
  UNIQUE(method, path)
)

CREATE TABLE IF NOT EXISTS routesResponses(
  id SERIAL PRIMARY KEY,
  statusCode INT NOT NULL DEFAULT 200,
  bodyData JSONB,
  UNIQUE(statusCode, bodyData)
)
CREATE TABLE IF NOT EXISTS routesRequests(
  id SERIAL PRIMARY KEY,
  queryData JSONB,
  paramData JSONB,
  bodyData JSONB,
  UNIQUE( queryData, paramData, bodyData)
)

CREATE TABLE IF NOT EXISTS workingRoutes(
    id SERIAL PRIMARY KEY,
    routeId INT NOT NULL REFERENCES routes(id),
    requestId INT NOT NULL REFERENCES routesRequests(id),
    responseId INT NOT NULL REFERENCES routesResponses(id),
    parentId INT REFERENCES workingRoutes(id),
    appId INT NOT NULL REFERENCES apps(id),
    UNIQUE(routeId, appId, parentId, requestId, responseId)
)
