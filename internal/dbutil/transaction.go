package dbutil

import "gorm.io/gorm"

// WithTransaction executes fn inside a database transaction.
//
// Usage:
//
//	err := dbutil.WithTransaction(s.db, func(tx *gorm.DB) error {
//
//		personRepo := s.personRepo.WithDB(tx).ForClient(clientID)
//		jobRepo := s.jobRepo.WithDB(tx).ForClient(clientID)
//
//		if err := personRepo.Create(&person); err != nil {
//			return err
//		}
//
//		if err := jobRepo.Create(&job); err != nil {
//			return err
//		}
//
//		return nil
//	})
//
// If fn returns an error, the transaction is rolled back.
// If fn returns nil, the transaction is committed.
func WithTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}