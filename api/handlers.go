package api

import (
	"container/heap"
	"log"
	"strings"
	"urfunavigator/index/models"
	"urfunavigator/index/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *API) MainHandler(c *fiber.Ctx) error {
	return c.SendString("OK")
}

func (s *API) FloorHandler(c *fiber.Ctx) error {
	_, floorExist := c.Queries()["floor"]
	institute, instituteExist := c.Queries()["institute"]

	floor := c.QueryInt("floor")

	if !floorExist || !instituteExist {
		log.Println("Request Floor without floor or institute")
		return c.Status(fiber.StatusBadRequest).SendString("Request must contain floor and institute query parameters")
	}

	floorData, err := s.Store.GetFloor(floor, institute)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetFloor")
	}

	response := models.FloorResponse{
		Institute: floorData.Institute,
		Floor:     floorData.Floor,
		Width:     floorData.Width,
		Height:    floorData.Height,
		Audiences: floorData.Audiences,
		Service:   floorData.Service,
	}

	return c.JSON(response)
}

func (s *API) InstituteHandler(c *fiber.Ctx) error {
	url, urlExist := c.Queries()["url"]

	if !urlExist {
		log.Println("Request Institute without url")
		return c.Status(fiber.StatusBadRequest).SendString("Request must contain url query parameters")
	}

	instituteData, err := s.Store.GetInstitute(url)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetInstitute")
	}

	iconData, iconErr := s.Store.GetInstituteIcons([]string{instituteData.Icon})
	if iconErr != nil {
		log.Println(iconErr)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetInstituteIcons")
	}
	if len(iconData) != 1 {
		log.Println("There is too many or no media with id")
		log.Println(iconData)
		return c.Status(fiber.StatusNotFound).SendString("Cannot find media by id")
	}

	response := models.InstituteResponse{
		Name:            instituteData.Name,
		DisplayableName: instituteData.DisplayableName,
		MinFloor:        instituteData.MinFloor,
		MaxFloor:        instituteData.MaxFloor,
		Url:             instituteData.Url,
		Latitude:        instituteData.Latitude,
		Longitude:       instituteData.Longitude,
		Icon:            iconData[0],
	}

	return c.JSON(response)
}

func (s *API) InstitutesHandler(c *fiber.Ctx) error {

	institutesData, err := s.Store.GetInstitutes()
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetInstitute")
	}

	iconIds := []string{}
	for _, institute := range institutesData {
		iconIds = append(iconIds, institute.Icon)
	}

	iconsData, iconsErr := s.Store.GetInstituteIcons(iconIds)
	if iconsErr != nil {
		log.Println(iconsErr)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetInstituteIcons")
	}
	if len(iconsData) != len(institutesData) {
		log.Printf("IconsData length = %d and InstitutesData length = %d", len(iconsData), len(institutesData))
		return c.Status(fiber.StatusNotFound).SendString("For some of the institutes icons not founded")
	}

	response := []models.InstituteResponse{}
	for i, institue := range institutesData {
		response = append(response, models.InstituteResponse{
			Name:            institue.Name,
			DisplayableName: institue.DisplayableName,
			MinFloor:        institue.MinFloor,
			MaxFloor:        institue.MaxFloor,
			Url:             institue.Url,
			Latitude:        institue.Latitude,
			Longitude:       institue.Longitude,
			Icon:            iconsData[i],
		})
	}

	return c.JSON(response)
}

func (s *API) PointsHandler(c *fiber.Ctx) error {
	queries := c.Queries()

	typeParam, typeExist := queries["type"]
	instituteParam, instituteExist := queries["institute"]
	_, floorExist := queries["floor"]
	nameParam, nameExist := queries["name"]
	_, lengthExist := queries["length"]

	floorParam := c.QueryInt("floor")
	lengthParam := c.QueryInt("length")

	pointsFilter := []models.PointsFilters{
		models.GetPointsFilter("types", bson.M{"$in": []string{typeParam}}, typeExist),
		models.GetPointsFilter("institute", instituteParam, instituteExist),
		models.GetPointsFilter("floor", floorParam, floorExist),
		models.GetPointsFilter("names", bson.M{
			"$in": []primitive.Regex{
				{
					Pattern: nameParam,
					Options: "i",
				},
			},
		}, nameExist),
	}

	limit := lengthParam
	if !lengthExist {
		limit = 40
	}
	points, err := s.Store.GetPoints(pointsFilter, limit)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetPoints")
	}

	return c.JSON(points)
}

func (s *API) PointIdHandler(c *fiber.Ctx) error {
	id, idExist := c.Queries()["id"]
	if !idExist {
		log.Println("Request Point by id without floor or institute")
		return c.Status(fiber.StatusBadRequest).SendString("Request must contain id query parameters")
	}

	point, err := s.Store.GetPoint(id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetPoint")
	}
	return c.JSON(point)
}

func (s *API) PathHandler(c *fiber.Ctx) error {
	from, fromExist := c.Queries()["from"]
	to, toExist := c.Queries()["to"]
	if !fromExist || !toExist {
		log.Println("Request Path without from or to")
		return c.Status(fiber.StatusBadRequest).SendString("Request must contain from and to query parameters")
	}

	toVisit := models.PriorityQueue{}
	paths := make(map[string]models.GraphPoint)
	costs := make(map[string]float64)

	heap.Init(&toVisit)

	start, startErr := s.Store.GetPoint(from)
	end, endErr := s.Store.GetPoint(to)
	if startErr != nil {
		log.Println(startErr)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in fetching from point")
	}
	if endErr != nil {
		log.Println(endErr)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in fetching to point")
	}
	if start.Institute != end.Institute {
		log.Println("Attempt to find path between point in different institutes")
		return c.Status(fiber.StatusInternalServerError).SendString("Both point must belong to same institute")
	}

	floorGraph, floorGraphErr := s.Store.GetGraph(start.Floor, start.Institute)
	stairs, stairErr := s.Store.GetStairs(start.Institute)
	if floorGraphErr != nil {
		log.Println(floorGraphErr)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in fetching graph")
	}
	if stairErr != nil {
		log.Println(floorGraphErr)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in fetching stairs")
	}

	graph := map[int]map[string]models.GraphPoint{
		start.Floor: floorGraph,
	}

	heap.Push(&toVisit, &models.QueueItem{
		Value: models.QueueItemValue{
			Id:    start.Id,
			Floor: start.Floor,
		},
		Priority: 0,
	})
	costs[start.Id] = 0

	for len(toVisit) != 0 {
		current := heap.Pop(&toVisit).(*models.QueueItem)
		if current.Value.Id == end.Id {
			break
		}

		currentPoint := graph[current.Value.Floor][current.Value.Id]
		neighbours := utils.GetNeighbours(currentPoint, stairs)
		for _, next := range neighbours {
			nextGraphFloor, nextGraphFloorExist := graph[next.Value.Floor]
			if !nextGraphFloorExist {
				graphFloor, floorGraphErr := s.Store.GetGraph(next.Value.Floor, start.Institute)
				if floorGraphErr != nil {
					log.Println(floorGraphErr)
					return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in fetching graph")
				}
				graph[next.Value.Floor] = graphFloor
				nextGraphFloor = graphFloor
			}

			nextPoint := nextGraphFloor[next.Value.Id]
			newCost := costs[currentPoint.Id] + utils.GraphCost(currentPoint, nextPoint)

			lastCost, isVisited := costs[nextPoint.Id]
			if !isVisited || newCost < lastCost {
				costs[nextPoint.Id] = newCost
				next.Priority = int(newCost) + int(utils.Heuristic(end, nextPoint))
				heap.Push(&toVisit, &next)
				paths[nextPoint.Id] = currentPoint
			}
		}
	}

	return c.JSON(utils.RestorePath(paths, start, end))
}

func (s *API) ObjectHandler(c *fiber.Ctx) error {
	iconName := c.Params("icon")
	if !strings.HasSuffix(iconName, ".svg") {
		log.Println("Request Object with unsupported type")
		return c.Status(fiber.StatusBadRequest).SendString("This file type is not supported")
	}

	obj, err := s.ObjectStore.GetFile(iconName)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Cannot get file from Object Storage")
	}

	c.Attachment(iconName)
	return c.Send(obj)
}
