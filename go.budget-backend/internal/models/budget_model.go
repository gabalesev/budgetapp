package models

import "time"

type BudgetModel struct {
	BaseModel
	Title       string           `gorm:"index;type:varchar(255);not null" json:"title"`
	Slug        string           `gorm:"index;type:varchar(255);not null;uniqueIndex:unique_user_id_slug_year_month" json:"slug"`
	UserID      uint             `gorm:"column:user_id;uniqueIndex:unique_user_id_slug_year_month;not null" json:"user_id"`
	Description *string          `gorm:"type:text" json:"description"`
	Amount      float64          `gorm:"type:decimal(10,2);not null" json:"amount"`
	Categories  []*CategoryModel `gorm:"constraints:OnDelete:CASCADE;many2many:budget_categories;joinForeignKey:"`
	Date        time.Time        `gorm:"type:datetime;not null" json:"date"`
	Month       uint             `gorm:"type:TINYINT UNSIGNED;not null;index:idx_month_year;uniqueIndex:unique_user_id_slug_year_month" json:"month"`
	Year        uint16           `gorm:"type:INT UNSIGNED;not null;index:idx_month_year;uniqueIndex:unique_user_id_slug_year_month" json:"year"`
}

func (BudgetModel) TableName() string {
	return "budgets"
}
