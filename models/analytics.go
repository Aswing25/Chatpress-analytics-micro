package models

import (
	"time"
)

type User struct {
	ID      int       `json:"id"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	Created time.Time `json:"created_at"`
}

type ApiRequest struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	SessionID      string    `json:"session_id"`
	CostPerRequest float64   `json:"cost_per_request"`
	Timestamp      time.Time `json:"timestamp"`
}

type ChatSession struct {
	UserID int `json:"user_id"`
}

type ChatMessage struct {
	ID        int       `json:"id"`
	SessionID int       `json:"session_id"`
	Sender    string    `json:"sender"`
	CreatedAt time.Time `json:"created_at"`
}

// Response models
type APIUsageCostResponse struct {
	TotalCost     float64 `json:"total_cost"`
	LastMonthCost float64 `json:"last_month_cost"`
	TotalRequest  string  `json:"total_request"`
	AvgCost       float64 `json:"avg_cost"`
	PeakDay       string  `json:"peak_day"`
	Change        float64 `json:"change"`
	ChangeUpDown  bool    `json:"change_up_down"`
}

type MonthlyUsageItem struct {
	Month      string  `json:"month"`
	Value      string  `json:"value"`
	Raw        int     `json:"raw"`
	Percentage float64 `json:"percentage"`
}

type MonthlyUsageResponse struct {
	FormattedChart []MonthlyUsageItem `json:"formatted_chart"`
	Change         float64            `json:"change"`
	ChangeUpDown   bool               `json:"change_up_down"`
}

type OverallStatusResponse struct {
	TotalChats      string  `json:"total_chats"`
	UniqueUsers     string  `json:"unique_users"`
	AvgResponseTime float64 `json:"avg_response_time"`
	APIUsage        string  `json:"api_usage"`
}
