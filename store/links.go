package store

import (
	"fmt"

	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *MongoDB) AddBuildingFloor(buildingID bson.ObjectID, floor models.Floor) (*bson.ObjectID, error) {
	building, err := s.Buildings.GetOne(buildingID)
	if err != nil {
		return nil, fmt.Errorf("building not found: %w", err)
	}

	floor.BuildingId = buildingID
	id, err := s.Floors.InsertOne(floor)
	if err != nil {
		return nil, err
	}

	if err := s.Buildings.PushToArray(building.Id, "floors", *id); err != nil {
		return nil, fmt.Errorf("link floor to building: %w", err)
	}

	return id, nil
}

func (s *MongoDB) AddBuildingColorSchemaLink(buildingID bson.ObjectID, schema models.BuildingColorSchema) (*bson.ObjectID, error) {
	building, err := s.Buildings.GetOne(buildingID)
	if err != nil {
		return nil, fmt.Errorf("building not found: %w", err)
	}

	id, err := s.BuildingColorSchemes.InsertOne(schema)
	if err != nil {
		return nil, err
	}

	if err := s.Buildings.PushToArray(building.Id, "colorSchemes", *id); err != nil {
		return nil, fmt.Errorf("link color schema to building: %w", err)
	}

	return id, nil
}

func (s *MongoDB) AddBuildingGps(buildingID bson.ObjectID, gps models.BuildingGps) error {
	if _, err := s.Buildings.GetOne(buildingID); err != nil {
		return fmt.Errorf("building not found: %w", err)
	}

	return s.Buildings.PushToArray(buildingID, "gps", gps)
}

func (s *MongoDB) RemoveBuildingFloorLink(buildingID, floorID bson.ObjectID) error {
	if _, err := s.Buildings.GetOne(buildingID); err != nil {
		return fmt.Errorf("building not found: %w", err)
	}

	floor, err := s.Floors.GetOne(floorID)
	if err != nil {
		return fmt.Errorf("floor not found: %w", err)
	}
	if floor.BuildingId != buildingID {
		return fmt.Errorf("floor does not belong to building")
	}

	if err := s.Buildings.PullFromArray(buildingID, "floors", floorID); err != nil {
		return err
	}

	_, err = s.RemoveFloor(floorID)
	return err
}

func (s *MongoDB) RemoveBuildingColorSchemaLink(buildingID, schemaID bson.ObjectID) error {
	if _, err := s.Buildings.GetOne(buildingID); err != nil {
		return fmt.Errorf("building not found: %w", err)
	}

	if err := s.Buildings.PullFromArray(buildingID, "colorSchemes", schemaID); err != nil {
		return err
	}

	_, err := s.RemoveBuildingColorSchema(schemaID)
	return err
}

func (s *MongoDB) RemoveBuildingGps(buildingID bson.ObjectID, floorID bson.ObjectID) error {
	if _, err := s.Buildings.GetOne(buildingID); err != nil {
		return fmt.Errorf("building not found: %w", err)
	}

	return s.Buildings.PullFromArray(buildingID, "gps", bson.M{"floorId": floorID})
}

func (s *MongoDB) AddFloorRoom(floorID bson.ObjectID, room models.Room) (*bson.ObjectID, error) {
	floor, err := s.Floors.GetOne(floorID)
	if err != nil {
		return nil, fmt.Errorf("floor not found: %w", err)
	}

	id, err := s.Rooms.InsertOne(room)
	if err != nil {
		return nil, err
	}

	if err := s.Floors.PushToArray(floor.Id, "rooms", *id); err != nil {
		return nil, fmt.Errorf("link room to floor: %w", err)
	}

	return id, nil
}

func (s *MongoDB) AddFloorService(floorID bson.ObjectID, service models.Service) (*bson.ObjectID, error) {
	floor, err := s.Floors.GetOne(floorID)
	if err != nil {
		return nil, fmt.Errorf("floor not found: %w", err)
	}

	id, err := s.Services.InsertOne(service)
	if err != nil {
		return nil, err
	}

	if err := s.Floors.PushToArray(floor.Id, "services", *id); err != nil {
		return nil, fmt.Errorf("link service to floor: %w", err)
	}

	return id, nil
}

func (s *MongoDB) AddFloorGraphPoint(floorID bson.ObjectID, point models.GraphPoint) (*bson.ObjectID, error) {
	floor, err := s.Floors.GetOne(floorID)
	if err != nil {
		return nil, fmt.Errorf("floor not found: %w", err)
	}

	point.BuildingId = floor.BuildingId
	point.FloorId = floorID

	id, err := s.GraphPoints.InsertOne(point)
	if err != nil {
		return nil, err
	}

	if err := s.Floors.PushToArray(floor.Id, "graph", *id); err != nil {
		return nil, fmt.Errorf("link graph point to floor: %w", err)
	}

	return id, nil
}

func (s *MongoDB) RemoveFloorRoomLink(floorID, roomID bson.ObjectID) error {
	if _, err := s.Floors.GetOne(floorID); err != nil {
		return fmt.Errorf("floor not found: %w", err)
	}

	if err := s.Floors.PullFromArray(floorID, "rooms", roomID); err != nil {
		return err
	}

	_, err := s.RemoveRoom(roomID)
	return err
}

func (s *MongoDB) RemoveFloorServiceLink(floorID, serviceID bson.ObjectID) error {
	if _, err := s.Floors.GetOne(floorID); err != nil {
		return fmt.Errorf("floor not found: %w", err)
	}

	if err := s.Floors.PullFromArray(floorID, "services", serviceID); err != nil {
		return err
	}

	_, err := s.RemoveService(serviceID)
	return err
}

func (s *MongoDB) RemoveFloorGraphPointLink(floorID, pointID bson.ObjectID) error {
	if _, err := s.Floors.GetOne(floorID); err != nil {
		return fmt.Errorf("floor not found: %w", err)
	}

	if err := s.Floors.PullFromArray(floorID, "graph", pointID); err != nil {
		return err
	}

	_, err := s.RemoveGraphPoint(pointID)
	return err
}
