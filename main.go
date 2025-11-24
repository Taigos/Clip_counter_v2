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

func writePrefixAndNumber(pref string, num int) {
	mu.Lock()
	defer mu.Unlock()

	text := pref + strconv.Itoa(num)
	clipboard.Write(clipboard.FmtText, []byte(text))
	fmt.Printf("Значение в буфере%s\n", text)
}

func onCtrlV() {
	mu.Lock()
	p := prefix
	c := counter
	mu.Unlock()

	writePrefixAndNumber(p, c)

	//fmt.Printf("Буфер обновлён: %s%d\n", p, c)

	// Обновляем глобальный счётчик
	mu.Lock()
	counter = c + 1
	mu.Unlock()
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

func main() {
	initClipboard()

	// Инициализируем с нулевым значением
	writePrefixAndNumber(prefix, counter)
	fmt.Printf("Инициализировано: %s%d\n", prefix, counter)
	fmt.Println("--- Нажимай Ctrl + V или C в консоли для смены начального значения ---")
	fmt.Println("--- Нажимай Ctrl + Shift + Q для выхода ---")

	gohook.Register(gohook.KeyDown, []string{"q", "ctrl", "shift"}, func(e gohook.Event) {
		fmt.Println("\nExit...")
		gohook.End()
	})

	gohook.Register(gohook.KeyDown, []string{"ctrl", "v"}, func(e gohook.Event) {
		onCtrlV()
	})

	// Запускаем горутину для слушания ввода в консоли
	go func() {
		for {
			var input string
			fmt.Scanln(&input)

			// Если пользователь ввёл "c" или "C" в консоли
			if strings.ToLower(input) == "c" {
				inputPrefixFromConsole()
			}
		}
	}()

	s := gohook.Start()
	defer gohook.End()

	<-gohook.Process(s)
	fmt.Println("Program finished")
}
