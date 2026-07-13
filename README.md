# Compute Server for the Woki Project

This module allows users to run GPU-intensive computational tasks on a remote server.

## Running the Project

Create the local backend environment file from the development template, then
fill in your local database, Redis, authentication, and mail credentials:

```bash
cp .env.development.example .env
```

Use the Makefile included in the project to start the backend:

```bash
make run
```

To start the frontend, run:

```bash
cd web
bun run dev
```

Vite automatically reads `web/.env.development` for local development and
`web/.env.production` for production builds. The production frontend calls
`https://mback.teckstate.com`.

Production backend settings are documented in `.env.production.example`.
Public settings and secrets are injected by Azure Container Apps; real
production secrets must never be committed to an environment file.
