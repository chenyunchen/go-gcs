package entity

// PreviewUrl is the structure for preview url response
type PreviewedUrl struct {
	Type        string `json:"type" validate:"required"`
	Icon        string `json:"icon" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"desc" validate:"required"`
	Image       string `json:"pic" validate:"required"`
	Url         string `json:"url" validate:"required"`
}
