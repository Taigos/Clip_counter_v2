package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	gohook "github.com/robotn/gohook"
	"golang.design/x/clipboard"
)

var mu sync.Mutex

var (
	prefix  string = "" // Изначально пусто
	counter int    = 1  // Начинаем с 1
)

func initClipboard() {
	if err := clipboard.Init(); err != nil {
		log.Fatalf("Ошибка инициализации буфера: %v", err)
	}
}

func clearClipboard() {
	clipboard.Write(clipboard.FmtText, []byte(""))
}

func onCtrlV2() {
	mu.Lock()
	defer mu.Unlock()

	text := prefix + strconv.Itoa(counter)
	fmt.Printf("Отправили в буфер: %s\n", text)
	clipboard.Write(clipboard.FmtText, []byte(text))

	fmt.Printf("Прочитали из буфера: %s\n", clipboard.Read(clipboard.FmtText))
	counter++
}

// Функция для ввода префикса в консоли
func inputPrefixFromConsole() {
	fmt.Print("\nВведите новый префикс (или пусто для удаления): ")
	var initValue string
	var newPrefix string
	var newCounter int
	fmt.Scanln(&initValue)

	if initValue == "" {
		newPrefix = ""
		newCounter = 1
	} else {
		// Найти где начинается число в конце строки
		i := len(initValue) - 1
		for i >= 0 && initValue[i] >= '0' && initValue[i] <= '9' {
			i--
		}

		// Если нет цифр в конце, вернуть текущие значения
		if i == len(initValue)-1 {
			newCounter = 1
			newPrefix = initValue
		} else {
			newPrefix = initValue[:i+1]
			numStr := initValue[i+1:]
			num, err := strconv.Atoi(numStr)
			if err != nil {
				newCounter = 1
			} else {
				newCounter = num
			}
		}
	}

	mu.Lock()
	prefix = newPrefix
	counter = newCounter
	mu.Unlock()

	fmt.Printf("Префикс: '%s'\nСчетчик: '%d'\n", prefix, counter)
}
func registerHooks() {
	gohook.Register(gohook.KeyDown, []string{"q", "ctrl", "shift"}, func(e gohook.Event) {
		fmt.Println("\nExit...")
		gohook.End()
	})

	gohook.Register(gohook.KeyDown, []string{"ctrl", "v"}, func(e gohook.Event) {
		onCtrlV2()
	})
}
func main() {
	initClipboard()
	//test()

	clearClipboard()
	fmt.Println("--- Нажимай Ctrl + V или C в консоли для смены начального значения ---")
	fmt.Println("--- Нажимай Ctrl + Shift + Q для выхода ---")

	registerHooks()

	// Запускаем горутину для слушания ввода в консоли
	go func() {
		for {
			var input string
			fmt.Scanln(&input)

			// Если пользователь ввёл "c" или "C" в консоли
			if strings.ToLower(input) == "c" {
				inputPrefixFromConsole()
				clearClipboard()
			}
		}
	}()

	s := gohook.Start()
	defer gohook.End()

	<-gohook.Process(s)
	fmt.Println("Program finished")
}

func test() {
	fmt.Printf("----------TEST----------\n")
	mu.Lock()
	defer mu.Unlock()

	old_text := clipboard.Read(clipboard.FmtText)
	fmt.Printf("Значение в буфере: %s\n", old_text)

	fmt.Printf("Пишем пустоту в буфер.\n")
	text := ""
	clipboard.Write(clipboard.FmtText, []byte(text))
	old_text = clipboard.Read(clipboard.FmtText)
	fmt.Printf("Значение в буфере: %s\n", old_text)

	fmt.Printf("Пишем 1 в буфер.\n")
	num := 1
	text = strconv.Itoa(num)
	clipboard.Write(clipboard.FmtText, []byte(text))
	old_text = clipboard.Read(clipboard.FmtText)
	fmt.Printf("Значение в буфере: %s\n", old_text)

	fmt.Printf("Увеличиваем на 1\n")
	num += 1
	text = strconv.Itoa(num)
	clipboard.Write(clipboard.FmtText, []byte(text))
	old_text = clipboard.Read(clipboard.FmtText)
	fmt.Printf("Значение в буфере: %s\n", old_text)

	fmt.Printf("Зписываем 3 значения подряд\n")
	text = "Первое"
	clipboard.Write(clipboard.FmtText, []byte(text))
	text = "Второе"
	clipboard.Write(clipboard.FmtText, []byte(text))
	text = "Третье"
	clipboard.Write(clipboard.FmtText, []byte(text))
	old_text = clipboard.Read(clipboard.FmtText)
	fmt.Printf("Последнее значение в буфере: %s\n", old_text)
	old_text = clipboard.Read(clipboard.FmtText)
	fmt.Printf("Последнее значение в буфере: %s\n", old_text)

	fmt.Printf("----------TEST----------\n")
}
