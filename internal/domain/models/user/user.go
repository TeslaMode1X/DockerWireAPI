package user

import (
	"github.com/gofrs/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"updatedAt,omitempty"`
}

type Entity struct {
	ID       string `json:"ID"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
