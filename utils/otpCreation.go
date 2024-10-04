package utils

import (
	"crypto/rand"
	"log"
	"math/big"
	"net/smtp"
)

// GenerateOTP generates a 6-digit OTP
func GenerateOTP() (string, error) {
	otp := ""
	for i := 0; i < 6; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		otp += n.String()
	}
	return otp, nil
}

// SendOTPEmail sends an email with the OTP
func SendOTPEmail(to string, otp string) error {
	from := "aadilnanoth@gmail.com"
	password := "jzvffexxenxzbixe"
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Your OTP Code\r\n" +
		"\r\n" +
		"Your OTP code is: " + otp + "\r\n" +
		"This OTP is valid for 10 minutes.\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send OTP email: %v", err)
		return err
	}

	return nil
}
