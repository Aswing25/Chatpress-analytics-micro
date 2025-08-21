package services

import (
	"chatpress-analytics/db"
	"chatpress-analytics/models"
)

type AnalyticsService struct {
	repo *db.AnalyticsRepository
}

func NewAnalyticsService(repo *db.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) GetAPIUsageCost(userID int) (*models.APIUsageCostResponse, error) {
	return s.repo.GetAPIUsageCost(userID)
}

func (s *AnalyticsService) GetMonthlyUsage(userID int) (*models.MonthlyUsageResponse, error) {
	return s.repo.GetMonthlyUsage(userID)
}

func (s *AnalyticsService) GetOverallStatus(userID int) (*models.OverallStatusResponse, error) {
	return s.repo.GetOverallStatus(userID)
}
