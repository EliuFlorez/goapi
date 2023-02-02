package model

import (
	"goapi/database"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Crm struct {
	gorm.Model

	// Columns
	Name   string       `gorm:"default:null"`
	Entity string       `gorm:"default:null"`
	Oauth  pgtype.JSONB `gorm:"type:jsonb;default:'{}';not null"`
}

func (crm *Crm) Save() (*Crm, error) {
	err := database.Sql.Debug().FirstOrCreate(&crm).Error

	if err != nil {
		return &Crm{}, err
	}

	return crm, nil
}

func (crm *Crm) Update() error {
	err := database.Sql.Debug().Model(&crm).Updates(crm).Error

	if err != nil {
		return err
	}

	return nil
}

func (crm *Crm) UpdateAll(user *User) error {
	err := database.Sql.Debug().Model(&crm).Where("user_id = ?", user.ID).Updates(crm).Error

	if err != nil {
		return err
	}

	return nil
}

func FindCrmById(id int) (Crm, error) {
	var crm Crm

	err := database.Sql.Debug().Where("ID = ?", id).First(&crm).Error

	if err != nil {
		return Crm{}, err
	}

	return crm, nil
}
