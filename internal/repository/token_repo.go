package repository

import (
	"time"

	"github.com/RAF-SI-2025/EXBanka-3-Infrastructure/internal/models"

	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(token *models.Token) error {
	return r.db.Create(token).Error
}

// FindValid returns an unused, non-expired token of the given type
func (r *TokenRepository) FindValid(tokenStr, tokenType string) (*models.Token, error) {
	var token models.Token
	err := r.db.Preload("Employee").
		Where("token = ? AND type = ? AND used = false AND expires_at > ?",
			tokenStr, tokenType, time.Now()).
		First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepository) MarkUsed(id uint) error {
	return r.db.Model(&models.Token{}).Where("id = ?", id).Update("used", true).Error
}

// InvalidateEmployeeTokens marks all unused tokens of the given type for an employee as used
func (r *TokenRepository) InvalidateEmployeeTokens(employeeID uint, tokenType string) error {
	return r.db.Model(&models.Token{}).
		Where("employee_id = ? AND type = ? AND used = false", employeeID, tokenType).
		Update("used", true).Error
}
