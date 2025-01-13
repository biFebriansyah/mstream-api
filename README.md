# mstream-api

A backend service for a music streaming application that implements HLS (HTTP Live Streaming). Designed and developed for research and development purposes, with a focus on scalability, efficiency, and modern communication protocols.

## Features

-   FFMPEG for processing, encoding, and segmenting of media files .
-   CRUD operations for resource.
-   Integration with RabbitMQ for queue.
-   Real-time updates using WebSockets.
-   Secure and optimized for production.

## Tech Stack

-   **Language:** Golang
-   **Framework:** Go-fiber
-   **Database:** Postgresql
-   **Queue:** RabbitMQ
-   **Socket:** Websocket
-   **Other Tools:** Docker, Swagger, Jest

## Getting Started

### Prerequisites

-   Go (v1.23+)
-   Docker (optional)
-   Postgres instance (local or cloud)
-   RabbitMQ instance

### Project setup

```bash
# Clone the repository
$ git clone https://github.com/biFebriansyah/yapping-api-gateway.git

# Install Package
$ go get -u ./...

# using Make
$ make install
```

### Compile and run the project

```bash
# development
$ go run *.go

# watch mode
$ make run

# production mode
$ make build
```

### Run Migration

```bash
# create migration file
$ make migrate-init name=mstream-name-migration

# exec migration
$ make migrate-up

# reset migration
$ make migrate-reset
```

## Authors

-   [@biFebriansyah](https://www.github.com/biFebriansyah)

## License

[MIT](https://choosealicense.com/licenses/mit/)
