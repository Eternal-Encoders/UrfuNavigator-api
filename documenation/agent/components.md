# Components

Domain objects the API stores and returns. Public clients should use the response shapes in `documenation/api/specification.md`. Admin clients read and write the storage documents described here.

## Building

Collection `buildings`. Model `models.Building`.

A campus building: display name, external `url`, map `latitude` / `longitude`, icon filename, linked floor ids, linked color-schema ids, and GPS entries.

GPS entry (`BuildingGps`): `centreAltitude` and `floorId`.

Public `GET /api/building` and `GET /api/buildings` expand floors to `{ id, displayableName }`, expand color schemas, and replace the icon filename with a presigned `{ name, url, expiresAt }`. If any icon cannot be loaded, the request fails with `404`.

Admin routes keep id arrays and the raw icon filename.

Deleting a building deletes its floors, those floors' rooms, services, and graph points, and its color schemas.

## Floor

Collection `floors`. Model `models.Floor`.

Belongs to one `buildingId`. Has `elevation`, pixel `width` and `height`, id lists `rooms`, `services`, `graph`, and optional floor GPS (`altitude`, linear mapping, correction forces).

Public `GET /api/floor?id=` returns one floor with rooms, services, and graph points inlined. Admin `GET /admin_api/floors/:id` returns the document with id lists only.

Deleting a floor deletes its rooms, services, and graph points.

## Room

Collection `rooms`. Model `models.Room`.

A drawn region on a floor: polymorphic `shape`, optional `children` shapes, optional `pointId` and `type` linking it to a graph point, optional `colorSchema`, and `isBorder` / `isFill` flags.

## Service

Collection `services`. Model `models.Service`.

A drawn campus service (not a graph node): `shape`, optional `colorSchema`, `isBorder`, `isFill`.

## Graph point

Collection `graph_points`. Model `models.GraphPoint`.

A navigation node: `buildingId`, `floorId`, plan coordinates `x` and `y`, undirected-style `links` (other point ids), `types`, localized `names`, weekly `time`, `description`, `info`, and `isPassFree`.

`types` values:

`corridor`, `auditorium`, `dinning`, `exit`, `stair`, `toilet-m`, `toilet-w`, `cafe`, `vending`, `coworking`, `atm`, `wardrobe`, `print`, `deanery`, `students`, `other`.

Cross-building routes require points of type `exit`.

`names` items: `{ name, translations: [{ language, value }] }`.

`time` is seven weekday objects (`monday` … `sunday`), each `{ start, end, isDayOff }` with integer `start` and `end`.

Public list, get, search, and path responses use `GraphPointResponse` (audit fields omitted). Admin CRUD returns the storage document.

## Color schema

Collection `buildingColorSchemes`. Model `models.BuildingColorSchema`.

Named palette attached to buildings. Fields are the map design tokens; see `design-tokens.md`.

## User

Collection `users`. Model `models.User`.

`login`, Argon2id `passwordHash` (never returned), and `role`: `reader`, `writer`, or `admine`.

`POST /api/login` accepts `writer` and `admine` only. `reader` receives `403`. The same two roles pass the `/admin_api` middleware.

Admin user responses are `UserAdminResponse`: audit fields, `login`, `role`. Passwords are write-only.

## Icon

Not a MongoDB document. SVG files in MinIO at `building-icons/{filename}.svg`.

Filename rules: non-empty, no `/` or `\`, suffix `.svg`.

Public and admin `GET .../icons/:icon` return a presigned URL. Upload and delete exist only on `/admin_api`. Upload JSON field `image` is a base64 string (`[]byte` in Go).

## Shape

Polymorphic JSON object. Discriminator field `type`.

| `type` | Fields |
| --- | --- |
| `point` | `x`, `y` |
| `rectangle` | `x`, `y`, `width`, `height` |
| `poly` | `x`, `y`, `points` (array of point shapes) |
| `container` | `x`, `y`, `width`, `height`, `alignX`, `alignY`, `children` |
| `text` | `x`, `y`, `alignX`, `alignY`, `text` |
| `icon` | `x`, `y`, `width`, `height`, `icon` |
| `door` | `wallId`, `length`, `offset` |

`alignX`: `LEFT`, `RIGHT`, `CENTER`. `alignY`: `TOP`, `BOTTOM`, `CENTER`.

Unknown `type` fails JSON decode.
