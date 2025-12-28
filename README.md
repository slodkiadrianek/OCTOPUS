# OCTOPUS

OCTOPUS is a backend application written in Go. It gives you opportunity to control your self hosted apps.

## Project Structure

The project is organized into the following main directories:

-   `cmd`: Contains the main applications for the different components of the system.
    -   `api`: The main entry point for the API server.
    -   `worker`: The main entry point for the background worker.
-   `internal`: Contains the internal logic of the application, separated by domain.
    -   `api`: Handles API-related logic, including routing, handlers, and server setup.
    -   `config`: Manages application configuration, including database and cache settings.
    -   `DTO`: Data Transfer Objects used for passing data between layers.
    -   `middleware`: Contains HTTP middleware for things like CORS, error handling, and rate limiting.
    -   `models`: Defines the data models for the application.
    -   `repository`: The layer responsible for database interactions.
    -   `schema`: Contains database schema definitions.
    -   `services`: Implements the business logic of the application.
    -   `utils`: Provides utility functions, such as logging and helpers.
-   `deployments`: Contains deployment configurations, such as Docker files.
-   `docs`: For project documentation.
-   `tests`: For application tests.

## Getting Started

To get started with the project, you will need to have Go installed on your system.

1.  **Clone the repository:**

    ```bash
    git clone https://github.com/slodkiadrianek/OCTOPUS
    ```

2.  **Install dependencies:**

    ```bash
    go mod tidy
    ```

3.  **Run the application:**

    To run a specific component, use the `go run` command:

    ```bash
    # Run the API server
    go run cmd/api/main.go

    # Run the worker
    go run cmd/worker/main.go

    ```

## Docker 
1.  **Go to proper the proper directory:**

    ```bash
    # API docker config location
    cd deployments/docker/api/ 

    # Worker docker config location
    cd deployments/docker/worker/
    ```

2.  **Build the app image:**

    ```bash
    # Build API image
    docker build -f Dockerfile.api -t octopus-api ../../../

    # Build worker image
    docker build -f Dockerfile.worker -t octopus-worker ../../../
    ```

3.  **Run the application as the container:**

    ```bash
    # Run API image
    docker run -d -p 3040:3040 --name octopus-api octopus-api

    # Run worker image
    docker run -d -p 3041:3041 --name octopus-worker octopus-worker
    ```

## Configuration

This project uses a `.env` file for configuration. Create a `.env` file in the root of the project and add the following variables:

```
Port=8080
JWTSecret=your-secret
DBLink=your-db-link
CacheLink=your-cache-link
DockerHost=your-docker-host
```

## Features

- Create account
- Rate limiter
- Import docker containers
- Add apps from hand
- Checking statuses of apps
- You can get notifications through webhooks like slack or discord
- Check server info
- Get server metrics
- Get routes responses in background using worker
- Get routes statuses

## documentation

Swagger file is in docs directory.

## Testing

Under development.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
