#Configuration Management API

A simple REST API for storing and managing configuration data. You can store **configs** (name + version + key-value parameters) and **config groups** (groups of configs with optional labels for filtering).

---

## What it does

- **Configs**: Store key-value settings (e.g. database host, port) with a name and version.
- **Config groups**: Group multiple configs together. Each config in a group can have **labels** (e.g. `environment:development`, `team:backend`) so you can list or delete configs by label.
- **Rate limiting**: Limits how many requests each client can make (token-bucket, per IP).
- **In-memory storage**: Data is kept in memory. Restarting the server clears it (no database).

---

## Prerequisites

- **Go 1.19** or newer  
- Optional: **Docker** and **Docker Compose** to run in a container

---

## Run locally

1. Open a terminal in the project folder and download dependencies:

   ```bash
   go mod download
   ```

2. Start the server:

   ```bash
   go run .
   ```

The API listens on **http://localhost:8000**.

---

## Run with Docker

From the project root:

```bash
docker-compose up --build
```

The API is available at **http://localhost:8000**. To run in the background:

```bash
docker-compose up -d --build
```

---

## Environment variables

| Variable           | Default | Description                                      |
|--------------------|---------|--------------------------------------------------|
| `RATE_LIMIT_RPS`   | `5`     | Approx. rate (tokens per 5 seconds) for limiting |
| `RATE_LIMIT_BURST` | `10`    | Max burst (bucket capacity) per client           |

---

## API overview

All responses are JSON. The server uses **rate limiting**; too many requests return `429 Too Many Requests` with a `Retry-After` header.

### Configs (standalone)

| Method | Path                      | Description              |
|--------|---------------------------|--------------------------|
| GET    | `/configs`                | List all configs         |
| GET    | `/configs/{name}/{version}` | Get one config        |
| POST   | `/configs`                | Create a config          |
| DELETE | `/configs/{name}/{version}` | Delete a config       |

**Example — create a config:**

```bash
curl -X POST http://localhost:8000/configs \
  -H "Content-Type: application/json" \
  -d '{"name":"db_config","version":1,"parameters":[{"key":"host","value":"localhost"},{"key":"port","value":"5432"}]}'
```

**Example — get a config:**

```bash
curl http://localhost:8000/configs/db_config/1
```

---

### Config groups

| Method | Path                                                    | Description                          |
|--------|---------------------------------------------------------|--------------------------------------|
| GET    | `/groups`                                               | List all groups                      |
| GET    | `/groups/{name}/{version}`                              | Get one group                        |
| POST   | `/groups`                                               | Create a group                       |
| DELETE | `/groups/{name}/{version}`                              | Delete a group                       |
| GET    | `/groups/{name}/{version}/configs`                      | List all configs in the group        |
| GET    | `/groups/{name}/{version}/configs?labels=k1:v1;k2:v2`   | List configs that match the labels   |
| GET    | `/groups/{name}/{version}/configs/{configName}`         | Get one config in the group          |
| POST   | `/groups/{name}/{version}/configs`                      | Add a config to the group            |
| DELETE | `/groups/{name}/{version}/configs/{configName}`         | Remove one config from the group     |
| DELETE | `/groups/{name}/{version}/configs?labels=k1:v1;k2:v2`   | Remove configs that match the labels |

**Labels** in query strings use the format: `key1:value1;key2:value2` (semicolon-separated).

**Example — create a group with a config:**

```bash
# Create group
curl -X POST http://localhost:8000/groups \
  -H "Content-Type: application/json" \
  -d '{"name":"web_configs","version":1,"configs":[]}'

# Add config to group
curl -X POST http://localhost:8000/groups/web_configs/1/configs \
  -H "Content-Type: application/json" \
  -d '{"name":"web_server","parameters":[{"key":"port","value":"8080"}],"labels":[{"key":"environment","value":"development"}]}'
```

**Example — get configs by labels:**

```bash
curl "http://localhost:8000/groups/web_configs/1/configs?labels=environment:development"
```

---

## Project structure

```
ars/
├── main.go              # Server setup, routes, rate limiter
├── go.mod / go.sum      # Go module and dependencies
├── Dockerfile           # Multi-stage build for the API
├── docker-compose.yml   # Run the API in Docker
├── handlers/            # HTTP handlers (config + config group)
├── model/               # Data types and repository interfaces
├── services/            # Business logic
└── repositories/        # Storage (in-memory; Consul code present but unused)
```

- **handlers**: Parse HTTP, call services, return JSON.
- **services**: Implement create/get/delete and label filtering.
- **repositories**: Define how configs and groups are stored (currently in-memory).
- **model**: Config, ConfigGroup, GroupConfig, parameters, labels.

---

## Tech stack

- **Go 1.19**
- **Gorilla Mux** for routing
- **Docker** (Alpine) for deployment

---

## Note

Data is **not persisted**. Everything is stored in memory, so restarting the server (or the container) removes all configs and groups. The app seeds a sample config and group on startup so you can try the API immediately.
