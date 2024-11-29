package models


type Currency struct {
    ID      uint    `gorm:"primaryKey;autoIncrement" json:"id"`
    Currency string  `gorm:"type:varchar(3);not null" json:"currency"`
    Rate    float64 `gorm:"type:decimal(10,4);not null" json:"rate"`
}
