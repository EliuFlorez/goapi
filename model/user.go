package model

import (
	"goapi/database"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	// Columns
	FirstName string `gorm:"default:null"`
	LastName  string `gorm:"default:null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Phone     string `gorm:"default:null"`
	Password  string `gorm:"size:255;not null;" json:"-"`

	// Two Factor Auth
	SignInTwofa     bool   `gorm:"default:false"`
	TwofaCode       string `gorm:"default:null,index:idx_twofa_code"`
	TwofaCodeAt     time.Time
	TwofaCodeSentAt time.Time

	// Change email
	EmailToken  string `gorm:"default:null,index:idx_email_token"`
	EmailAt     time.Time
	EmailSentAt time.Time

	// Recoverable
	PasswordToken  string `gorm:"default:null,index:idx_password_token"`
	PasswordAt     time.Time
	PasswordSentAt time.Time

	// Trackable
	SignInCount     int
	CurrentSignInAt time.Time
	LastSignInAt    time.Time
	CurrentSignInIp string `gorm:"default:null"`
	LastSignInIp    string `gorm:"default:null"`

	// Confirmable
	ConfirmationToken  string `gorm:"default:null,index:idx_confirmation_token"`
	ConfirmationAt     time.Time
	ConfirmationSentAt time.Time
	ConfirmationEmail  bool `gorm:"default:false"`

	// Invintation
	InvitationToken  string `gorm:"default:null,index:idx_invitation_token"`
	InvitationAt     time.Time
	InvitationSentAt time.Time

	// Lockable
	FailedAttempts int
	LockedToken    string `gorm:"default:null,index:idx_locked_token"`
	LockedAt       time.Time
	LockedSentAt   time.Time

	// Relations
	Roles       []Role
	Permissions []Permission
	Accounts    []*Account `gorm:"many2many:user_accounts;references:ID;"`
}

type UserResponse struct {
	ID           uint      `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	SignInTwofa  bool      `json:"signin_twofa"`
	LastSignInAt time.Time `json:"last_signin_at"`

	// Relations
	Accounts []Account `json:"accounts"`
}

func (user *User) Save() (*User, error) {
	err := database.Sql.Debug().Create(&user).Error

	if err != nil {
		return &User{}, err
	}

	return user, nil
}

func (user *User) Update() error {
	err := database.Sql.Debug().Model(&user).Updates(user).Error

	if err != nil {
		return err
	}

	return nil
}

func FindUserByEmail(email string) (User, error) {
	var user User

	err := database.Sql.Debug().Where("email = ?", email).First(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func FindUserByCode(code string) (User, error) {
	var user User

	err := database.Sql.Debug().Where("twofa_code = ?", code).First(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func FindUserByColumn(column string, token string) (User, error) {
	var user User

	err := database.Sql.Debug().Where(column+" = ?", token).First(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (user *User) UpdateUserByToken(column string, token string) error {
	currentTime := time.Now()

	err := database.Sql.Debug().Model(&user).Updates(map[string]interface{}{
		column + "_token":   token,
		column + "_at":      currentTime,
		column + "_sent_at": currentTime.Add(time.Hour * 4),
	}).Error

	if err != nil {
		return err
	}

	return nil
}

func FindUserByToken(column string, token string) (User, error) {
	var user User

	err := database.Sql.Debug().Where(column+" = ?", token).First(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func FindUserById(id uint) (User, error) {
	var user User

	err := database.Sql.Debug().Where("ID = ?", id).First(&user).Error

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func CountUser() int {
	var user User

	result := database.Sql.Debug().Find(&user)

	if result.Error != nil {
		return 0
	}

	return int(result.RowsAffected)
}
