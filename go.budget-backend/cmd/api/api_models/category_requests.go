package api_models

type CategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"max=255"`
	IsCustom    bool   `json:"is_custom" default:"true"`
}
