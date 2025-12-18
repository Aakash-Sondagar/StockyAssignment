package main

import (
	"database/sql"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var db *sql.DB
var log = logrus.New()

type RewardRequest struct {
	UserID      int     `json:"user_id" binding:"required"`
	StockSymbol string  `json:"stock_symbol" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required"`
}

type DailyStockResponse struct {
	StockSymbol string  `json:"stock_symbol"`
	Quantity    float64 `json:"total_quantity"`
}

type StatsResponse struct {
	TotalSharesByStock map[string]float64 `json:"total_shares_by_stock"`
	CurrentPortfolioValueINR float64      `json:"current_portfolio_value_inr"`
}

func GetCurrentStockPrice(symbol string) float64 {
	return 500.0 + rand.Float64()*(3000.0-500.0)
}

func createReward(c *gin.Context) {
	var req RewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentPrice := GetCurrentStockPrice(req.StockSymbol)
	stockCost := currentPrice * req.Quantity
	brokerageFee := stockCost * 0.01
	totalCompanyCost := stockCost + brokerageFee

	tx, err := db.Begin()
	if err != nil {
		log.Error("Failed to start transaction")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var rewardID int
	err = tx.QueryRow(`
		INSERT INTO rewards (user_id, stock_symbol, quantity) 
		VALUES ($1, $2, $3) RETURNING id`, 
		req.UserID, req.StockSymbol, req.Quantity).Scan(&rewardID)

	if err != nil {
		tx.Rollback()
		log.Error("Failed to insert reward: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record reward"})
		return
	}

	_, err = tx.Exec(`INSERT INTO ledger (reward_id, account_type, description, amount, currency) VALUES ($1, 'USER_ASSET', 'Stock Reward', $2, $3)`,
		rewardID, req.Quantity, req.StockSymbol)

	_, err = tx.Exec(`INSERT INTO ledger (reward_id, account_type, description, amount, currency) VALUES ($1, 'COMPANY_CASH', 'Stock Purchase', $2, 'INR')`,
		rewardID, -stockCost)

	_, err = tx.Exec(`INSERT INTO ledger (reward_id, account_type, description, amount, currency) VALUES ($1, 'FEE_EXPENSE', 'Brokerage/Taxes', $2, 'INR')`,
		rewardID, -brokerageFee)

	if err != nil {
		tx.Rollback()
		log.Error("Ledger entry failed: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ledger error"})
		return
	}

	tx.Commit()
	log.WithFields(logrus.Fields{
		"user_id": req.UserID,
		"stock":   req.StockSymbol,
		"cost":    totalCompanyCost,
	}).Info("Reward processed successfully")

	c.JSON(http.StatusCreated, gin.H{"message": "Reward created", "reward_id": rewardID})
}

func getTodayStocks(c *gin.Context) {
	userID := c.Param("userId")
	
	rows, err := db.Query(`
		SELECT stock_symbol, SUM(quantity) 
		FROM rewards 
		WHERE user_id = $1 AND rewarded_at >= CURRENT_DATE 
		GROUP BY stock_symbol`, userID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
		return
	}
	defer rows.Close()

	var results []DailyStockResponse
	for rows.Next() {
		var r DailyStockResponse
		rows.Scan(&r.StockSymbol, &r.Quantity)
		results = append(results, r)
	}

	c.JSON(http.StatusOK, results)
}

func getUserStats(c *gin.Context) {
	userID := c.Param("userId")

	rows, err := db.Query(`SELECT stock_symbol, SUM(quantity) FROM rewards WHERE user_id = $1 GROUP BY stock_symbol`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Error"})
		return
	}
	defer rows.Close()

	holdings := make(map[string]float64)
	var totalValueINR float64

	for rows.Next() {
		var sym string
		var qty float64
		rows.Scan(&sym, &qty)
		
		holdings[sym] = qty
		
		currentPrice := GetCurrentStockPrice(sym) 
		totalValueINR += (qty * currentPrice)
	}

	resp := StatsResponse{
		TotalSharesByStock: holdings,
		CurrentPortfolioValueINR: totalValueINR,
	}

	c.JSON(http.StatusOK, resp)
}

func getHistoricalINR(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Returns historical valuation based on daily snapshots"})
}

func main() {
	log.SetFormatter(&logrus.JSONFormatter{})

	connStr := "user=postgres dbname=assignment sslmode=disable password=postgres"
	var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }

    r := gin.Default()

    r.POST("/reward", createReward)
    r.GET("/today-stocks/:userId", getTodayStocks)
    r.GET("/stats/:userId", getUserStats)
    r.GET("/historical-inr/:userId", getHistoricalINR)

    log.Info("Server starting on port 8080")
    r.Run(":8080")
}