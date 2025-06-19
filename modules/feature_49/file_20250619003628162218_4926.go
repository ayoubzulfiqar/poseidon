package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Stock struct {
	Symbol string
	Price  float64
	Change float64
}

func (s *Stock) updatePrice() {
	oldPrice := s.Price
	changePercent := (rand.Float64()*2 - 1) * 0.01
	s.Price += s.Price * changePercent
	if s.Price < 0.01 {
		s.Price = 0.01
	}
	s.Change = s.Price - oldPrice
}

func clearScreen() {
	os.Stdout.WriteString("\033[H\033[2J")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	stocks := []*Stock{
		{Symbol: "AAPL", Price: 150.25},
		{Symbol: "GOOG", Price: 2500.10},
		{Symbol: "MSFT", Price: 300.50},
		{Symbol: "AMZN", Price: 3200.75},
		{Symbol: "TSLA", Price: 850.00},
		{Symbol: "NFLX", Price: 400.00},
		{Symbol: "NVDA", Price: 600.00},
		{Symbol: "FB", Price: 350.00},
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	fmt.Println("Starting stock ticker simulation... Press Ctrl+C to exit.")

	for range ticker.C {
		clearScreen()

		var tickerOutput strings.Builder
		for _, stock := range stocks {
			stock.updatePrice()

			priceStr := strconv.FormatFloat(stock.Price, 'f', 2, 64)
			changeValStr := strconv.FormatFloat(stock.Change, 'f', 2, 64)

			var changeIndicator string
			if stock.Change > 0 {
				changeIndicator = fmt.Sprintf("▲ %s", changeValStr)
			} else if stock.Change < 0 {
				changeIndicator = fmt.Sprintf("▼ %s", changeValStr[1:])
			} else {
				changeIndicator = "— 0.00"
			}

			tickerOutput.WriteString(fmt.Sprintf("%s: $%s (%s) | ", stock.Symbol, priceStr, changeIndicator))
		}
		finalOutput := strings.TrimSuffix(tickerOutput.String(), " | ")
		fmt.Println(finalOutput)
	}
}

// Additional implementation at 2025-06-19 00:37:09
package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Stock represents a single stock in the market
type Stock struct {
	Symbol          string
	Name            string
	Price           float64
	Volatility      float64 // Factor for price change (e.g., 0.05 for 5% max change)
	mu              sync.RWMutex // Mutex for concurrent access to Price
	LastPriceChange float64 // To indicate if price went up or down
}

// Portfolio represents a user's holdings and cash
type Portfolio struct {
	Cash     float64
	Holdings map[string]int // Stock Symbol -> Quantity
	mu       sync.RWMutex // Mutex for concurrent access to Cash and Holdings
}

// Market holds all available stocks
type Market struct {
	Stocks map[string]*Stock
}

// NewStock creates a new Stock instance
func NewStock(symbol, name string, initialPrice, volatility float64) *Stock {
	return &Stock{
		Symbol:     symbol,
		Name:       name,
		Price:      initialPrice,
		Volatility: volatility,
	}
}

// NewPortfolio creates a new Portfolio instance
func NewPortfolio(initialCash float64) *Portfolio {
	return &Portfolio{
		Cash:     initialCash,
		Holdings: make(map[string]int),
	}
}

// NewMarket creates a new Market instance with predefined stocks
func NewMarket() *Market {
	stocks := make(map[string]*Stock)
	stocks["GOOG"] = NewStock("GOOG", "Google Inc.", 1500.00, 0.02)
	stocks["AAPL"] = NewStock("AAPL", "Apple Inc.", 170.00, 0.03)
	stocks["MSFT"] = NewStock("MSFT", "Microsoft Corp.", 300.00, 0.025)
	stocks["AMZN"] = NewStock("AMZN", "Amazon.com Inc.", 100.00, 0.04)
	stocks["TSLA"] = NewStock("TSLA", "Tesla Inc.", 250.00, 0.05)
	return &Market{Stocks: stocks}
}

// updatePrice simulates a price change for a stock
func (s *Stock) updatePrice() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate a random change based on volatility
	changePercent := (rand.Float64()*2 - 1) * s.Volatility // -volatility to +volatility
	changeAmount := s.Price * changePercent
	newPrice := s.Price + changeAmount

	// Ensure price doesn't go below a minimum
	if newPrice < 0.01 {
		newPrice = 0.01
	}

	s.LastPriceChange = newPrice - s.Price
	s.Price = newPrice
}

// runMarketUpdates continuously updates stock prices and simulates news events
func runMarketUpdates(market *Market, updateInterval time.Duration) {
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for range ticker.C {
		// Update all stock prices
		for _, stock := range market.Stocks {
			stock.updatePrice()
		}

		// Simulate a random news event
		if rand.Float64() < 0.1 { // 10% chance of a news event per update cycle
			symbols := make([]string, 0, len(market.Stocks))
			for s := range market.Stocks {
				symbols = append(symbols, s)
			}
			if len(symbols) > 0 {
				affectedSymbol := symbols[rand.Intn(len(symbols))]
				affectedStock := market.Stocks[affectedSymbol]

				affectedStock.mu.Lock()
				newsEffect := (rand.Float64()*2 - 1) * 0.1 // Up to 10% boost/drop
				affectedStock.Price *= (1 + newsEffect)
				affectedStock.LastPriceChange = affectedStock.Price * newsEffect // Approximate
				if affectedStock.Price < 0.01 {
					affectedStock.Price = 0.01
				}
				affectedStock.mu.Unlock()

				fmt.Printf("\n[NEWS] %s: %s %s! Price changed by %.2f%%\n",
					affectedStock.Symbol,
					affectedStock.Name,
					map[bool]string{true: "surges", false: "plummets"}[newsEffect > 0],
					math.Abs(newsEffect*100),
				)
			}
		}
	}
}

// displayMarket prints the current stock prices
func displayMarket(market *Market) {
	fmt.Println("\n--- Market Prices ---")
	fmt.Printf("%-8s %-20s %-10s %-10s\n", "Symbol", "Name", "Price", "Change")
	for _, stock := range market.Stocks {
		stock.mu.RLock() // Read lock
		changeStr := ""
		if stock.LastPriceChange > 0 {
			changeStr = fmt.Sprintf("▲ %.2f", stock.LastPriceChange)
		} else if stock.LastPriceChange < 0 {
			changeStr = fmt.Sprintf("▼ %.2f", math.Abs(stock.LastPriceChange))
		} else {
			changeStr = "  0.00"
		}
		fmt.Printf("%-8s %-20s $%-9.2f %-10s\n", stock.Symbol, stock.Name, stock.Price, changeStr)
		stock.mu.RUnlock() // Read unlock
	}
	fmt.Println("---------------------")
}

// displayPortfolio prints the user's cash and holdings
func displayPortfolio(portfolio *Portfolio, market *Market) {
	portfolio.mu.RLock() // Read lock for portfolio
	defer portfolio.mu.RUnlock()

	fmt.Println("\n--- Your Portfolio ---")
	fmt.Printf("Cash: $%.2f\n", portfolio.Cash)

	totalHoldingsValue := 0.0
	if len(portfolio.Holdings) > 0 {
		fmt.Printf("%-8s %-10s %-10s %-10s\n", "Symbol", "Shares", "Current Price", "Value")
		for symbol, quantity := range portfolio.Holdings {
			if stock, ok := market.Stocks[symbol]; ok {
				stock.mu.RLock() // Read lock for stock price
				currentValue := stock.Price * float64(quantity)
				totalHoldingsValue += currentValue
				fmt.Printf("%-8s %-10d $%-9.2f $%-9.2f\n", symbol, quantity, stock.Price, currentValue)
				stock.mu.RUnlock() // Read unlock for stock price
			}
		}
	} else {
		fmt.Println("You currently hold no stocks.")
	}
	fmt.Printf("Total Holdings Value: $%.2f\n", totalHoldingsValue)
	fmt.Printf("Total Portfolio Value (Cash + Holdings): $%.2f\n", portfolio.Cash+totalHoldingsValue)
	fmt.Println("----------------------")
}

// buyStock handles the purchase of stocks
func buyStock(portfolio *Portfolio, market *Market, symbol string, quantity int) {
	stock, ok := market.Stocks[strings.ToUpper(symbol)]
	if !ok {
		fmt.Printf("Error: Stock '%s' not found.\n", symbol)
		return
	}
	if quantity <= 0 {
		fmt.Println("Error: Quantity must be positive.")
		return
	}

	stock.mu.RLock() // Read lock for stock price
	cost := stock.Price * float64(quantity)
	stock.mu.RUnlock() // Read unlock

	portfolio.mu.Lock() // Write lock for portfolio
	defer portfolio.mu.Unlock()

	if portfolio.Cash < cost {
		fmt.Printf("Error: Insufficient cash. You have $%.2f, need $%.2f.\n", portfolio.Cash, cost)
		return
	}

	portfolio.Cash -= cost
	portfolio.Holdings[stock.Symbol] += quantity
	fmt.Printf("Successfully bought %d shares of %s for $%.2f. Remaining cash: $%.2f.\n", quantity, stock.Symbol, cost, portfolio.Cash)
}

// sellStock handles the selling of stocks
func sellStock(portfolio *Portfolio, market *Market, symbol string, quantity int) {
	stock, ok := market.Stocks[strings.ToUpper(symbol)]
	if !ok {
		fmt.Printf("Error: Stock '%s' not found.\n", symbol)
		return
	}
	if quantity <= 0 {
		fmt.Println("Error: Quantity must be positive.")
		return
	}

	portfolio.mu.Lock() // Write lock for portfolio
	defer portfolio.mu.Unlock()

	heldQuantity, hasShares := portfolio.Holdings[stock.Symbol]
	if !hasShares || heldQuantity < quantity {
		fmt.Printf("Error: You only have %d shares of %s, cannot sell %d.\n", heldQuantity, stock.Symbol, quantity)
		return
	}

	stock.mu.RLock() // Read lock for stock price
	revenue := stock.Price * float64(quantity)
	stock.mu.RUnlock() // Read unlock

	portfolio.Cash += revenue
	portfolio.Holdings[stock.Symbol] -= quantity
	if portfolio.Holdings[stock.Symbol] == 0 {
		delete(portfolio.Holdings, stock.Symbol)
	}
	fmt.Printf("Successfully sold %d shares of %s for $%.2f. New cash: $%.2f.\n", quantity, stock.Symbol, revenue, portfolio.Cash)
}

// printHelp displays available commands
func printHelp() {
	fmt.Println("\n--- Commands ---")
	fmt.Println("  buy <SYMBOL> <QUANTITY> - Buy shares of a stock (e.g., buy GOOG 10)")
	fmt.Println("  sell <SYMBOL> <QUANTITY> - Sell shares of a stock (e.g., sell AAPL 5)")
	fmt.Println("  view - Display current market prices and your portfolio")
	fmt.Println("  help - Show this help message")
	fmt.Println("  exit - Exit the simulation")
	fmt.Println("----------------")
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	market := NewMarket()
	portfolio := NewPortfolio(10000.00) // Starting cash

	// Start market updates in a goroutine
	go runMarketUpdates(market, 2*time.Second) // Update every 2 seconds

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Welcome to the Go Stock Ticker Simulation!")
	printHelp()

	for {
		// For better display, you might clear the screen here on Unix-like systems:
		// fmt.Print("\033[H\033[2J")
		displayMarket(market)
		displayPortfolio(portfolio, market)

		fmt.Print("\nEnter command (type 'help' for options): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		command := strings.ToLower(parts[0])

		switch command {
		case "buy":
			if len(parts) == 3 {
				symbol := parts[1]
				quantity, err := strconv.Atoi(parts[2])
				if err != nil {
					fmt.Println("Error: Invalid quantity. Please enter a number.")
					continue
				}
				buyStock(portfolio, market, symbol, quantity)
			} else {
				fmt.Println("Usage: buy <SYMBOL> <QUANTITY>")
			}
		case "sell":
			if len(parts) == 3 {
				symbol := parts[1]
				quantity, err := strconv.Atoi(parts[2])
				if err != nil {
					fmt.Println("Error: Invalid quantity. Please enter a number.")
					continue
				}
				sellStock(portfolio, market, symbol, quantity)
			} else {
				fmt.Println("Usage: sell <SYMBOL> <QUANTITY>")
			}
		case "view":
			// displayMarket and displayPortfolio are already called at the top of the loop
			// This command just forces a refresh without waiting for the next market update
			fmt.Println("Refreshing display...")
		case "help":
			printHelp()
		case "exit":
			fmt.Println("Exiting stock ticker simulation. Goodbye!")
			return
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}
		// Small sleep to prevent rapid updates if user enters invalid commands quickly
		time.Sleep(500 * time.Millisecond)
	}
}

// Additional implementation at 2025-06-19 00:38:35
package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Stock struct {
	Symbol        string
	Name          string
	CurrentPrice  float64
	PreviousPrice float64
}

type Portfolio struct {
	Holdings map[string]int
	Cash     float64
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func updateStockPrice(s *Stock) {
	s.PreviousPrice = s.CurrentPrice
	changePercent := (rand.Float64() - 0.5) * 0.05
	changeAmount := s.CurrentPrice * changePercent
	s.CurrentPrice += changeAmount

	if s.CurrentPrice < 0.01 {
		s.CurrentPrice = 0.01
	}
}

func displayTicker(stocks map[string]*Stock) {
	fmt.Println("----------------------------------------------------")
	fmt.Println("               GO STOCK TICKER SIMULATION           ")
	fmt.Println("----------------------------------------------------")
	fmt.Printf("%-8s %-20s %-10s %-8s\n", "SYMBOL", "NAME", "PRICE", "CHANGE")
	fmt.Println("----------------------------------------------------")

	for _, stock := range stocks {
		change := stock.CurrentPrice - stock.PreviousPrice
		changeStr := ""
		if change > 0 {
			changeStr = fmt.Sprintf("▲ %.2f", change)
		} else if change < 0 {
			changeStr = fmt.Sprintf("▼ %.2f", -change)
		} else {
			changeStr = "— 0.00"
		}
		fmt.Printf("%-8s %-20s $%-9.2f %-8s\n", stock.Symbol, stock.Name, stock.CurrentPrice, changeStr)
	}
	fmt.Println("----------------------------------------------------")
}

func displayPortfolio(p *Portfolio, stocks map[string]*Stock) {
	fmt.Println("\n--- YOUR PORTFOLIO ---")
	fmt.Printf("Cash: $%.2f\n", p.Cash)
	fmt.Println("Holdings:")
	if len(p.Holdings) == 0 {
		fmt.Println("  No stocks owned.")
	} else {
		totalPortfolioValue := 0.0
		for symbol, quantity := range p.Holdings {
			if stock, ok := stocks[symbol]; ok {
				value := stock.CurrentPrice * float64(quantity)
				fmt.Printf("  %s (%s): %d shares @ $%.2f/share = $%.2f\n", symbol, stock.Name, quantity, stock.CurrentPrice, value)
				totalPortfolioValue += value
			}
		}
		fmt.Printf("Total Stock Value: $%.2f\n", totalPortfolioValue)
		fmt.Printf("Total Net Worth: $%.2f\n", p.Cash+totalPortfolioValue)
	}
	fmt.Println("----------------------")
}

func handleBuyCommand(input string, p *Portfolio, stocks map[string]*Stock) {
	parts := strings.Fields(input)
	if len(parts) != 3 || strings.ToLower(parts[0]) != "buy" {
		fmt.Println("Invalid buy command. Usage: buy <symbol> <quantity>")
		return
	}

	symbol := strings.ToUpper(parts[1])
	quantity, err := strconv.Atoi(parts[2])
	if err != nil || quantity <= 0 {
		fmt.Println("Invalid quantity. Must be a positive integer.")
		return
	}

	stock, ok := stocks[symbol]
	if !ok {
		fmt.Printf("Stock '%s' not found.\n", symbol)
		return
	}

	cost := stock.CurrentPrice * float64(quantity)
	if p.Cash < cost {
		fmt.Printf("Insufficient cash. You need $%.2f but have $%.2f.\n", cost, p.Cash)
		return
	}

	p.Cash -= cost
	p.Holdings[symbol] += quantity
	fmt.Printf("Successfully bought %d shares of %s for $%.2f.\n", quantity, symbol, cost)
}

func handleSellCommand(input string, p *Portfolio, stocks map[string]*Stock) {
	parts := strings.Fields(input)
	if len(parts) != 3 || strings.ToLower(parts[0]) != "sell" {
		fmt.Println("Invalid sell command. Usage: sell <symbol> <quantity>")
		return
	}

	symbol := strings.ToUpper(parts[1])
	quantity, err := strconv.Atoi(parts[2])
	if err != nil || quantity <= 0 {
		fmt.Println("Invalid quantity. Must be a positive integer.")
		return
	}

	stock, ok := stocks[symbol]
	if !ok {
		fmt.Printf("Stock '%s' not found.\n", symbol)
		return
	}

	currentHoldings, hasStock := p.Holdings[symbol]
	if !hasStock || currentHoldings < quantity {
		fmt.Printf("You only own %d shares of %s. Cannot sell %d.\n", currentHoldings, symbol, quantity)
		return
	}

	revenue := stock.CurrentPrice * float64(quantity)
	p.Cash += revenue
	p.Holdings[symbol] -= quantity
	if p.Holdings[symbol] == 0 {
		delete(p.Holdings, symbol)
	}
	fmt.Printf("Successfully sold %d shares of %s for $%.2f.\n", quantity, symbol, revenue)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	stocks := map[string]*Stock{
		"GOOG": {Symbol: "GOOG", Name: "Alphabet Inc.", CurrentPrice: 1500.00, PreviousPrice: 1500.00},
		"AAPL": {Symbol: "AAPL", Name: "Apple Inc.", CurrentPrice: 170.00, PreviousPrice: 170.00},
		"MSFT": {Symbol: "MSFT", Name: "Microsoft Corp.", CurrentPrice: 350.00, PreviousPrice: 350.00},
		"AMZN": {Symbol: "AMZN", Name: "Amazon.com Inc.", CurrentPrice: 130.00, PreviousPrice: 130.00},
		"TSLA": {Symbol: "TSLA", Name: "Tesla Inc.", CurrentPrice: 250.00, PreviousPrice: 250.00},
		"NVDA": {Symbol: "NVDA", Name: "NVIDIA Corp.", CurrentPrice: 480.00, PreviousPrice: 480.00},
	}

	portfolio := &Portfolio{
		Holdings: make(map[string]int),
		Cash:     10000.00,
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Go Stock Ticker Simulation!")
	fmt.Println("Type 'buy <symbol> <quantity>' to buy stocks (e.g., 'buy AAPL 10').")
	fmt.Println("Type 'sell <symbol> <quantity>' to sell stocks (e.g., 'sell AAPL 5').")
	fmt.Println("Type 'q' or 'quit' to exit.")
	fmt.Println("Press Enter to continue...")
	reader.ReadString('\n')

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			clearScreen()
			for _, stock := range stocks {
				updateStockPrice(stock)
			}
			displayTicker(stocks)
			displayPortfolio(portfolio, stocks)
			fmt.Print("\nEnter command (buy/sell/q): ")
		}
	}()

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		lowerInput := strings.ToLower(input)

		if lowerInput == "q" || lowerInput == "quit" {
			fmt.Println("Exiting simulation. Goodbye!")
			return
		} else if strings.HasPrefix(lowerInput, "buy ") {
			handleBuyCommand(input, portfolio, stocks)
		} else if strings.HasPrefix(lowerInput, "sell ") {
			handleSellCommand(input, portfolio, stocks)
		} else {
			fmt.Println("Unknown command. Type 'buy', 'sell', or 'q'.")
		}
	}
}

// Additional implementation at 2025-06-19 00:40:08
package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

type Stock struct {
	Symbol    string
	Name      string
	Price     float64
	LastChange float64
}

type Market struct {
	Stocks []*Stock
	mu     sync.Mutex
	rng    *rand.Rand
}

func NewMarket() *Market {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	return &Market{
		Stocks: []*Stock{
			{Symbol: "GOOG", Name: "Alphabet Inc.", Price: 1500.00, LastChange: 0.0},
			{Symbol: "AAPL", Name: "Apple Inc.", Price: 170.00, LastChange: 0.0},
			{Symbol: "MSFT", Name: "Microsoft Corp.", Price: 350.00, LastChange: 0.0},
			{Symbol: "AMZN", Name: "Amazon.com Inc.", Price: 140.00, LastChange: 0.0},
			{Symbol: "TSLA", Name: "Tesla Inc.", Price: 250.00, LastChange: 0.0},
			{Symbol: "NVDA", Name: "NVIDIA Corp.", Price: 450.00, LastChange: 0.0},
			{Symbol: "META", Name: "Meta Platforms Inc.", Price: 300.00, LastChange: 0.0},
			{Symbol: "NFLX", Name: "Netflix Inc.", Price: 400.00, LastChange: 0.0},
			{Symbol: "ADBE", Name: "Adobe Inc.", Price: 500.00, LastChange: 0.0},
			{Symbol: "CRM", Name: "Salesforce Inc.", Price: 200.00, LastChange: 0.0},
		},
		rng: r,
	}
}

func (m *Market) updatePrices() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, stock := range m.Stocks {
		changePercent := (m.rng.Float64()*2 - 1) * 0.005
		newPrice := stock.Price * (1 + changePercent)

		if newPrice < 0.01 {
			newPrice = 0.01
		}

		stock.LastChange = newPrice - stock.Price
		stock.Price = newPrice
	}
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func (m *Market) displayTicker() {
	clearScreen()
	fmt.Println("-------------------------------------------------------------------")
	fmt.Println("                     GO STOCK TICKER SIMULATION                    ")
	fmt.Println("-------------------------------------------------------------------")
	fmt.Printf("%-8s %-25s %-10s %-10s\n", "SYMBOL", "NAME", "PRICE", "CHANGE")
	fmt.Println("-------------------------------------------------------------------")

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, stock := range m.Stocks {
		changeStr := fmt.Sprintf("%.2f", stock.LastChange)
		if stock.LastChange > 0 {
			changeStr = "+" + changeStr
		}
		fmt.Printf("%-8s %-25s $%-9.2f %-10s\n", stock.Symbol, stock.Name, stock.Price, changeStr)
	}
	fmt.Println("-------------------------------------------------------------------")
	fmt.Printf("Last Updated: %s\n", time.Now().Format("15:04:05"))
	fmt.Println("Press Ctrl+C to exit.")
}

func (m *Market) runSimulation(updateInterval time.Duration) {
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for range ticker.C {
		m.updatePrices()
		m.displayTicker()
	}
}

func main() {
	market := NewMarket()

	go market.runSimulation(1 * time.Second)

	select {}
}