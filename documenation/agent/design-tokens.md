# Design tokens

This service has no UI stylesheet. Map colors are data: `BuildingColorSchema` documents linked from buildings.

A public building response includes resolved schemas. Admin responses return schema ids on the building and full documents from `/admin_api/color-schemes`.

## Schema

```json
{
  "id": "ObjectID hex",
  "displayableName": "string",
  "accentColor": "string",
  "whiteColorTheme": {},
  "darkColorTheme": {}
}
```

Admin documents also include `createdAt`, `updatedAt`, `author`, `lastUpdatedBy`.

`accentColor` and every theme color are opaque strings (typically CSS colors). The API does not validate format.

## Theme fields

Both `whiteColorTheme` and `darkColorTheme` use the same object. Omitted pointers are absent in JSON.

| Token | Meaning for a map client |
| --- | --- |
| `buildingBorder` | Building outline |
| `buildingFill` | Building fill |
| `buildingBackground` | Area behind the building |
| `roomBorder` | Default room outline |
| `roomFill` | Default room fill |
| `roomText` | Default room label |
| `roomTypeBorder` | Outline by graph-point type |
| `roomTypeFill` | Fill by graph-point type |
| `roomTypeText` | Label color by graph-point type |

Type maps are keyed by `GraphPointType`:

`corridor`, `auditorium`, `dinning`, `exit`, `stair`, `toilet-m`, `toilet-w`, `cafe`, `vending`, `coworking`, `atm`, `wardrobe`, `print`, `deanery`, `students`, `other`.

Rooms and services may point at a schema with `colorSchema`. `isBorder` and `isFill` say whether that element draws a stroke and a fill. The API does not resolve a room's schema into concrete colors; the client joins `colorSchema` to `building.colorSchemes`.

## Where tokens appear

- Public building payload: `colorSchemes[]` with both themes inlined.
- Admin building payload: `colorSchemes` is an array of ids.
- Create or replace a schema: `POST` or `PUT /admin_api/color-schemes`.
- Create and attach in one call: `POST /admin_api/buildings/:id/color-schemes`.
- Detach and delete: `DELETE /admin_api/buildings/:id/color-schemes/:schemaId`.
