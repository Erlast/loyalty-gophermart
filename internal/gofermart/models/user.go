package models

import (
	"net/http"
	"time"
)

// User представляет пользователя системы.
type User struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Password  string    `json:"password"`
	Login     string    `json:"login"`
	ID        int64     `json:"id"`
}

// Bind для User необходим для использования с chi/render.
func (c *User) Bind(r *http.Request) error {
	// Здесь можно добавить валидацию полей
	return nil
}

// Credentials представляет данные для входа пользователя.
type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Bind для Credentials необходим для использования с chi/render.
func (c *Credentials) Bind(r *http.Request) error {
	// Здесь можно добавить валидацию полей
	return nil
}
