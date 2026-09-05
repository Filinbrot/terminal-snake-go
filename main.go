package main

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/eiannone/keyboard"
)

type Point struct {
	X int
	Y int
}

// ANSI Цвета:
// \033[30m - Черный текст
// \033[31m - Красный текст
// \033[32m - Зеленый текст
// \033[47m - Белый фон
// \033[0m  - Сброс стилей

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		// Очистка экрана и заливка фона белым
		fmt.Print("\033[2J\033[47m")
	}
}

func genMap(mapka [][]rune, snake []Point, apples []Point) {
	for i := range mapka {
		for j := range mapka[i] {
			mapka[i][j] = 'E'
		}
	}

	for _, a := range apples {
		mapka[a.Y][a.X] = 'A'
	}

	for _, s := range snake {
		mapka[s.Y][s.X] = 'S'
	}
}

func genApples(mapka [][]rune, snake []Point) []Point {
	size := len(mapka)
	maxApples := int(float64(size) * 0.2)
	if maxApples < 1 {
		maxApples = 1
	}

	count := rand.IntN(maxApples) + 1
	var apples []Point

	for i := 0; i < count; i++ {
		for {
			x := rand.IntN(size)
			y := rand.IntN(size)

			inSnake := false
			for _, s := range snake {
				if s.X == x && s.Y == y {
					inSnake = true
					break
				}
			}

			if !inSnake {
				apples = append(apples, Point{X: x, Y: y})
				break
			}
		}
	}

	return apples
}

func output(mapka [][]rune) {
	fmt.Print("\033[H")
	for i := range mapka {
		// Черные рамки на белом фоне
		line := "\033[30;47m|\033[0m"
		for j := range mapka[i] {
			switch mapka[i][j] {
			case 'E':
				// Пустая клетка (белый фон)
				line += "\033[47m  \033[0m"
			case 'A':
				// Красное яблоко на белом фоне
				line += "\033[31;47m○ \033[0m"
			case 'S':
				// Зеленая змейка на белом фоне
				line += "\033[32;47m● \033[0m"
			}
		}
		// Черные рамки на белом фоне + сброс цвета в конце
		line += "\033[30;47m|\033[0m"
		fmt.Println(line)
	}
}

func control(direction chan rune) {
	if err := keyboard.Open(); err != nil {
		return
	}
	defer keyboard.Close()

	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			return
		}

		var nextDir rune

		switch {
		case char == 'w' || char == 'W' || key == keyboard.KeyArrowUp:
			nextDir = 'N'
		case char == 's' || char == 'S' || key == keyboard.KeyArrowDown:
			nextDir = 'S'
		case char == 'a' || char == 'A' || key == keyboard.KeyArrowLeft:
			nextDir = 'W'
		case char == 'd' || char == 'D' || key == keyboard.KeyArrowRight:
			nextDir = 'E'
		default:
			continue
		}

		select {
		case direction <- nextDir:
		default:
			<-direction
			direction <- nextDir
		}
	}
}

func main() {
	direction := make(chan rune, 1)
	direction <- 'N'

	var size int
	fmt.Print("Змейка alpha 0.7\n\nВведите размер карты: ")

	for {
		fmt.Scan(&size)
		if size >= 10 {
			break
		}
		fmt.Printf("Вы ввели %d, что меньше 10-ти. Введите >=10: ", size)
	}

	snake := []Point{
		{X: size / 2, Y: size / 2},
		{X: size / 2, Y: (size / 2) + 1},
		{X: size / 2, Y: (size / 2) + 2},
	}

	mapka := make([][]rune, size)
	for i := range mapka {
		mapka[i] = make([]rune, size)
	}

	apples := genApples(mapka, snake)
	genMap(mapka, snake, apples)

	clearScreen()
	output(mapka)

	go control(direction)

	dir := 'N'

	for {
		time.Sleep(200 * time.Millisecond)

		select {
		case newDir := <-direction:
			if (newDir == 'N' && dir != 'S') ||
				(newDir == 'S' && dir != 'N') ||
				(newDir == 'W' && dir != 'E') ||
				(newDir == 'E' && dir != 'W') {
				dir = newDir
			}
		default:
		}
		tail := snake[len(snake)-1]

		for i := len(snake) - 1; i > 0; i-- {
			snake[i] = snake[i-1]
		}

		switch dir {
		case 'N':
			snake[0].Y--
		case 'S':
			snake[0].Y++
		case 'W':
			snake[0].X--
		case 'E':
			snake[0].X++
		}

		if snake[0].X >= size || snake[0].X < 0 || snake[0].Y >= size || snake[0].Y < 0 || mapka[snake[0].Y][snake[0].X] == 'S' {
			// Сброс цвета терминала при завершении игры
			fmt.Print("\033[0m\033[2J\033[H")
			fmt.Printf("КОНЕЦ ИГРЫ\nЗмейка врезалась!\nДлина змейки: %d\n", len(snake))
			time.Sleep(5 * time.Second)
			return
		}

		for i, a := range apples {
			if snake[0].X == a.X && snake[0].Y == a.Y {
				snake = append(snake, tail)
				apples = append(apples[:i], apples[i+1:]...)
				break
			}
		}

		if len(apples) == 0 {
			apples = genApples(mapka, snake)
		}

		genMap(mapka, snake, apples)
		output(mapka)
	}
}
