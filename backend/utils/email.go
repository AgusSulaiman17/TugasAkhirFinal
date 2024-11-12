package utils

import (
    "gopkg.in/gomail.v2"
    "fmt"
    "os"
)

func SendEmail(recipient string, subject string, body string) error {
    // Ambil konfigurasi dari variabel lingkungan
    from := os.Getenv("EMAIL_FROM")
    password := os.Getenv("EMAIL_PASSWORD")
    smtpServer := "smtp.gmail.com"
    smtpPort := 587

    // Pastikan variabel lingkungan sudah diatur
    if from == "" || password == "" {
        return fmt.Errorf("email credentials not set in environment variables")
    }

    // Membuat pesan email
    m := gomail.NewMessage()
    m.SetHeader("From", from)
    m.SetHeader("To", recipient)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)

    // Membuat dialer untuk mengirim email
    d := gomail.NewDialer(smtpServer, smtpPort, from, password)

    // Mengirim email
    if err := d.DialAndSend(m); err != nil {
        fmt.Println("Failed to send email:", err)
        return err
    }

    fmt.Println("Email sent successfully")
    return nil
}
