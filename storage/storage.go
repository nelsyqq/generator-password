package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"password-generator/models"
	"time"
)

const historyFile = "password_history.json"

// SavePassword сохраняет пароль в JSON файл
func SavePassword(password, purpose string, config models.PasswordConfig) error { // Добавлен параметр purpose
	// Читаем существующую историю
	history, err := readHistory()
	if err != nil {
		return err
	}

	// Создаем новую запись
	newPassword := models.GeneratedPassword{
		Password:  password,
		Purpose:   purpose, // Сохраняем назначение
		Config:    config,
		CreatedAt: time.Now(),
	}

	// Добавляем в историю
	history.Passwords = append(history.Passwords, newPassword)

	// Сохраняем обратно в файл
	return writeHistory(history)
}

// readHistory читает историю из JSON файла
func readHistory() (*models.PasswordHistory, error) {
	history := &models.PasswordHistory{}

	// Проверяем существует ли файл
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		// Файл не существует, возвращаем пустую историю
		return history, nil
	}

	// Читаем файл
	data, err := os.ReadFile(historyFile)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %v", err)
	}

	// Парсим JSON
	err = json.Unmarshal(data, history)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	return history, nil
}

// writeHistory записывает историю в JSON файл
func writeHistory(history *models.PasswordHistory) error {
	// Конвертируем в JSON с красивым форматированием
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка конвертации в JSON: %v", err)
	}

	// Записываем в файл
	err = os.WriteFile(historyFile, data, 0644)
	if err != nil {
		return fmt.Errorf("ошибка записи в файл: %v", err)
	}

	return nil
}

// GetHistory возвращает историю паролей
func GetHistory() (*models.PasswordHistory, error) {
	return readHistory()
}

// ClearHistory очищает историю паролей
func ClearHistory() error {
	history := &models.PasswordHistory{}
	return writeHistory(history)
}

// UpdatePasswordInHistory обновляет пароль в истории
func UpdatePasswordInHistory(updatedPassword *models.GeneratedPassword) error {
	history, err := readHistory()
	if err != nil {
		return err
	}

	// Ищем пароль для обновления по времени создания (уникальный идентификатор)
	for i, pwd := range history.Passwords {
		if pwd.CreatedAt.Equal(updatedPassword.CreatedAt) {
			history.Passwords[i] = *updatedPassword
			return writeHistory(history)
		}
	}

	return fmt.Errorf("пароль не найден в истории")
}

// DeletePasswordFromHistory удаляет пароль из истории
func DeletePasswordFromHistory(passwordToDelete *models.GeneratedPassword) error {
	history, err := readHistory()
	if err != nil {
		return err
	}

	// Ищем пароль для удаления по времени создания
	for i, pwd := range history.Passwords {
		if pwd.CreatedAt.Equal(passwordToDelete.CreatedAt) {
			// Удаляем пароль из среза
			history.Passwords = append(history.Passwords[:i], history.Passwords[i+1:]...)
			return writeHistory(history)
		}
	}

	return fmt.Errorf("пароль не найден в истории")
}