// models/buku.go
package models

import (
    "time"
)

type Buku struct {
    ID          int            `gorm:"column:id_buku;primaryKey"`
    Judul       string         `gorm:"type:varchar(100);not null"`
    IdPenulis   int            `gorm:"column:id_penulis"`
    IdGenre     int            `gorm:"column:id_genre"`
    Deskripsi   string         `gorm:"type:text"`
    Jumlah      int            `gorm:"not null"`
    Gambar      string         `gorm:"type:varchar(255)"`
    Status      bool           `gorm:"default:true"`
    DibuatPada  time.Time      `gorm:"column:dibuat_pada;default:CURRENT_TIMESTAMP"`
    DiperbaruiPada time.Time   `gorm:"column:diperbarui_pada;default:CURRENT_TIMESTAMP"`
}

// TableName memberi tahu GORM untuk menggunakan nama tabel 'books'
func (Buku) TableName() string {
    return "books"
}
