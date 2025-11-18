package generator

import (
	"fmt"
	"password-generator/models"
	"password-generator/security"
	"password-generator/storage"
	"strings"
)

// PasswordGenerator - структура генератора паролей
type PasswordGenerator struct {
	lowercase string
	uppercase string
	digits    string
	symbols   string
}

// NewPasswordGenerator создает новый экземпляр генератора
func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{
		lowercase: "abcdefghijklmnopqrstuvwxyz",
		uppercase: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		digits:    "0123456789",
		symbols:   "!@#$%^&*()-_=+[]{}|;:,.<>?/",
	}
}

// GenerateWithOptions генерирует пароль и сохраняет его
func (pg *PasswordGenerator) GenerateWithOptions(config models.PasswordConfig, purpose string) (string, error) {
	// Если назначение не указано, используем значение по умолчанию
	if strings.TrimSpace(purpose) == "" {
		purpose = "Не указано"
	}

	allChars := ""
	typesUsed := 0

	if config.UseLower {
		allChars += pg.lowercase
		typesUsed++
	}
	if config.UseUpper {
		allChars += pg.uppercase
		typesUsed++
	}
	if config.UseDigits {
		allChars += pg.digits
		typesUsed++
	}
	if config.UseSymbols {
		allChars += pg.symbols
		typesUsed++
	}

	if allChars == "" {
		return "", fmt.Errorf("вы не выбрали ни один тип символов")
	}

	password := pg.guaranteeCharacterTypes(config)

	for i := typesUsed; i < config.Length; i++ {
		index := security.GetRandomIndex(len(allChars))
		password += string(allChars[index])
	}

	password = pg.shuffleString(password)

	// Сохраняем пароль в историю с указанием назначения
	err := storage.SavePassword(password, purpose, config)
	if err != nil {
		return "", fmt.Errorf("ошибка сохранения пароля: %v", err)
	}

	return password, nil
}

func (pg *PasswordGenerator) guaranteeCharacterTypes(config models.PasswordConfig) string {
	result := ""

	if config.UseLower {
		index := security.GetRandomIndex(len(pg.lowercase))
		result += string(pg.lowercase[index])
	}
	if config.UseUpper {
		index := security.GetRandomIndex(len(pg.uppercase))
		result += string(pg.uppercase[index])
	}
	if config.UseDigits {
		index := security.GetRandomIndex(len(pg.digits))
		result += string(pg.digits[index])
	}
	if config.UseSymbols {
		index := security.GetRandomIndex(len(pg.symbols))
		result += string(pg.symbols[index])
	}

	return result
}

func (pg *PasswordGenerator) shuffleString(s string) string {
	runes := []rune(s)
	for i := len(runes) - 1; i > 0; i-- {
		j := security.GetRandomIndex(i + 1)
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}