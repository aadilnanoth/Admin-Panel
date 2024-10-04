package models

import (
	"time"
)

// User struct defines the schema for the users table
type User struct {
	ID           int       `json:"id"`                              // Primary Key (auto-increment)
	FirstName    string    `json:"first_name" binding:"required"`   // User's first name (required)
	LastName     string    `json:"last_name" binding:"required"`    // User's last name (required)
	Email        string    `json:"email" binding:"required,email"`  // Email (required, must be a valid email)
	Password     string    `json:"-"`                               // Password (hashed, not included in JSON response)
	PhoneNumber  string    `json:"phone_number" binding:"required"` // Phone number (used for OTP verification)
	Status       string    `json:"status"`                          // Account status (e.g., "pending", "active", "blocked")
	OTPCode      string    `json:"otp_code"`                        // OTP code for verification
	OTPExpiresAt time.Time `json:"otp_expires_at"`                  // Expiration time of the OTP
	CreatedAt    time.Time `json:"created_at"`                      // Timestamp of user creation
	UpdatedAt    time.Time `json:"updated_at"`                      // Timestamp of last update
}
