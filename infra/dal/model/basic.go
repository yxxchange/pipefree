package model

type Basic struct {
	Id        int64 `json:"id" gorm:"column:id"`
	CreatedAt int64 `json:"created_at" gorm:"column:created_at"`
	UpdatedAt int64 `json:"updated_at" gorm:"column:updated_at"`
	IsDel     bool  `json:"is_del" gorm:"column:is_del"` // soft delete flag
}
