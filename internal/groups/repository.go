package groups

import (
	"training-app/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

func (r *Repository) WithDB(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// NOTE: groups is a global lookup table (admin, training
// manager, coordinator, user) and has no client_id, so
// unlike other repositories there is no ForClient here.

func (r *Repository) GetAll() ([]models.Group, error) {

	var groups []models.Group

	err := r.db.
		Order("id").
		Find(&groups).Error

	return groups, err
}

func (r *Repository) GetByID(id int) (*models.Group, error) {

	var group models.Group

	err := r.db.
		First(&group, id).Error

	if err != nil {
		return nil, err
	}

	return &group, nil
}

// GetAllowed returns every group with an id greater than or
// equal to minGroupID, ordered by id. A user may only assign
// other users to groups at or below his own group.
func (r *Repository) GetAllowed(minGroupID int) ([]models.Group, error) {

	var groups []models.Group

	err := r.db.
		Where("id >= ?", minGroupID).
		Order("id").
		Find(&groups).Error

	return groups, err
}
