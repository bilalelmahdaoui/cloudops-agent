# CloudOps Agent

CloudOps Agent is a small cloud operations dashboard with an AI assistant. It can inspect cloud services, read logs, compare CPU usage, and run lifecycle actions through explicit tools.

The project uses an in-memory cloud adapter, so it runs locally without a cloud account or database.

## Screenshots

> **Screenshot placeholder:** main dashboard and cloud service list.

> **Screenshot placeholder:** agent running a service lifecycle action.

## Stack

### Backend

- Go and `net/http`
- OpenAI Responses API
- Official MCP Go SDK
- In-memory cloud adapter
- Standard Go tests

### Frontend

- React
- TypeScript
- Vite
- React Markdown with GitHub Flavored Markdown support
- Plain CSS

## Architecture

The backend follows a lightweight Ports and Adapters structure:

```text
HTTP handlers                 MCP server
      |                            |
      v                            v
Application use cases <---- Chat orchestration
      |                            |
      v                            v
Repository ports              LLM / tool ports
      ^                            ^
      |                            |
Fake cloud adapter       OpenAI and MCP adapters
```

The application layer does not depend on OpenAI, HTTP, or the fake cloud implementation. The chat use case only coordinates the provider-neutral LLM and tool-provider ports. MCP handlers reuse the same cloud-service use cases as the HTTP API, so the dashboard and agent share one in-memory state.

```text
cloudops-agent-backend/
├── cmd/api/                 Dependency wiring and server startup
└── internal/
    ├── adapters/
    │   ├── cloud/           In-memory cloud implementation
    │   ├── http/            HTTP handlers and CORS
    │   ├── llm/             OpenAI Responses API adapter
    │   └── mcp/             Local MCP server and client
    ├── application/         Cloud and chat use cases
    ├── config/              Environment configuration
    ├── domain/              Cloud-service model
    └── ports/               Repository, LLM, and tool contracts

cloudops-agent-frontend/
└── src/
    ├── components/          Chat and cloud-service UI
    ├── hooks/               Cloud-service state and polling
    ├── services/            Backend API clients
    └── types/               Frontend models
```

## Features

- List and inspect cloud services
- View status, CPU usage, and service logs
- Restart, stop, and start services from the dashboard
- Run the same operations through natural-language chat
- Resolve services by ID or partial name
- Keep short conversation context between messages
- Persist chat history in browser storage
- Refresh service state while lifecycle operations are running
- Render assistant Markdown safely

## Local setup

### Requirements

- Go 1.27 or compatible
- Node.js 22+
- npm
- An OpenAI API key

### 1. Configure and run the backend

```bash
cd cloudops-agent-backend
cp .env.example .env
```

Edit `.env`:

```dotenv
OPENAI_API_KEY=your_api_key
OPENAI_MODEL=gpt-5-mini
```

Install dependencies and start the API:

```bash
go mod download
go run ./cmd/api
```

The API runs at `http://localhost:8080`.

### 2. Run the frontend

In another terminal:

```bash
cd cloudops-agent-frontend
cp .env.example .env
npm install
npm run dev
```

Open `http://localhost:5173`.

Set the backend URL in `cloudops-agent-frontend/.env`:

```dotenv
VITE_API_URL=http://localhost:8080
```

## HTTP API

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `GET` | `/cloud-services` | List all services |
| `GET` | `/cloud-services/{id}` | Get one service |
| `POST` | `/cloud-services/{id}/restart` | Restart a service |
| `POST` | `/cloud-services/{id}/shutdown` | Stop a service |
| `POST` | `/cloud-services/{id}/start` | Start a service |
| `POST` | `/chat` | Send a message to the agent |

Chat request:

```json
{
  "message": "Restart the backend service",
  "history": []
}
```

Chat response:

```json
{
  "message": "The backend service was restarted successfully."
}
```

## Agent tools

The local MCP server exposes six tools:

- `get_cloud_service`
- `get_all_cloud_services`
- `find_cloud_services`
- `restart_cloud_service`
- `shutdown_cloud_service`
- `start_cloud_service`

Tool execution stays behind the MCP boundary. The OpenAI adapter does not access the cloud repository directly.

## Checks

Backend:

```bash
cd cloudops-agent-backend
gofmt -w $(find cmd internal -name '*.go')
go test ./...
go test -race ./...
```

Frontend:

```bash
cd cloudops-agent-frontend
npm run build
npm run lint
```

## Current limitations

- Cloud services are simulated and stored in memory.
- Service state resets when the backend restarts.
- The API has no authentication or persistence layer.
- CORS is configured for the local Vite development origin.
