package models

import (
	"time"
)

type BaseModeler interface {
	GetID() uint
	SetID(uint)
}

type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *BaseModel) GetID() uint {
	return b.ID
}

func (b *BaseModel) SetID(id uint) {
	b.ID = id
}
