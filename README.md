# Microservices Project with Node.js, Go, and Python

This project is a demonstration of a microservices architecture using:
- **Node.js (NestJS)** for a general-purpose backend service.
- **Go (Gin)** for a high-performance service.
- **Python (FastAPI)** for an AI/ML service.
- **PostgreSQL** as the database.
- **pgAdmin** for database management.
- **Docker** and **Docker Compose** for containerization and orchestration.

## Project Structure

```
.
├── docker-compose.yml      # Orchestrates all the services
├── go-service/             # Go (Gin) microservice
│   ├── Dockerfile
│   ├── go.mod
│   └── main.go
├── nodejs-service/         # Node.js (NestJS) microservice
│   ├── Dockerfile
│   ├── package.json
│   ├── tsconfig.json
│   └── src/
│       ├── main.ts
│       ├── app.module.ts
│       ├── app.controller.ts
│       └── app.service.ts
├── python-ai-service/      # Python (FastAPI) microservice
│   ├── Dockerfile
│   ├── requirements.txt
│   └── main.py
└── README.md
```

## How to Run

1.  **Prerequisites:** Make sure you have [Docker](https://www.docker.com/products/docker-desktop) installed and running on your system.

2.  **Build and Run the Services:**
    Open a terminal in the project root and run the following command:
    ```bash
    docker-compose up --build
    ```
    This command will build the Docker images for each service (if they don't exist) and then start all the containers.

3.  **Access the Services:**
    Once all the containers are running, you can access the services at the following URLs:

    - **Node.js Service:** [http://localhost:3000](http://localhost:3000)
    - **Go Service:** [http://localhost:8080](http://localhost:8080)
    - **Python AI Service:** [http://localhost:8000](http://localhost:8000)
    - **pgAdmin:** [http://localhost:8888](http://localhost:8888)
      - **Login Email:** `admin@example.com`
      - **Login Password:** `admin`

      To connect to the PostgreSQL database from pgAdmin, you will need to add a new server with the following details:
      - **Host:** `postgres` (this is the service name from `docker-compose.yml`)
      - **Port:** `5432`
      - **Maintenance database:** `mydb`
      - **Username:** `user`
      - **Password:** `password`

## How to Stop

To stop all the running containers, press `Ctrl + C` in the terminal where `docker-compose` is running. To remove the containers and volumes, you can run:
```bash
docker-compose down -v
``` 