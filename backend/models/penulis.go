package models

import (
    "time"
)

type Penulis struct {
    IDPenulis   uint      `gorm:"primaryKey;column:id_penulis" json:"id_penulis"`
    Nama        string    `gorm:"column:nama" json:"nama"`
    Biografi    string    `gorm:"column:biografi" json:"biografi"`
    DibuatPada  time.Time `gorm:"column:dibuat_pada;autoCreateTime" json:"dibuat_pada"`
    DiperbaruiPada time.Time `gorm:"column:diperbarui_pada;autoUpdateTime" json:"diperbarui_pada"`
}
