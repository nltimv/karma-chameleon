package db

import "gorm.io/gorm"

func handleTransaction[T any](fn func(tx *gorm.DB) (T, error)) (T, error) {
	var result T
	err := db.Transaction(func(tx *gorm.DB) error {
		var err2 error
		result, err2 = fn(tx)
		return err2
	},
	)
	return result, err
}
