# Go Recipes
A collection of practical Go recipes, examples, and reusable backend solutions.
This repository explores Go development through practical implementations and backend technologies.


## Tech Stack
- **Go (Golang)**
- **Fiber** – Web framework
- **PostgreSQL** – Database
- **JWT** – Authentication
- **Docker** – Containerization


## Getting Started
### 1. Start Database
```bash
docker compose up -d db
```
### 2. Run in Debug Mode
Open the project in VS Code and run the configured **Launch JSON** from the Run and Debug panel.
### 3. Run with Docker
To build and run the entire application using Docker Compose:
```bash
docker compose up -d --build
```