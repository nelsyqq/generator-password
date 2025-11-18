package models

import "time"

// PasswordConfig - конфигурация для генерации пароля
type PasswordConfig struct {
	Length     int  `json:"length"`
	UseLower   bool `json:"use_lower"`
	UseUpper   bool `json:"use_upper"`
	UseDigits  bool `json:"use_digits"`
	UseSymbols bool `json:"use_symbols"`
}

// GeneratedPassword - сгенерированный пароль с метаданными
type GeneratedPassword struct {
	Password  string         `json:"password"`
	Purpose   string         `json:"purpose"` // Новое поле: для чего пароль
	Config    PasswordConfig `json:"config"`
	CreatedAt time.Time      `json:"created_at"`
}

// PasswordHistory - история сгенерированных паролей
type PasswordHistory struct {
	Passwords []GeneratedPassword `json:"passwords"`
}