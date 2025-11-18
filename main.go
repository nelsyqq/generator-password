package main

import (
	"bufio"
	"fmt"
	"os"
	"password-generator/generator"
	"password-generator/models"
	"password-generator/storage"
	"password-generator/utils"
	"strconv"
	"strings"
)

func main() {
	// Создаем экземпляр генератора
	passGen := generator.NewPasswordGenerator()
	reader := bufio.NewReader(os.Stdin)

	for {
		utils.ShowHeader()

		// Показываем меню
		fmt.Println("1. Сгенерировать новый пароль")
		fmt.Println("2. Показать историю паролей")
		fmt.Println("3. Очистить историю")
		fmt.Println("4. Выйти")
		fmt.Print("\nВыберите действие (1-4): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			generatePassword(passGen, reader)
		case "2":
			showPasswordHistoryMenu(reader) // Изменено на новую функцию меню
		case "3":
			clearPasswordHistory(reader)
		case "4":
			utils.ShowGoodbye()
			return
		default:
			fmt.Println("Неверный выбор! Пожалуйста, выберите от 1 до 4.")
			fmt.Print("Нажмите Enter чтобы продолжить...")
			reader.ReadString('\n')
		}
	}
}

// generatePassword генерирует новый пароль
func generatePassword(passGen *generator.PasswordGenerator, reader *bufio.Reader) {
	// Получаем настройки от пользователя
	config, purpose := getPasswordConfigFromUser(reader)

	// Генерируем пароль
	password, err := passGen.GenerateWithOptions(config, purpose)
	if err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		fmt.Print("Нажмите Enter чтобы продолжить...")
		reader.ReadString('\n')
		return
	}

	// Показываем результат
	utils.ShowPassword(password, config.Length, purpose)

	fmt.Print("Нажмите Enter чтобы продолжить...")
	reader.ReadString('\n')
}

// getPasswordConfigFromUser получает настройки от пользователя
func getPasswordConfigFromUser(reader *bufio.Reader) (models.PasswordConfig, string) {
	var config models.PasswordConfig
	var purpose string

	// Запрашиваем назначение пароля ПЕРВЫМ
	fmt.Print("Для чего этот пароль (например: Gmail, Facebook, банк): ")
	purpose, _ = reader.ReadString('\n')
	purpose = strings.TrimSpace(purpose)

	// Запрашиваем длину пароля
	config.Length = askForLength(reader)

	// Запрашиваем остальные настройки
	config.UseLower = askYesNo("Использовать строчные буквы (a-z)? (y/n): ", reader)
	config.UseUpper = askYesNo("Использовать заглавные буквы (A-Z)? (y/n): ", reader)
	config.UseDigits = askYesNo("Использовать цифры (0-9)? (y/n): ", reader)
	config.UseSymbols = askYesNo("Использовать специальные символы (!@#$% и т.д.)? (y/n): ", reader)

	return config, purpose
}

// askForLength запрашивает длину пароля с обработкой ввода
func askForLength(reader *bufio.Reader) int {
	for {
		fmt.Print("Введите длину пароля (4-50): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		length, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Ошибка: пожалуйста, введите только цифры!")
			continue
		}

		if length < 4 {
			fmt.Println("Ошибка: слишком короткий! Минимум 4 символа для безопасности.")
			continue
		}

		if length > 50 {
			fmt.Println("Ошибка: слишком длинный! Максимум 50 символов.")
			continue
		}

		return length
	}
}

// askYesNo задает вопрос да/нет с обработкой ввода
func askYesNo(prompt string, reader *bufio.Reader) bool {
	for {
		fmt.Print(prompt)
		answer, _ := reader.ReadString('\n')
		answer = strings.ToLower(strings.TrimSpace(answer))

		switch answer {
		case "y", "д", "yes", "да":
			return true
		case "n", "н", "no", "нет":
			return false
		default:
			fmt.Println("Ошибка: введите 'y' (да) или 'n' (нет)")
		}
	}
}

// showPasswordHistoryMenu показывает меню работы с историей паролей
func showPasswordHistoryMenu(reader *bufio.Reader) {
	history, err := storage.GetHistory()
	if err != nil {
		fmt.Printf("Ошибка при чтении истории: %v\n", err)
		fmt.Print("Нажмите Enter чтобы продолжить...")
		reader.ReadString('\n')
		return
	}

	if len(history.Passwords) == 0 {
		fmt.Println("\nИстория паролей пуста")
		fmt.Print("Нажмите Enter чтобы продолжить...")
		reader.ReadString('\n')
		return
	}

	for {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("              ИСТОРИЯ ПАРОЛЕЙ               ")
		fmt.Println(strings.Repeat("=", 50))

		// Показываем список паролей с номерами
		for i := len(history.Passwords) - 1; i >= 0; i-- {
			pwd := history.Passwords[i]
			displayIndex := len(history.Passwords) - i
			fmt.Printf("%d. Для: %s\n", displayIndex, pwd.Purpose)
			fmt.Printf("   Пароль: %s\n", pwd.Password)
			fmt.Printf("   Создан: %s\n", pwd.CreatedAt.Format("2006-01-02 15:04:05"))
			fmt.Println(strings.Repeat("-", 50))
		}

		// Меню действий
		fmt.Println("\nДействия:")
		fmt.Println("1. Выбрать пароль для редактирования")
		fmt.Println("2. Вернуться в главное меню")
		fmt.Print("Выберите действие (1-2): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			selectPasswordForEditing(history, reader)
			// Обновляем историю после возможных изменений
			history, err = storage.GetHistory()
			if err != nil {
				fmt.Printf("Ошибка при обновлении истории: %v\n", err)
				return
			}
		case "2":
			return
		default:
			fmt.Println("Неверный выбор! Пожалуйста, выберите 1 или 2.")
		}
	}
}

// selectPasswordForEditing позволяет выбрать пароль для редактирования
func selectPasswordForEditing(history *models.PasswordHistory, reader *bufio.Reader) {
	fmt.Print("\nВведите номер пароля для редактирования: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	index, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Ошибка: пожалуйста, введите только цифры!")
		return
	}

	// Конвертируем номер в индекс массива (показываем в обратном порядке)
	if index < 1 || index > len(history.Passwords) {
		fmt.Printf("Ошибка: номер должен быть от 1 до %d\n", len(history.Passwords))
		return
	}

	// Получаем реальный индекс (показываем последние первыми)
	realIndex := len(history.Passwords) - index
	selectedPassword := &history.Passwords[realIndex]

	editPasswordMenu(selectedPassword, reader)
}

// editPasswordMenu показывает меню редактирования пароля
func editPasswordMenu(password *models.GeneratedPassword, reader *bufio.Reader) {
	for {
		fmt.Println("\n" + strings.Repeat("=", 50))
		fmt.Println("             РЕДАКТИРОВАНИЕ ПАРОЛЯ            ")
		fmt.Println(strings.Repeat("=", 50))
		fmt.Printf("Для: %s\n", password.Purpose)
		fmt.Printf("Пароль: %s\n", password.Password)
		fmt.Printf("Длина: %d символов\n", password.Config.Length)
		fmt.Printf("Создан: %s\n", password.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println(strings.Repeat("-", 50))

		fmt.Println("\nЧто вы хотите изменить?")
		fmt.Println("1. Изменить назначение (для чего пароль)")
		fmt.Println("2. Удалить этот пароль из истории")
		fmt.Println("3. Вернуться к истории")
		fmt.Print("Выберите действие (1-3): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			editPasswordPurpose(password, reader)
		case "2":
			if deletePasswordFromHistory(password, reader) {
				return // Возвращаемся если пароль удален
			}
		case "3":
			return
		default:
			fmt.Println("Неверный выбор! Пожалуйста, выберите от 1 до 3.")
		}
	}
}

// editPasswordPurpose изменяет назначение пароля
func editPasswordPurpose(password *models.GeneratedPassword, reader *bufio.Reader) {
	fmt.Print("Введите новое назначение пароля: ")
	newPurpose, _ := reader.ReadString('\n')
	newPurpose = strings.TrimSpace(newPurpose)

	if newPurpose == "" {
		fmt.Println("Назначение не может быть пустым!")
		return
	}

	oldPurpose := password.Purpose
	password.Purpose = newPurpose

	// Сохраняем изменения
	err := storage.UpdatePasswordInHistory(password)
	if err != nil {
		fmt.Printf("Ошибка при сохранении изменений: %v\n", err)
		// Откатываем изменения в памяти
		password.Purpose = oldPurpose
	} else {
		fmt.Println("Назначение успешно изменено!")
	}
}

// deletePasswordFromHistory удаляет пароль из истории
func deletePasswordFromHistory(password *models.GeneratedPassword, reader *bufio.Reader) bool {
	if askYesNo("Вы уверены, что хотите удалить этот пароль из истории? (y/n): ", reader) {
		err := storage.DeletePasswordFromHistory(password)
		if err != nil {
			fmt.Printf("Ошибка при удалении пароля: %v\n", err)
			return false
		}
		fmt.Println("Пароль успешно удален из истории!")
		return true
	}
	return false
}

// clearPasswordHistory очищает историю паролей
func clearPasswordHistory(reader *bufio.Reader) {
	if askYesNo("Вы уверены, что хотите очистить всю историю паролей? (y/n): ", reader) {
		err := storage.ClearHistory()
		if err != nil {
			fmt.Printf("Ошибка при очистке истории: %v\n", err)
		} else {
			fmt.Println("История паролей успешно очищена!")
		}
	}
	fmt.Print("Нажмите Enter чтобы продолжить...")
	reader.ReadString('\n')
}