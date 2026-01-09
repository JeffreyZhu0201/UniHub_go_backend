package repo

import (
	"unihub/internal/model"

	"gorm.io/gorm"
)

type OpenRepository interface {
	CreateDeveloper(dev *model.Developer) error
	GetDeveloperBySecret(secret string) (*model.Developer, error)
	CreateApp(app *model.App) error
}

type openRepository struct {
	db *gorm.DB
}

func NewOpenRepository(db *gorm.DB) OpenRepository {
	return &openRepository{db: db}
}

func (r *openRepository) CreateDeveloper(dev *model.Developer) error {
	return r.db.Create(dev).Error
}

func (r *openRepository) GetDeveloperBySecret(secret string) (*model.Developer, error) {
	var dev model.Developer
	if err := r.db.Where("secret = ?", secret).First(&dev).Error; err != nil {
		return nil, err
	}
	return &dev, nil
}

func (r *openRepository) CreateApp(app *model.App) error {
	return r.db.Create(app).Error
}
