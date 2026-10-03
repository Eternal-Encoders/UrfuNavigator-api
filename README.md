# UrfuNavigator API

HTTP API for campus navigation: buildings, floor plans, named places, search, and indoor routes between graph points. An admin API edits the same data.

Go service (Fiber) with MongoDB for documents, MinIO for building-icon SVGs, JWT for admin access, and A* for routes. Search uses a MongoDB Atlas Search index named `point_search`.

## Documentation

| Document | Use |
| --- | --- |
| [documenation/api/specification.md](documenation/api/specification.md) | HTTP contract to hand to agents in other projects |
| [documenation/agent/project.md](documenation/agent/project.md) | Architecture, environment, and data stores for work in this repo |
| [documenation/agent/components.md](documenation/agent/components.md) | Domain objects and how they link |
| [documenation/agent/design-tokens.md](documenation/agent/design-tokens.md) | Building color schemas used by map clients |
| [documenation/agent/pages.md](documenation/agent/pages.md) | Client flows and how responses connect |
| `docs/swagger.json` | Generated OpenAPI, UI at `/docs` while the server is running |

## Routes

- `GET {DEFAULT_PATH}/` — health check, plain text `OK`
- `/api` — public reads and `POST /api/login`
- `/admin_api` — JWT, roles `writer` and `admine`

Public building and floor responses expand linked data and return presigned icon URLs. Admin responses return stored documents (id lists, icon filename, audit fields).

## Environment

Copy these into `.env` or export them. A variable may also be supplied as `{NAME}_FILE` pointing at a file whose contents are the value.

Required:

- `PORT` — listen address, for example `:5000`
- `CORS` — allowed origins, separated by ` | `
- `DEFAULT_PATH` — prefix for the health route only
- `DATABASE_URI` — MongoDB URI
- `DATABASE_COLLECTION` — database name
- `BUCKET_ENDPOINT` — MinIO endpoint
- `BUCKET_ACCESS_KEY` — MinIO access key
- `BUCKET_SECRET_KEY` — MinIO secret key
- `BUCKET_NAME` — bucket name
- `JWT_SECRET` — HMAC secret for admin tokens

Optional:

- `MODE` — default `DEV`
- `LOG_LEVEL` — default `info`
- `LOG_FORMAT` — default `text`
- `S3_REQUEST_TIMEOUT_SEC` — default `60`
- `S3_PRESIGN_TTL_SEC` — default `3600`
- `S3_MAX_RETRIES` — default `2`
- `JWT_EXPIRES_IN_HOURS` — default `24`

MongoDB collections: `users`, `buildings`, `buildingColorSchemes`, `floors`, `graph_points`, `rooms`, `services`. Icons are stored as `building-icons/{name}.svg`.

## Run

Development, with [Air](https://github.com/air-verse/air):

```bash
air -c .air.windows.conf
```

Or:

```bash
go run .
```

Regenerate OpenAPI after handler comment changes:

```bash
go generate
```
