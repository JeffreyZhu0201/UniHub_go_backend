package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"unihub/internal/model"
	"unihub/internal/repo"
)

type OpenService interface {
	RegisterDeveloper(name, email string) (*model.Developer, error)
	CreateApp(devSecret, appName string) (*model.App, error)
}

type openService struct {
	openRepo repo.OpenRepository
}

func NewOpenService(openRepo repo.OpenRepository) OpenService {
	return &openService{openRepo: openRepo}
}

func (s *openService) RegisterDeveloper(name, email string) (*model.Developer, error) {
	// Generate Secret
	bytes := make([]byte, 32)
	rand.Read(bytes)
	secret := hex.EncodeToString(bytes)

	dev := model.Developer{
		Name:   name,
		Email:  email,
		Secret: secret,
	}

	if err := s.openRepo.CreateDeveloper(&dev); err != nil {
		return nil, errors.New("邮箱已被注册或系统错误")
	}

	return &dev, nil
}

func (s *openService) CreateApp(devSecret, appName string) (*model.App, error) {
	dev, err := s.openRepo.GetDeveloperBySecret(devSecret)
	if err != nil {
		return nil, errors.New("无效的开发者密钥")
	}

	// Generate AppID, AppSecret
	appIDData := make([]byte, 8)
	rand.Read(appIDData)
	appID := hex.EncodeToString(appIDData)

	appSecretData := make([]byte, 16)
	rand.Read(appSecretData)
	appSecret := hex.EncodeToString(appSecretData)

	app := model.App{
		DeveloperID: dev.ID,
		Name:        appName,
		AppID:       appID,
		AppSecret:   appSecret,
		RateLimit:   60, // Default 60 req/min
	}

	if err := s.openRepo.CreateApp(&app); err != nil {
		return nil, err
	}

	return &app, nil
}
