package models

import (
	"net/http"
	"time"
)

// User представляет пользователя системы.
type User struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
