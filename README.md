# Terminal Snake Game (Go)

Консольная Змейка, написанная на чистом Go с использованием стандартной библиотеки и ANSI escape-последовательностей для UI.

## Технические особенности

- **Асинхронный ввод:** Считывание клавиш вынесено в отдельную горутину. Передача команд управления в основной игровой цикл реализована через буферизированный канал `chan rune` и неблокирующий `select/default`.
- **Игровой движок:** Отрисовка поля через матрицу `[][]rune` и координатную сетку `Point{X, Y}`.
- **Цветной CLI:** ANSI-коды для динамической очистки экрана и раскраски элементов (зеленая змейка, красные яблоки, белый фон).
- **Кроссплатформенность:** Обработка особенностей очистки консоли под Windows (`cmd /c cls`) и POSIX-системы.

## Управление

**Стрелки** — сдвиг направления.

## Запуск в режиме разработки
```bash
go run main.go
```

## Кросс-компиляция Windows/MacOS/Linux:

**Windows(x64 Powershell):**
```powershell
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o bin/snake.exe main.go
```

**Linux (x64):**
```bash
GOOS=linux GOARCH=amd64 go build -o bin/snake_linux main.go
```

**macOS (x64):**
```bash
GOOS=darwin GOARCH=amd64 go build -o bin/snake_mac main.go
