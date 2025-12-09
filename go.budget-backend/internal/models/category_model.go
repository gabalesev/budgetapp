package models

type CategoryModel struct {
	BaseModel
	Slug     string `gorm:"unique;type:varchar(200);not null" json:"slug"`
	Name     string `gorm:"unique;type:varchar(200);not null" json:"name"`
	IsCustom bool   `gorm:"type:bool;default:false" json:"is_custom"`
}

func (CategoryModel) TableName() string {
	return "categories"
}
