# EXBanka Infrastructure

Infrastructure repository for the EXBanka microservices system.

This repository contains configuration required to run the system
locally using Docker.

------------------------------------------------------------------------

## Project structure

    EXBanka-3-Infrastructure
    ├─ docker-compose.yml
    ├─ .env.example
    ├─ README.md
    ├─ nginx/
    │   └─ nginx.conf
    └─ scripts/

------------------------------------------------------------------------

## Requirements

You need to install:

-   Docker Desktop
-   Git

Optional but recommended:

-   VS Code

------------------------------------------------------------------------

## Running the system

Clone all repositories in the same parent folder:

    project-folder
    ├─ EXBanka-3-Backend
    ├─ EXBanka-3-Frontend
    └─ EXBanka-3-Infrastructure

Go to the infrastructure repository:

    cd EXBanka-3-Infrastructure

Start the infrastructure:

    docker compose up

This will start:

-   PostgreSQL database
-   Nginx API gateway

------------------------------------------------------------------------

## Services

  Service         Port
  --------------- ------
  nginx gateway   8080
  postgres        5432

------------------------------------------------------------------------

## Testing infrastructure

Open in your browser:

    http://localhost:8080

If everything works you should see:

    Infrastructure running

------------------------------------------------------------------------

## Environment variables

Environment variables template:

    .env.example

Copy it when needed:

    cp .env.example .env

------------------------------------------------------------------------

## Future services

Later the following services will be connected:

-   auth-service
-   employee-service
-   notification-service
-   frontend

These services will be routed through the nginx gateway.

------------------------------------------------------------------------

## Authors

RAF Software Engineering 2025
