package models

type Genre struct {
    IDGenre int    `gorm:"primaryKey" json:"id_genre"`
    Nama    string `json:"nama"`
}
