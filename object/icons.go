package object

import (
	"context"
	"fmt"
	"strings"
	"time"
	"urfunavigator/index/logger"
	"urfunavigator/index/models"
)

const BuildingIconsPrefix = "building-icons/"

const svgContentType = "image/svg+xml"

func ValidateIconFileName(fileName string) error {
	if fileName == "" {
		return fmt.Errorf("icon filename is required")
	}

	if strings.Contains(fileName, "/") || strings.Contains(fileName, "\\") {
		return fmt.Errorf("icon filename must not contain path separators")
	}

	if !strings.HasSuffix(fileName, ".svg") {
		return fmt.Errorf("icon must be an svg file")
	}

	return nil
}

func buildingIconKey(fileName string) (string, error) {
	if err := ValidateIconFileName(fileName); err != nil {
		return "", err
	}

	return BuildingIconsPrefix + fileName, nil
}

func (s *MinIOS3) GetBuildingIconURL(fileName string) (models.BuildingIconURL, error) {
	key, err := buildingIconKey(fileName)
	if err != nil {
		return models.BuildingIconURL{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()

	presigned, err := s.Client.PresignedGetObject(ctx, s.BucketName, key, s.presignTTL, nil)
	if err != nil {
		return models.BuildingIconURL{}, fmt.Errorf("create presigned url: %w", err)
	}

	expiresAt := time.Now().Add(s.presignTTL)
	logger.Debug("building icon presigned url generated", "icon", fileName, "expires_at", expiresAt)

	return models.BuildingIconURL{
		URL:       presigned.String(),
		ExpiresAt: expiresAt,
	}, nil
}

func (s *MinIOS3) CheckBuildingIcon(fileName string) (bool, error) {
	key, err := buildingIconKey(fileName)
	if err != nil {
		return false, err
	}

	return s.objectExists(key)
}

func (s *MinIOS3) PutBuildingIcon(fileName string, fileData []byte) error {
	key, err := buildingIconKey(fileName)
	if err != nil {
		return err
	}

	return s.putObject(key, fileData, svgContentType)
}

func (s *MinIOS3) RemoveBuildingIcon(fileName string) error {
	key, err := buildingIconKey(fileName)
	if err != nil {
		return err
	}

	return s.removeObject(key)
}
