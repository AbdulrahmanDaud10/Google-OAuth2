package api

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary-key;"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not bull"`
	Password  string    `gorm:"not null"`
	Role      string    `gorm:"type:varchar(20);default:'user';"`
	Photo     string    `gorm:"default:'default.png';"`
	Verified  bool      `gorm:"default:false;"`
	Provider  string    `gorm:"default:local;"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

// RegisterUserInput struct will be used to validate the request payload during the signup process
type RegisteredUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginUserInput struct will be used to validate the request payload during the sign-in process.
type LoginUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Role      string `json:"role,omitempty"`
	Provider  string `json:"provider,omitempty"`
	Photo     string `json:"photo,omitempty"`
	Verified  bool   `json:"verified,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// function BeforeCreate uses NewV4 to return random generated UUID before creating a user
func (user *User) BeforeCreate(*gorm.DB) error {
	user.ID = uuid.NewV4()

	return nil
}

// FilteredResponse function to omit sensitive fields from the data returned by GORM.
func FilteredResponse(user *User) UserResponse {
	return UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		Verified:  user.Verified,
		Photo:     user.Photo,
		Provider:  user.Provider,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
