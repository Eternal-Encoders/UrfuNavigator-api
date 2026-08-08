package store

import (
	"fmt"

	"urfunavigator/index/models"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *MongoDB) deleteFloorContents(floor *models.Floor) error {
	if _, err := s.Rooms.DeleteMany(floor.Rooms); err != nil {
		return fmt.Errorf("delete floor rooms: %w", err)
	}
	if _, err := s.Services.DeleteMany(floor.Services); err != nil {
		return fmt.Errorf("delete floor services: %w", err)
	}
	if _, err := s.GraphPoints.DeleteMany(floor.Graph); err != nil {
		return fmt.Errorf("delete floor graph points: %w", err)
	}
	return nil
}

func (s *MongoDB) deleteFloorsContents(floors map[bson.ObjectID]*models.Floor) error {
	var roomIDs, serviceIDs, graphIDs []bson.ObjectID

	for _, floor := range floors {
		if floor == nil {
			continue
		}
		roomIDs = append(roomIDs, floor.Rooms...)
		serviceIDs = append(serviceIDs, floor.Services...)
		graphIDs = append(graphIDs, floor.Graph...)
	}

	if _, err := s.Rooms.DeleteMany(roomIDs); err != nil {
		return fmt.Errorf("delete floors rooms: %w", err)
	}
	if _, err := s.Services.DeleteMany(serviceIDs); err != nil {
		return fmt.Errorf("delete floors services: %w", err)
	}
	if _, err := s.GraphPoints.DeleteMany(graphIDs); err != nil {
		return fmt.Errorf("delete floors graph points: %w", err)
	}

	return nil
}

func (s *MongoDB) deleteBuildingContents(building *models.Building) error {
	if len(building.Floors) > 0 {
		floors, err := s.Floors.GetMany(building.Floors)
		if err != nil {
			return fmt.Errorf("load building floors: %w", err)
		}
		if err := s.deleteFloorsContents(floors); err != nil {
			return err
		}
		if _, err := s.Floors.DeleteMany(building.Floors); err != nil {
			return fmt.Errorf("delete building floors: %w", err)
		}
	}

	if _, err := s.BuildingColorSchemes.DeleteMany(building.ColorSchemes); err != nil {
		return fmt.Errorf("delete building color schemes: %w", err)
	}

	return nil
}
