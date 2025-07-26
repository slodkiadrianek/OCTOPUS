# Config/cahe.go


- create rate-limiter
- Rate limiter have to be nuanced

# Router (internal/api/routes/router.go)

- **Implement Path Parameter Support:** The current router uses a simple map for routes, which doesn't support dynamic path parameters (e.g., `/users/:id`). This is a critical feature for REST APIs. You could use a more advanced routing library or implement your own logic to handle this.

- **Add Route Grouping:** Introduce the concept of route groups to apply middleware and path prefixes to a set of routes. This helps in organizing the code and avoiding repetition.
  ```go
  // Example of how it could be used
  api := r.Group("/api/v1")
  api.Use(AuthMiddleware)
  api.Get("/users", getAllUsers)
  ```

- **Implement a "Method Not Allowed" Handler:** Currently, if a route exists but the wrong HTTP method is used, the router returns a 404 Not Found. It should return a 405 Method Not Allowed with an `Allow` header specifying the allowed methods.

- **Add Automatic `OPTIONS` Method Handling:** For routes that have at least one method defined, the router should automatically handle `OPTIONS` requests and respond with the `Allow` header. This is important for CORS pre-flight requests.

- **Integrate CORS Middleware:** Use the `cors.go` middleware you have. You can apply it globally or on a per-group/per-route basis.

- **Refactor HTTP Method Functions:** The `Post`, `Patch`, `Put`, and `Delete` methods have repetitive code for adding the `MethodCheckMiddleware`. This can be refactored into a helper function. Also, the `Delete` method uses the string `"Delete"` instead of `http.MethodDelete`. This should be corrected. The `Get` method has the `MethodCheckMiddleware` commented out, which should be addressed.

- **Add a "Not Found" Handler Function:** Instead of handling the 404 case directly in `ServeHTTP`, allow a custom `http.Handler` to be set for handling not-found routes.

- **Static File Serving:** Add a function to serve static files from a directory. This is useful for serving assets for a web application.

- **Graceful Shutdown:** While not strictly a router function, ensuring the server that uses the router can shut down gracefully is important. The `api/main.go` should be checked for this.