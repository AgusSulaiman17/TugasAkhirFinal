package models

import (
	"time"

	"gorm.io/gorm"
)

type Peminjaman struct {
	IDPeminjaman int       `gorm:"primaryKey" json:"id_peminjaman"`
	IDUser       int       `gorm:"not null" json:"id_user"`
	IDBuku       int       `gorm:"not null" json:"id_buku"`
	TanggalPinjam time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"tanggal_pinjam"`
	DurasiHari   int       `gorm:"not null" json:"durasi_hari"`
	JamKembali   string `json:"jam_kembali"`
	TanggalKembali time.Time `json:"tanggal_kembali"`
	StatusKembali bool      `gorm:"default:false" json:"status_kembali"`
	Denda          int       `gorm:"default:0" json:"denda"`
	DibuatPada   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"dibuat_pada"`
	DiperbaruiPada time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"diperbarui_pada"`
}

func (p *Peminjaman) BeforeCreate(tx *gorm.DB) (err error) {
	// Set tanggal kembali otomatis berdasarkan durasi pinjam
	p.TanggalKembali = p.TanggalPinjam.Add(time.Duration(p.DurasiHari) * time.Hour * 24)
	return
}
