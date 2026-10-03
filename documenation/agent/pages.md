# Pages and connections

There is no HTML UI in this repository. The surfaces below are the flows a map client or an admin client walks. Endpoint details live in `documenation/api/specification.md`.

## Public map

```
GET /api/buildings
        │  list of buildings, icons, color schemes, floor summaries
        ▼
GET /api/building?id=
        │  one building, same shape
        ▼
GET /api/floor?id=
        │  rooms, services, graph points for that floor
        ▼
GET /api/icons/:icon
           presigned SVG when the client needs the file again
```

Floor summaries on a building are `{ id, displayableName }`. The floor id is the key for `GET /api/floor`.

Room `pointId` and graph point `id` connect a drawn room to a navigation node. Graph point `links` connect nodes on the plan. `buildingId` and `floorId` on a point place that node on a building and floor.

## Search and a single point

```
GET /api/search?name=&length=
        │  Atlas Search, max 40
        ▼
GET /api/point?id=
        │  one graph point
        ▼
GET /api/points?buildingId=&floorId=&type=&name=&length=
           filtered list, default limit 40
```

`GET /api/points` name filter is a case-insensitive regex on `names.name`. `GET /api/search` uses the Atlas index `point_search` (autocomplete on `names`, text on description, info, types, displayableName) and only returns points that have at least one name.

## Route

```
chosen start point id ──┐
                        ├── GET /api/path?from=&to=
chosen end point id  ───┘
```

The result groups segments by building id, then floor id. Each segment is an ordered list of graph points. Same-building routes stay in one building. Different buildings route through `exit` points; missing exits fail the request.

A client draws each segment on the matching floor from `GET /api/floor`.

## Admin editor

```
POST /api/login
        │  token for role writer or admine
        ▼
Authorization: Bearer <token>
        │
        ├── /admin_api/buildings ── link floors, color schemas, GPS
        ├── /admin_api/floors ── link rooms, services, graph points
        ├── /admin_api/rooms | services | graph-points | color-schemes
        ├── /admin_api/users
        └── /admin_api/upload_icon | delete_icon | icons/:icon
```

Prefer the nested create routes when the new document must appear in the parent list:

- `POST /admin_api/buildings/:id/floors`
- `POST /admin_api/buildings/:id/color-schemes`
- `POST /admin_api/floors/:id/rooms`
- `POST /admin_api/floors/:id/services`
- `POST /admin_api/floors/:id/graph-points`

Standalone `POST /admin_api/floors` (and rooms, services, points, schemas) inserts a document and does not update the parent array.

Nested `DELETE` unlinks and deletes the child. Deleting a building or floor cascades to descendants. See `components.md`.

## Health

`GET {DEFAULT_PATH}/` returns plain text `OK`. It is not under `/api`.
