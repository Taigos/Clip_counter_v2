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

func initClipboard() {
	if err := clipboard.Init(); err != nil {
		log.Fatalf("Ошибка инициализации буфера: %v", err)
	}
}

func readNumberFromClipboard() int {
	mu.Lock()
	defer mu.Unlock()

	data := clipboard.Read(clipboard.FmtText)
	if data == nil {
		return 0
	}
	text := strings.TrimSpace(string(data))
	num, err := strconv.Atoi(text)
	if err != nil {
		return 0
	}
	return num
}

func writeNumberToClipboard(num int) {
	mu.Lock()
	defer mu.Unlock()

	clipboard.Write(clipboard.FmtText, []byte(strconv.Itoa(num)))
}

func onCtrlV() {
	num := readNumberFromClipboard()
	if num == 0 {
		num = 1
	} else {
		num++
	}
	writeNumberToClipboard(num)
	fmt.Printf("Буфер обновлён: %d\n", num)
}
func mainTwo() {
	fmt.Println("--- Please press ctrl + shift + q to stop gohook ---")
	gohook.Register(gohook.KeyDown, []string{"q", "ctrl", "shift"}, func(e gohook.Event) {
		fmt.Println("ctrl-shift-q")
		gohook.End()
	})

	gohook.Register(gohook.KeyDown, []string{"ctrl", "v"}, func(e gohook.Event) {
		onCtrlV()
	})

	s := gohook.Start()
	<-gohook.Process(s)

	evChan := gohook.Start()
	defer gohook.End()

	for ev := range evChan {
		fmt.Println("gohook: ", ev)
	}
}
func main() {
	mainTwo()
}
