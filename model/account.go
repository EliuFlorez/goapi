package model

import (
	"goapi/database"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model

	// Columns
	Name          string `gorm:"default:null,index:idx_name"`
	Email         string `gorm:"default:null,index:idx_email"`
	Phone         string `gorm:"default:null"`
	Address       string `gorm:"default:null"`
	AddressNumber string `gorm:"default:null"`
	City          string `gorm:"default:null"`
	ProviceState  string `gorm:"default:null"`
	PortalCode    string `gorm:"default:null"`
	Country       string `gorm:"default:null"`

	// Relations
	Users []*User `gorm:"many2many:user_accounts;references:ID;"`
}

type UserAccount struct {
	gorm.Model

	// Columns
	UserID    uint `gorm:"primarykey"`
	AccountID uint `gorm:"primarykey"`
}

func AllAccounts(user *User, limit int) ([]Account, error) {
	var accounts []Account

	err := database.Sql.Limit(limit).Model(&user).Association("Accounts").Find(&accounts)

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (account *Account) All(user *User, limit int) ([]Account, error) {
	var accounts []Account

	err := database.Sql.Limit(limit).Model(&user).Association("Accounts").Find(&accounts)

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (account *Account) Total(user *User) int64 {
	return database.Sql.Model(&user).Association("Accounts").Count()
}

func (a *Account) Get(user *User, id int) (Account, error) {
	var account Account

	err := database.Sql.Where("ID = ?", id).Limit(1).Model(&user).Association("Accounts").Find(&account)

	if err != nil {
		return Account{}, err
	}

	return account, nil
}

func (user *User) AddAccountToUser(name string) error {
	newAccount := Account{
		Name: name,
		Users: []*User{
			user,
		},
	}

	err := database.Sql.Debug().Create(&newAccount).Error

	if err != nil {
		return err
	}

	return nil
}

func (account *Account) UpdateAccountByUser(user *User) error {
	err := database.Sql.Debug().Model(&account).Updates(account).Error

	if err != nil {
		return err
	}

	return nil
}

func (account *Account) DestroyAccountByUser(user *User) error {
	err := database.Sql.Debug().Model(&user).Association("Accounts").Delete(&account)

	if err != nil {
		return err
	}

	return nil
}
