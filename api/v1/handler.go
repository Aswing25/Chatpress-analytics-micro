package v1

import (
	"chatpress-analytics/db"
	"chatpress-analytics/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	analyticsService *services.AnalyticsService
	jwtSecret        string
}

func NewHandler(repo *db.AnalyticsRepository, jwtSecret string) *Handler {
	return &Handler{
		analyticsService: services.NewAnalyticsService(repo),
		jwtSecret:        jwtSecret,
	}
}

func (h *Handler) getCurrentUser(c *gin.Context) (int, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return 0, gin.Error{}
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
		return 0, gin.Error{}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userIDFloat, exists := claims["user_id"]; exists {
			if userID, ok := userIDFloat.(float64); ok {
				return int(userID), nil
			}
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
	return 0, gin.Error{}
}

func (h *Handler) GetAPIUsageCost(c *gin.Context) {
	userID, err := h.getCurrentUser(c)
	if err != nil {
		return
	}

	result, err := h.analyticsService.GetAPIUsageCost(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetMonthlyUsage(c *gin.Context) {
	userID, err := h.getCurrentUser(c)
	if err != nil {
		return
	}

	result, err := h.analyticsService.GetMonthlyUsage(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) GetOverallStatus(c *gin.Context) {
	userID, err := h.getCurrentUser(c)
	if err != nil {
		return
	}

	result, err := h.analyticsService.GetOverallStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
