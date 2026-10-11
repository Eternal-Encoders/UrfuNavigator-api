# UrfuNavigator API

Campus navigation backend for Ural Federal University maps. It serves building and floor geometry, named graph points, search, and indoor routes, and it exposes an authenticated admin API for editing that data.

Snapshot date: 2026-10-03.

## Stack

- Language: Go 1.25 (`module urfunavigator/index`)
- HTTP: Fiber v3
- Database: MongoDB (driver v2). One database name from `DATABASE_COLLECTION`; collections are fixed in `store/collections.go`
- Icons: MinIO (S3-compatible), objects under `building-icons/`
- Auth: JWT (HS256) after Argon2id password check
- Routing: A* in `geo/`, including cross-building routes through exit points
- Search: MongoDB Atlas Search index `point_search` on graph points
- OpenAPI: Swagger UI at `/docs`, spec generated into `docs/swagger.json`

Entry point is `main.go`. Routes are registered in `api/api.go`.

## Runtime layout

```
HTTP
  GET  {DEFAULT_PATH}/     health, plain text "OK"
  /api/...                 public read API and login
  /admin_api/...           JWT + writer or admine role
  /docs                    Swagger UI

MongoDB collections
  users, buildings, buildingColorSchemes, floors,
  graph_points, rooms, services

MinIO
  building-icons/{filename}.svg
```

`DEFAULT_PATH` prefixes only the health check. Public and admin routes are always `/api` and `/admin_api`.

## Environment

Required (process env, or the same name with `_FILE` pointing at a secret file):

| Variable | Role |
| --- | --- |
| `PORT` | Listen address passed to Fiber, for example `:5000` |
| `CORS` | Allowed origins, split by ` \| ` (space, pipe, space) |
| `DEFAULT_PATH` | Prefix for the health route only |
| `DATABASE_URI` | MongoDB URI |
| `DATABASE_COLLECTION` | Database name |
| `BUCKET_ENDPOINT` | MinIO endpoint |
| `BUCKET_ACCESS_KEY` | MinIO access key |
| `BUCKET_SECRET_KEY` | MinIO secret key |
| `BUCKET_NAME` | Bucket name |
| `JWT_SECRET` | HMAC secret for admin tokens |

Optional:

| Variable | Default |
| --- | --- |
| `MODE` | `DEV` |
| `LOG_LEVEL` | `info` |
| `LOG_FORMAT` | `text` |
| `S3_REQUEST_TIMEOUT_SEC` | `60` |
| `S3_PRESIGN_TTL_SEC` | `3600` |
| `S3_MAX_RETRIES` | `2` |
| `JWT_EXPIRES_IN_HOURS` | `24` |

`.env` is loaded when present. Missing required variables stop the process.

## Identifier and error conventions

- Every stored entity uses a MongoDB ObjectID. JSON fields named `id` (and foreign keys) are 24-character hex strings.
- Public JSON omits audit fields. Admin JSON for stored documents includes `createdAt`, `updatedAt` (BSON field `updateAt`), `author`, and `lastUpdatedBy`.
- Create handlers stamp `createdAt`, `updatedAt`, and `lastUpdatedBy` from the JWT user. The path id overwrites `id` on update.
- Success bodies are JSON. Error bodies are plain text, not a JSON error object.
- Missing documents on admin reads and writes return `404` with text `Not found`. Invalid ObjectIDs return `400`.

## Documentation map

| Path | Audience |
| --- | --- |
| `documenation/agent/` | Agents working in this repository: components, color tokens, client flows |
| `documenation/api/specification.md` | Agents in other projects that only need the HTTP contract |
| `docs/swagger.json` | Generated OpenAPI. Regenerate with `go generate` in the module root |
| `README.md` | How to run the service |

## Latest work

2026-10-03: wrote the agent notes in this folder, the portable API specification, and the README. No application code changed.
