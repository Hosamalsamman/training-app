package countries

import "training-app/db"

type Repository struct {
}

func (r *Repository) GetAll() ([]Country, error) {
	var response []Country
	err := db.DB.Find(&response).Error
	return response, err
}

func (r Repository) GetByID(id int) (*Country, error) {
	var response Country
	err := db.DB.First(&response, id).Error
	return &response, err
}
