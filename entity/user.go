package entity

type User struct {
	ID       uint64  `gorm:"primary_key:auto_increment" json:"id"`
	Name     string  `gorm:"type:varchar(255)" json:"name"`
	Email    string  `gorm:"type:varchar(255);unique" json:"email"`
	Password string  `gorm:"->;<-;not null" json:"-"`
	Token    string  `gorm:"-" json:"token,omitempty"`
	Books    *[]Book `json:"books,omitempty"`
	//PsySession *[]PsySession `json:"psySession,omitempty"`
	//Role 	 string  `gorm:"type:varchar(255)" json:"role"`
}
