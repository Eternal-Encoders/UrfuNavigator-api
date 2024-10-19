package api

type FloorQuery struct {
	Floor     int    `query:"floor,required"`
	Institute string `query:"institute,required"`
}

type InstituteQuery struct {
	Institute string `query:"institute,required"`
}

type PointsQuery struct {
	Institute *string `query:"institute"`
	Floor     *int    `query:"floor"`
	Type      *string `query:"type"`
	Name      *string `query:"name"`
	Length    *int    `query:"length"`
}

type PointIdQuery struct {
	Id string `query:"id,required"`
}

type PathQuery struct {
	From string `query:"from,required"`
	To   string `query:"to,required"`
}

type SearchQuery struct {
	Name   string `query:"name,required"`
	Length *int   `query:"length"`
}
