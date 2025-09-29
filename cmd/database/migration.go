package database

import (
	model "goPromotion/pkg/model"

	"gorm.io/gorm"
)

func Migration(db *gorm.DB) error {

	modelList := []interface{}{
		&model.Category{},
		&model.Product{},
		&model.Payment{},
		&model.Condition{},
		&model.Promotion{},
		&model.Order{},
		&model.OrderDetail{},
	}

	for _, m := range modelList {
		if err := db.AutoMigrate(m); err != nil {
			return err
		}
	}

	return nil
}
