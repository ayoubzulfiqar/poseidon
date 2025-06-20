package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	width           = 40
	height          = 20
	initialSnakeLen = 3
	initialDelay    = 150 * time.Millisecond
	minDelay        = 50 * time.Millisecond
	delayDecrease   = 5 * time.Millisecond
)

type Point struct {
	X int
	Y int
}

type Snake struct {
	Body      []Point
	Direction Point
}

type Game struct {
	Screen    tcell.Screen
	Snake     Snake
	Food      Point
	Score     int
	GameOver  bool
	Delay     time.Duration
	InputChan chan tcell.Event
}

func NewGame(s tcell.Screen) *Game {
	rand.Seed(time.Now().UnixNano())

	snakeBody := make([]Point, initialSnakeLen)
	for i := 0; i < initialSnakeLen; i++ {
		snakeBody[i] = Point{X: width/2 - i, Y: height / 2}
	}

	game := &Game{
		Screen: s,
		Snake: Snake{
			Body:      snakeBody,
			Direction: Point{X: 1, Y: 0}, // Initial direction: right
		},
		Score:     0,
		GameOver:  false,
		Delay:     initialDelay,
		InputChan: make(chan tcell.Event),
	}

	game.generateFood()
	return game
}

func (g *Game) draw(p Point, char rune, style tcell.Style) {
	g.Screen.SetContent(p.X, p.Y, char, nil, style)
}

func (g *Game) drawBorder() {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
	for x := 0; x < width; x++ {
		g.draw(Point{X: x, Y: 0}, '═', style)
		g.draw(Point{X: x, Y: height - 1}, '═', style)
	}
	for y := 0; y < height; y++ {
		g.draw(Point{X: 0, Y: y}, '║', style)
		g.draw(Point{X: width - 1, Y: y}, '║', style)
	}
	g.draw(Point{X: 0, Y: 0}, '╔', style)
	g.draw(Point{X: width-1, Y: 0}, '╗', style)
	g.draw(Point{X: 0, Y: height-1}, '╚', style)
	g.draw(Point{X: width-1, Y: height-1}, '╝', style)
}

func (g *Game) drawSnake() {
	headStyle := tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack)
	bodyStyle := tcell.StyleDefault.Foreground(tcell.ColorLimeGreen).Background(tcell.ColorBlack)

	for i, p := range g.Snake.Body {
		if i == 0 {
			g.draw(p, '●', headStyle) // Head
		} else {
			g.draw(p, 'o', bodyStyle) // Body
		}
	}
}

func (g *Game) drawFood() {
	style := tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorBlack)
	g.draw(g.Food, 'F', style)
}

func (g *Game) drawScore() {
	scoreStr := fmt.Sprintf("Score: %d", g.Score)
	style := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorBlack)
	for i, r := range scoreStr {
		g.draw(Point{X: width/2 - len(scoreStr)/2 + i, Y: 0}, r, style)
	}
}

func (g *Game) clearScreen() {
	g.Screen.Clear()
}

func (g *Game) handleInput() {
	for {
		ev := g.Screen.PollEvent()
		if ev == nil {
			continue
		}
		g.InputChan <- ev
	}
}

func (g *Game) processInput() {
	select {
	case ev := <-g.InputChan:
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC, tcell.KeyRune:
				if ev.Rune() == 'q' {
					g.GameOver = true
				}
			case tcell.KeyUp:
				if g.Snake.Direction.Y == 0 { // Cannot reverse direction
					g.Snake.Direction = Point{X: 0, Y: -1}
				}
			case tcell.KeyDown:
				if g.Snake.Direction.Y == 0 {
					g.Snake.Direction = Point{X: 0, Y: 1}
				}
			case tcell.KeyLeft:
				if g.Snake.Direction.X == 0 {
					g.Snake.Direction = Point{X: -1, Y: 0}
				}
			case tcell.KeyRight:
				if g.Snake.Direction.X == 0 {
					g.Snake.Direction = Point{X: 1, Y: 0}
				}
			}
		}
	default:
		// No input, continue
	}
}

func (g *Game) update() {
	if g.GameOver {
		return
	}

	head := g.Snake.Body[0]
	newHead := Point{
		X: head.X + g.Snake.Direction.X,
		Y: head.Y + g.Snake.Direction.Y,
	}

	if g.checkCollision(newHead) {
		g.GameOver = true
		return
	}

	g.Snake.Body = append([]Point{newHead}, g.Snake.Body...)

	if newHead == g.Food {
		g.Score++
		if g.Delay > minDelay {
			g.Delay -= delayDecrease
		}
		g.generateFood()
	} else {
		g.Snake.Body = g.Snake.Body[:len(g.Snake.Body)-1] // Remove tail
	}
}

func (g *Game) checkCollision(head Point) bool {
	// Wall collision
	if head.X <= 0 || head.X >= width-1 || head.Y <= 0 || head.Y >= height-1 {
		return true
	}

	// Self-collision
	for i, p := range g.Snake.Body {
		if i == 0 { // Skip head itself
			continue
		}
		if head == p {
			return true
		}
	}
	return false
}

func (g *Game) generateFood() {
	for {
		foodX := rand.Intn(width-2) + 1  // +1 and -1 for border
		foodY := rand.Intn(height-2) + 1 // +1 and -1 for border
		newFood := Point{X: foodX, Y: foodY}

		isValid := true
		for _, p := range g.Snake.Body {
			if newFood == p {
				isValid = false
				break
			}
		}
		if isValid {
			g.Food = newFood
			return
		}
	}
}

func (g *Game) gameOverScreen() {
	g.clearScreen()
	style := tcell.StyleDefault.Foreground(tcell.ColorRed).Background(tcell.ColorBlack)
	msg1 := "GAME OVER!"
	msg2 := fmt.Sprintf("Final Score: %d", g.Score)
	msg3 := "Press 'q' to quit."

	g.Screen.SetContent(width/2-len(msg1)/2, height/2-2, ' ', nil, style)
	for i, r := range msg1 {
		g.Screen.SetContent(width/2-len(msg1)/2+i, height/2-2, r, nil, style)
	}
	for i, r := range msg2 {
		g.Screen.SetContent(width/2-len(msg2)/2+i, height/2, r, nil, style)
	}
	for i, r := range msg3 {
		g.Screen.SetContent(width/2-len(msg3)/2+i, height/2+2, r, nil, style)
	}
	g.Screen.Show()

	for {
		ev := g.Screen.PollEvent()
		if ev == nil {
			continue
		}
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyRune && ev.Rune() == 'q' {
				return
			}
		}
	}
}

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.HideCursor()

	game := NewGame(s)
	go game.handleInput() // Start input goroutine

	for !game.GameOver {
		game.clearScreen()
		game.drawBorder()
		game.drawSnake()
		game.drawFood()
		game.drawScore()
		s.Show()

		game.processInput() // Process any pending input
		game.update()       // Update game state

		time.Sleep(game.Delay)
	}

	game.gameOverScreen()
	s.Fini()
}

// Additional implementation at 2025-06-20 00:35:04
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Coord struct {
	X, Y int
}

type GameState struct {
	screen    tcell.Screen
	snake     []Coord
	food      Coord
	direction Coord
	score     int
	gameOver  bool
	paused    bool
	width     int
	height    int
	gameSpeed time.Duration
}

const (
	initialSnakeLength = 3
	initialGameSpeed   = 200 * time.Millisecond
	minGameSpeed       = 50 * time.Millisecond
	speedIncreaseScore = 5
	speedDecreaseAmount = 10 * time.Millisecond
)

var (
	styleDefault = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	styleSnake   = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	styleFood    = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
	styleBorder  = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlue)
	styleMessage = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorYellow)
)

func main() {
	rand.Seed(time.Now().UnixNano())

	s, err := tcell.New

// Additional implementation at 2025-06-20 00:36:36
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	BoardWidth     = 40
	BoardHeight    = 20
	InitialSpeed   = 200 * time.Millisecond
	MinSpeed       = 50 * time.Millisecond
	SpeedIncrement = 10 * time.Millisecond
)

type Point struct {
	X, Y int
}

type Snake struct {
	Body      []Point
	Direction Point
	GrowCount int
}

type Food struct {
	Position Point
	Type     int
}

type Game struct {
	Screen    tcell.Screen
	Snake     *Snake
	Food      *Food
	Score     int
	HighScore int
	GameOver  bool
	Paused    bool
	Speed     time.Duration
	Quit      bool
	Level     int
}

var (
	styleDefault    = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	styleSnakeHead  = tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorBlack)
	styleSnakeBody  = tcell.StyleDefault.Background(tcell.ColorDarkGreen).Foreground(tcell.ColorBlack)
	styleFoodNormal = tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite)
	styleFoodBonus  = tcell.StyleDefault.Background(tcell.ColorYellow).Foreground(tcell.ColorBlack)
	styleGameOver   = tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite).Bold(true)
	styleInfo       = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
)

func NewGame() *Game {
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	s.SetStyle(styleDefault)
	s.Clear()

	game := &Game{
		Screen:    s,
		HighScore: loadHighScore(),
		Speed:     InitialSpeed,
		Level:     1,
	}
	game.Reset()
	return game
}

func (g *Game) Reset() {
	g.Snake = &Snake{
		Body:      []Point{{BoardWidth / 2, BoardHeight / 2}},
		Direction: Point{1, 0},
		GrowCount: 0,
	}
	g.Score = 0
	g.GameOver = false
	g.Paused = false
	g.Speed = InitialSpeed
	g.Level = 1
	g.generateFood()
}

func (g *Game) generateFood() {
	for {
		x := rand.Intn(BoardWidth)
		y := rand.Intn(BoardHeight)
		newFood := Food{Position: Point{x, y}}

		onSnake := false
		for _, p := range g.Snake.Body {
			if p.X == x && p.Y == y {
				onSnake = true
				break
			}
		}
		if !onSnake {
			if rand.Intn(10) < 2 {
				newFood.Type = 1
			} else {
				newFood.Type = 0
			}
			g.Food = &newFood
			return
		}
	}
}

func (g *Game) Update() {
	if g.GameOver || g.Paused {
		return
	}

	head := g.Snake.Body[0]
	newHead := Point{head.X + g.Snake.Direction.X, head.Y + g.Snake.Direction.Y}

	if newHead.X < 0 || newHead.X >= BoardWidth || newHead.Y < 0 || newHead.Y >= BoardHeight {
		g.GameOver = true
		return
	}

	for i, p := range g.Snake.Body {
		if i == 0 {
			continue
		}
		if newHead.X == p.X && newHead.Y == p.Y {
			g.GameOver = true
			return
		}
	}

	if newHead.X == g.Food.Position.X && newHead.Y == g.Food.Position.Y {
		g.Snake.GrowCount += 3
		if g.Food.Type == 1 {
			g.Score += 20
			g.Speed = time.Duration(float64(g.Speed) * 0.9)
		} else {
			g.Score += 10
			g.Speed -= SpeedIncrement
		}
		if g.Speed < MinSpeed {
			g.Speed = MinSpeed
		}
		g.Level = int((InitialSpeed - g.Speed) / SpeedIncrement) + 1
		g.generateFood()
	}

	g.Snake.Body = append([]Point{newHead}, g.Snake.Body...)
	if g.Snake.GrowCount > 0 {
		g.Snake.GrowCount--
	} else {
		g.Snake.Body = g.Snake.Body[:len(g.Snake.Body)-1]
	}

	if g.Score > g.HighScore {
		g.HighScore = g.Score
		saveHighScore(g.HighScore)
	}
}

func (g *Game) Draw() {
	g.Screen.Clear()

	for i, p := range g.Snake.Body {
		if i == 0 {
			g.Screen.SetContent(p.X, p.Y, '█', nil, styleSnakeHead)
		} else {
			g.Screen.SetContent(p.X, p.Y, '▓', nil, styleSnakeBody)
		}
	}

	foodRune := '●'
	foodStyle := styleFoodNormal
	if g.Food.Type == 1 {
		foodRune = '★'
		foodStyle = styleFoodBonus
	}
	g.Screen.SetContent(g.Food.Position.X, g.Food.Position.Y, foodRune, nil, foodStyle)

	infoLine := fmt.Sprintf("Score: %d | High Score: %d | Level: %d | Speed: %dms | (P)ause (R)eset (Q)uit", g.Score, g.HighScore, g.Level, g.Speed.Milliseconds())
	g.drawText(0, BoardHeight+1, infoLine, styleInfo)

	if g.GameOver {
		g.drawCenteredText(BoardWidth/2, BoardHeight/2-1, "GAME OVER!", styleGameOver)
		g.drawCenteredText(BoardWidth/2, BoardHeight/2, fmt.Sprintf("Final Score: %d", g.Score), styleGameOver)
		g.drawCenteredText(BoardWidth/2, BoardHeight/2+1, "Press 'R' to Restart or 'Q' to Quit", styleGameOver)
	} else if g.Paused {
		g.drawCenteredText(BoardWidth/2, BoardHeight/2, "PAUSED", styleInfo)
		g.drawCenteredText(BoardWidth/2, BoardHeight/2+1, "Press 'P' to Resume", styleInfo)
	}

	g.Screen.Show()
}

func (g *Game) drawText(x, y int, text string, style tcell.Style) {
	for i, r := range text {
		g.Screen.SetContent(x+i, y, r, nil, style)
	}
}

func (g *Game) drawCenteredText(centerX, y int, text string, style tcell.Style) {
	startX := centerX - len(text)/2
	g.drawText(startX, y, text, style)
}

func (g *Game) HandleInput() {
	for {
		ev := g.Screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if g.GameOver {
				if ev.Key() == tcell.KeyRune && ev.Rune() == 'r' {
					g.Reset()
				} else if ev.Key() == tcell.KeyRune && ev.Rune() == 'q' {
					g.Quit = true
					return
				}
				continue
			}

			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				g.Quit = true
				return
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'p', 'P':
					g.Paused = !g.Paused
				case 'r', 'R':
					g.Reset()
				case 'q', 'Q':
					g.Quit = true
					return
				}
			case tcell.KeyUp:
				if g.Snake.Direction.Y == 0 {
					g.Snake.Direction = Point{0, -1}
				}
			case tcell.KeyDown:
				if g.Snake.Direction.Y == 0 {
					g.Snake.Direction = Point{0, 1}
				}
			case tcell.KeyLeft:
				if g.Snake.Direction.X == 0 {
					g.Snake.Direction = Point{-1, 0}
				}
			case tcell.KeyRight:
				if g.Snake.Direction.X == 0 {
					g.Snake.Direction = Point{1, 0}
				}
			}
		case *tcell.EventResize:
			g.Screen.Sync()
		}
	}
}

const highscoreFile = "snake_highscore.txt"

func loadHighScore() int {
	data, err := os.ReadFile(highscoreFile)
	if err != nil {
		return 0
	}
	var score int
	_, err = fmt.Sscanf(string(data), "%d", &score)
	if err != nil {
		return 0
	}
	return score
}

func saveHighScore(score int) {
	err := os.WriteFile(highscoreFile, []byte(fmt.Sprintf("%d", score)), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving high score: %v\n", err)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	game := NewGame()
	defer game.Screen.Fini()

	go game.HandleInput()

	for !game.Quit {
		game.Update()
		game.Draw()

		if game.GameOver || game.Paused {
			time.Sleep(100 * time.Millisecond)
		} else {
			time.Sleep(game.Speed)
		}
	}
}