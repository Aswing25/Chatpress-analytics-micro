package db

import (
	"chatpress-analytics/models"
	"database/sql"
	"time"
)

type AnalyticsRepository struct {
	db *DB
}

func NewAnalyticsRepository(db *DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) GetAPIUsageCost(userID int) (*models.APIUsageCostResponse, error) {
	now := time.Now()
	firstDayCurrent := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	firstDayLast := firstDayCurrent.AddDate(0, -1, 0)
	lastDayLast := firstDayCurrent.AddDate(0, 0, -1)

	// Total cost this month
	var thisMonthCost sql.NullFloat64
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(cost_per_request), 0) 
		FROM api_requests 
		WHERE user_id = $1 AND timestamp >= $2`,
		userID, firstDayCurrent).Scan(&thisMonthCost)
	if err != nil {
		return nil, err
	}

	// Last month cost
	var lastMonthCost sql.NullFloat64
	err = r.db.QueryRow(`
		SELECT COALESCE(SUM(cost_per_request), 0) 
		FROM api_requests 
		WHERE user_id = $1 AND timestamp >= $2 AND timestamp <= $3`,
		userID, firstDayLast, lastDayLast).Scan(&lastMonthCost)
	if err != nil {
		return nil, err
	}

	// Total requests
	var totalRequests int
	err = r.db.QueryRow(`
		SELECT COUNT(*) 
		FROM api_requests 
		WHERE user_id = $1 AND timestamp >= $2`,
		userID, firstDayCurrent).Scan(&totalRequests)
	if err != nil {
		return nil, err
	}

	// Daily costs for peak calculation
	rows, err := r.db.Query(`
		SELECT DATE(timestamp) as day, SUM(cost_per_request) as daily_cost
		FROM api_requests 
		WHERE user_id = $1 AND timestamp >= $2
		GROUP BY DATE(timestamp)
		ORDER BY daily_cost DESC
		LIMIT 1`,
		userID, firstDayCurrent)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peakDay string
	var peakCost float64
	if rows.Next() {
		var day time.Time
		err = rows.Scan(&day, &peakCost)
		if err != nil {
			return nil, err
		}
		peakDay = day.Format("2")
	}

	// Calculate averages and changes
	daysCount := int(now.Sub(firstDayCurrent).Hours()/24) + 1
	avgCost := 0.0
	if daysCount > 0 {
		avgCost = thisMonthCost.Float64 / float64(daysCount)
	}

	change := 0.0
	if lastMonthCost.Float64 > 0 {
		change = ((thisMonthCost.Float64 - lastMonthCost.Float64) / lastMonthCost.Float64) * 100
	} else if thisMonthCost.Float64 > 0 {
		change = 100.0
	}

	return &models.APIUsageCostResponse{
		TotalCost:     thisMonthCost.Float64,
		LastMonthCost: lastMonthCost.Float64,
		TotalRequest:  formatCount(totalRequests),
		AvgCost:       avgCost,
		PeakDay:       peakDay,
		Change:        change,
		ChangeUpDown:  change >= 0,
	}, nil
}

func (r *AnalyticsRepository) GetMonthlyUsage(userID int) (*models.MonthlyUsageResponse, error) {
	now := time.Now()
	var results []models.MonthlyUsageItem
	var totalRequests int

	// Get last 6 months data
	for i := 5; i >= 0; i-- {
		monthDate := now.AddDate(0, -i, 0)
		monthStart := time.Date(monthDate.Year(), monthDate.Month(), 1, 0, 0, 0, 0, monthDate.Location())
		monthEnd := monthStart.AddDate(0, 1, 0)

		var count int
		err := r.db.QueryRow(`
			SELECT COUNT(*) 
			FROM api_requests 
			WHERE user_id = $1 AND timestamp >= $2 AND timestamp < $3`,
			userID, monthStart, monthEnd).Scan(&count)
		if err != nil {
			return nil, err
		}

		results = append(results, models.MonthlyUsageItem{
			Month: monthDate.Format("Jan"),
			Value: formatCount(count),
			Raw:   count,
		})
		totalRequests += count
	}

	// Calculate percentages
	if totalRequests == 0 {
		totalRequests = 1 // avoid division by zero
	}
	for i := range results {
		results[i].Percentage = float64(results[i].Raw) / float64(totalRequests) * 100
	}

	// Calculate change
	change := 0.0
	changeUpDown := true
	if len(results) >= 2 {
		prev := results[len(results)-2].Raw
		curr := results[len(results)-1].Raw
		if prev > 0 {
			change = float64(curr-prev) / float64(prev) * 100
		} else if curr > 0 {
			change = 100.0
		}
		changeUpDown = curr >= prev
	}

	return &models.MonthlyUsageResponse{
		FormattedChart: results,
		Change:         change,
		ChangeUpDown:   changeUpDown,
	}, nil
}

func (r *AnalyticsRepository) GetOverallStatus(userID int) (*models.OverallStatusResponse, error) {
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Total chats this month
	var totalChats int
	err := r.db.QueryRow(`
		SELECT COUNT(DISTINCT cm.session_id)
		FROM chat_messages cm
		JOIN chat_sessions cs ON cm.session_id = cs.id
		WHERE cs.user_id = $1 AND cm.created_at >= $2 AND cm.sender = 'user'`,
		userID, firstDay).Scan(&totalChats)
	if err != nil {
		return nil, err
	}

	// Unique users (sessions)
	var uniqueUsers int
	err = r.db.QueryRow(`
		SELECT COUNT(DISTINCT session_id)
		FROM api_requests
		WHERE user_id = $1`,
		userID).Scan(&uniqueUsers)
	if err != nil {
		return nil, err
	}

	// Average response time
	avgResponseTime := 0.0
	rows, err := r.db.Query(`
		SELECT cm.created_at, cm.sender
		FROM chat_messages cm
		JOIN chat_sessions cs ON cm.session_id = cs.id
		WHERE cs.user_id = $1 AND cm.sender IN ('user', 'bot')
		ORDER BY cm.created_at DESC
		LIMIT 50`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []struct {
		CreatedAt time.Time
		Sender    string
	}
	for rows.Next() {
		var msg struct {
			CreatedAt time.Time
			Sender    string
		}
		err = rows.Scan(&msg.CreatedAt, &msg.Sender)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	// Calculate response times
	var responseTimes []float64
	var lastUserTime *time.Time
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg.Sender == "user" {
			lastUserTime = &msg.CreatedAt
		} else if msg.Sender == "bot" && lastUserTime != nil {
			delta := msg.CreatedAt.Sub(*lastUserTime).Seconds()
			if delta > 0 && delta < 60 {
				responseTimes = append(responseTimes, delta)
			}
			lastUserTime = nil
		}
	}

	if len(responseTimes) > 0 {
		var sum float64
		for _, rt := range responseTimes {
			sum += rt
		}
		avgResponseTime = sum / float64(len(responseTimes))
	}

	// API usage this month
	var apiUsage int
	err = r.db.QueryRow(`
		SELECT COUNT(*)
		FROM api_requests
		WHERE user_id = $1 AND timestamp >= $2`,
		userID, firstDay).Scan(&apiUsage)
	if err != nil {
		return nil, err
	}

	return &models.OverallStatusResponse{
		TotalChats:      formatCount(totalChats),
		UniqueUsers:     formatCount(uniqueUsers),
		AvgResponseTime: avgResponseTime,
		APIUsage:        formatCount(apiUsage),
	}, nil
}

func formatCount(count int) string {
	if count < 1000 {
		return string(rune(count))
	} else if count < 1000000 {
		return string(rune(count/1000)) + "K"
	} else {
		return string(rune(count/1000000)) + "M"
	}
}
