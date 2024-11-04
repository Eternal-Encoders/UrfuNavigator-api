package api

import (
	"log"
	"strings"
	"urfunavigator/index/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *API) MainHandler(c *fiber.Ctx) error {
	defer s.Logger.WriteLog(c.IP(), "MainHandler", c.AllParams())
	return c.SendString("OK")
}

func (s *API) FloorHandler(c *fiber.Ctx) error {
	defer s.Logger.WriteLog(c.IP(), "FloorHandler", c.Queries())

	var query FloorQuery
	if err := c.QueryParser(&query); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Request must contain floor and institute query parameters")
	}

	floorData, err := s.Store.GetFloor(query.Floor, query.Institute)
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
	defer s.Logger.WriteLog(c.IP(), "InstituteHandler", c.Queries())

	var query InstituteQuery
	if err := c.QueryParser(&query); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Request must contain url query parameters")
	}

	instituteData, err := s.Store.GetInstitute(query.Institute)
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
	defer s.Logger.WriteLog(c.IP(), "InstitutesHandler", c.AllParams())

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
	defer s.Logger.WriteLog(c.IP(), "PointsHandler", c.Queries())

	var query PointsQuery
	if err := c.QueryParser(&query); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Something went wrong in QueryParser")
	}

	pointsFilter := []models.PointsFilters{}
	if query.Type != nil {
		pointsFilter = append(pointsFilter, models.GetPointsFilter("types", bson.M{"$in": []string{*query.Type}}))
	}
	if query.Institute != nil {
		pointsFilter = append(pointsFilter, models.GetPointsFilter("institute", *query.Institute))
	}
	if query.Floor != nil {
		pointsFilter = append(pointsFilter, models.GetPointsFilter("floor", *query.Floor))
	}
	if query.Name != nil {
		pointsFilter = append(pointsFilter, models.GetPointsFilter("names", bson.M{
			"$in": []primitive.Regex{
				{
					Pattern: *query.Name,
					Options: "i",
				},
			},
		}))
	}

	limit := 40
	if query.Length != nil {
		limit = *query.Length
	}
	points, err := s.Store.GetPoints(pointsFilter, limit)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetPoints")
	}

	return c.JSON(points)
}

func (s *API) PointIdHandler(c *fiber.Ctx) error {
	defer s.Logger.WriteLog(c.IP(), "PointIdHandler", c.Queries())

	var query PointIdQuery
	if err := c.QueryParser(&query); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Request must contain id query parameters")
	}

	point, err := s.Store.GetPoint(query.Id)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in GetPoint")
	}
	return c.JSON(point)
}

func (s *API) PathHandler(c *fiber.Ctx) error {
	defer s.Logger.WriteLog(c.IP(), "PathHandler", c.Queries())

	var query PathQuery
	if err := c.QueryParser(&query); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Request must contain from and to query parameters")
	}

	start, startErr := s.Store.GetPoint(query.From)
	end, endErr := s.Store.GetPoint(query.To)
	if startErr != nil {
		log.Println(startErr)
		return c.Status(fiber.StatusInternalServerError).SendString(startErr.Error())
	}
	if endErr != nil {
		log.Println(endErr)
		return c.Status(fiber.StatusInternalServerError).SendString(endErr.Error())
	}

	path, pathErr := s.GEoService.FindPath(
		start,
		end,
		s.Store.GetGraph,
		s.Store.GetStairs,
		s.Store.GetEnters,
	)
	if pathErr != nil {
		log.Println(pathErr)
		return c.Status(fiber.StatusInternalServerError).SendString(pathErr.Error())
	}

	return c.JSON(map[string]any{
		"result": path,
	})
}

func (s *API) ObjectHandler(c *fiber.Ctx) error {
	defer s.Logger.WriteLog(c.IP(), "ObjectHandler", c.AllParams())

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

func (s *API) SearchHandler(c *fiber.Ctx) error {
	defer s.Logger.WriteLog(c.IP(), "SearchHandler", c.Queries())

	var query SearchQuery
	if err := c.QueryParser(&query); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Request must contain name query parameters")
	}
	length := 40

	if query.Length != nil {
		rawLength := *query.Length
		if rawLength > 40 {
			rawLength = 40
		}
		if rawLength < 1 {
			rawLength = 1
		}
		length = rawLength
	}

	points, err := s.Store.GetBySearchEngine(query.Name, length)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).SendString("Something went wrong in Search")
	}

	return c.JSON(points)
}
