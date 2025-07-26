# OCTOPUS

OCTOPUS is a backend application written in Go. It appears to be a multi-component system with an API, a worker, and an agent.

## Project Structure

The project is organized into the following main directories:

-   `cmd`: Contains the main applications for the different components of the system.
    -   `api`: The main entry point for the API server.
    -   `worker`: The main entry point for the background worker.
    -   `agent`: The main entry point for the agent.
-   `internal`: Contains the internal logic of the application, separated by domain.
    -   `api`: Handles API-related logic, including routing, handlers, and server setup.
    -   `worker`: Contains the logic for the background worker, including jobs and queues.
    -   `agent`: Holds the agent's core logic, such as data collection and sending.
    -   `config`: Manages application configuration, including database and cache settings.
    -   `DTO`: Data Transfer Objects used for passing data between layers.
    -   `middleware`: Contains HTTP middleware for things like CORS, error handling, and rate limiting.
    -   `models`: Defines the data models for the application.
    -   `repository`: The layer responsible for database interactions.
    -   `schema`: Contains database schema definitions.
    -   `services`: Implements the business logic of the application.
    -   `utils`: Provides utility functions, such as logging and helpers.
-   `pkg`: Contains packages that can be shared with other applications.
    -   `errors`: Defines custom error types.
    -   `notification`: A package for handling notifications.
-   `deployments`: Contains deployment configurations, such as Docker files.
-   `docs`: For project documentation.
-   `scripts`: For utility scripts.
-   `tests`: For application tests.

## Getting Started

To get started with the project, you will need to have Go installed on your system.

1.  **Clone the repository:**

    ```bash
    git clone <repository-url>
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

    # Run the agent
    go run cmd/agent/main.go
    ```

## Contributing

Please read the `todo.md` file for tasks that need to be done.