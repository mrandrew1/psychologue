package entity

import "time"

type PsySession struct {
	ID                   uint64    `gorm:"primary_key:auto_increment" json:"id"`
	ClientID             uint64    `gorm:"not null" json:"-"`
	Client               User      `gorm:"foreignkey:ClientID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
	SpecialistID         uint64    `gorm:"not null" json:"-"`
	Specialist           User      `gorm:"foreignkey:SpecialistID;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"-"`
	DateTime             string    `gorm:"type:varchar(255)" json:"datetime"`
	Notes_for_specialist string    `gorm:"type:string" json:"notes_for_specialist"`
	Notes_for_client     string    `gorm:"type:string" json:"notes_for_client"`
	CreatedAt            time.Time `json:"createdAt"`
}
