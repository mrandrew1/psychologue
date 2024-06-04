package dto

type PsySessionUpdateDTO struct {
	ID          uint64 `json:"id" form:"id" binding:"required"`
	Title       string `json:"title" form:"title" binding:"required"`
	Description string `json:"description" form:"description" binding:"required"`
	UserID      uint64 `json:"user_id,omitempty" form:"user_id,omitempty"`
}

type PsySessionCreateDTO struct {
	ClientID             uint64 `gorm:"not null" json:"client_id"`
	SpecialistID         uint64 `gorm:"not null" json:"specialist_id"`
	DateTime             string `gorm:"type:varchar(255)" json:"datetime"`
	Notes_for_specialist string `gorm:"type:string" json:"notes_for_specialist"`
	Notes_for_client     string `gorm:"type:string" json:"notes_for_client"`
}
