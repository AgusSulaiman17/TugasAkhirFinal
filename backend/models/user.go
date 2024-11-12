package models

import (
    "time"
)

type User struct {
    IDUser        int       `json:"id_user" gorm:"primary_key"` // Menandai ID sebagai primary key
    Nama          string    `json:"nama"`                       // Nama pengguna
    Email         string    `json:"email" gorm:"unique"`       // Email pengguna harus unik
    KataSandi     string    `json:"kata_sandi"`                 // Kata sandi (harus di-hash)
    Role          string    `json:"role"`                       // Role pengguna
    DibuatPada    time.Time `json:"dibuat_pada" gorm:"autoCreateTime"` // Waktu dibuat
    DiperbaruiPada time.Time `json:"diperbarui_pada" gorm:"autoUpdateTime"` // Waktu diperbarui
}

