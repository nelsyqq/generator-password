package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func ShowHeader() {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("           ГЕНЕРАТОР СЛУЧАЙНЫХ ПАРОЛЕЙ            ")
	fmt.Println(strings.Repeat("=", 50))
}

func ShowLoading(seconds int) {
	for i := 0; i < seconds; i++ {
		fmt.Print(".")
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Println(" Готово!")
}

func ShowPassword(password string, length int, purpose string) {
	fmt.Print("\nГенерируем пароль")
	ShowLoading(3)

	fmt.Println("\n" + strings.Repeat("*", 50))
	fmt.Printf("ВАШ НОВЫЙ ПАРОЛЬ:\n")
	fmt.Printf("Для: %s\n", purpose)
	fmt.Printf("Пароль: %s\n", password)
	fmt.Printf("Длина: %d символов\n", length)
	fmt.Println(strings.Repeat("*", 50))
}

func ShowGoodbye() {
	fmt.Println("\n" + strings.Repeat("*", 50))
	fmt.Println("Спасибо за использование нашего генератора!")
	fmt.Println("До новых встреч!")
	fmt.Println(strings.Repeat("*", 50))
	fmt.Print("Нажмите Enter для выхода...")

	// Ждем нажатия Enter перед выходом
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}