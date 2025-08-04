package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Stock struct {
	Symbol string
	Price  float64
	Change float64
}

func (s *Stock) updatePrice() {
	oldPrice := s.Price
	changePercent := (rand.Float64()*2 - 1) * 0.02
	s.Price *= (1 + changePercent)
	s.Change = s.Price - oldPrice
	if s.Price < 0.01 {
		s.Price = 0.01
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func displayTicker(stocks []*Stock) {
	clearScreen()
	fmt.Println("--- Go Stock Ticker ---")
	fmt.Println("Symbol\tPrice\tChange")
	fmt.Println("-----------------------")
	for _, s := range stocks {
		changeStr := fmt.Sprintf("%.2f", s.Change)
		if s.Change > 0 {
			changeStr = "+" + changeStr
		}
		fmt.Printf("%s\t%.2f\t%s\n", s.Symbol, s.Price, changeStr)
	}
	fmt.Println("-----------------------")
	fmt.Println("Updating every 1 second...")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	stocks := []*Stock{
		{Symbol: "GOOG", Price: 1500.00},
		{Symbol: "AAPL", Price: 170.00},
		{Symbol: "MSFT", Price: 350.00},
		{Symbol: "AMZN", Price: 140.00},
		{Symbol: "TSLA", Price: 250.00},
	}

	for {
		for _, s := range stocks {
			s.updatePrice()
		}
		displayTicker(stocks)
		time.Sleep(1 * time.Second)
	}
}

// Additional implementation at 2025-08-04 07:59:41
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

type Stock struct {
	Symbol     string
	Name       string
	Price      float64
	LastChange float64
}

type Market struct {
	Stocks map[string]*Stock
}

func NewMarket() *Market {
	return &Market{
		Stocks: make(map[string]*Stock),
	}
}

func (m *Market) AddStock(symbol, name string, initialPrice float64) {
	m.Stocks[symbol] = &Stock{
		Symbol:     symbol,
		Name:       name,
		Price:      initialPrice,
		LastChange: 0.0,
	}
}

func (s *Stock) updateStockPrice() {
	// Simulate price change: -2.5% to +2.5%
	changePercent := (rand.Float64() - 0.5) * 0.05
	changeAmount := s.Price * changePercent
	s.Price += changeAmount
	s.LastChange = changeAmount

	// Ensure price doesn't go below a trivial amount
	if s.Price < 0.01 {
		s.Price = 0.01
	}
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (m *Market) displayMarket() {
	clearScreen()
	fmt.Println("--------------------------------------------------")
	fmt.Println("                GO STOCK TICKER                   ")
	fmt.Println("--------------------------------------------------")
	fmt.Printf("%-8s %-20s %-10s %-10s\n", "SYMBOL", "NAME", "PRICE", "CHANGE")
	fmt.Println("--------------------------------------------------")

	for _, stock := range m.Stocks {
		changeStr := fmt.Sprintf("%.2f", stock.LastChange)
		if stock.LastChange > 0 {
			changeStr = "+" + changeStr
		}
		fmt.Printf("%-8s %-20s %-10.2f %-10s\n", stock.Symbol, stock.Name, stock.Price, changeStr)
	}
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Last Updated: %s\n", time.Now().Format("15:04:05"))
}

func (m *Market) RunSimulation(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Initial display
	for _, stock := range m.Stocks {
		stock.updateStockPrice() // Initial price update
	}
	m.displayMarket()

	for range ticker.C {
		for _, stock := range m.Stocks {
			stock.updateStockPrice()
		}
		m.displayMarket()
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	market := NewMarket()
	market.AddStock("GOOG", "Alphabet Inc.", 1500.00)
	market.AddStock("AAPL", "Apple Inc.", 170.00)
	market.AddStock("MSFT", "Microsoft Corp.", 300.00)
	market.AddStock("AMZN", "Amazon.com Inc.", 100.00)
	market.AddStock("TSLA", "Tesla Inc.", 250.00)
	market.AddStock("NVDA", "NVIDIA Corp.", 450.00)
	market.AddStock("NFLX", "Netflix Inc.", 400.00)
	market.AddStock("DIS", "Walt Disney Co.", 90.00)

	fmt.Println("Starting Go Stock Ticker Simulation. Press Ctrl+C to exit.")
	time.Sleep(2 * time.Second) // Give user time to read initial message

	market.RunSimulation(time.Second) // Update every 1 second
}