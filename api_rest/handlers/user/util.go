package user

import (
	"github.com/Eugenill/SmartScooter/api_rest/models"
	_import00 "github.com/sqlbunny/sqlbunny/types/null"
	"time"
)

type ReqUser struct {
	ID           models.UserID `json:"id"`
	Username     string        `json:"username" `
	Secret       string        `json:"secret" `
	PhoneNumber  string        `json:"phone_number"`
	ContactEmail string        `json:"email" `
	Admin        bool          `json:"admin" `
}

type ReqUsernames struct {
	Usernames []string `json:"usernames"`
}
type RespUser struct {
	ID          models.UserID  `json:"id"`
	Username    string         `json:"usernames,omitempty"`
	Secret      string         `json:"secret"`
	PhoneNumber string         `json:"phone_number"`
	Email       string         `json:"email"`
	Admin       bool           `json:"admin"`
	CreatedAt   time.Time      `json:"created_at"`
	IsDeleted   bool           `json:"is_deleted"`
	DeletedAt   _import00.Time `json:"deleted_at"`
}
