package database

import (
	"database/sql"
	"fmt"
	"log"
	"login_page/models"
	"os"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func DbConnect() {
	connstr := os.Getenv("DB_URL")
	var err error
	DB, err = sql.Open("postgres", connstr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Could not ping DB: %v", err)
	}

	fmt.Println("Successfully connected to the database")
}

// GetAllUsers returns all users from the database
func GetAllUsers() ([]models.User, error) {
	rows, err := DB.Query("SELECT id, email, role, status, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Email, &user.Status, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func CreateUser(db *sql.DB, user *models.User) error {
	// Hash the password before saving it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Set the timestamps
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// SQL query to insert the new user
	query := `
	INSERT INTO users (first_name, last_name, email, password, phone_number, status, otp_code, otp_expires_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
	`

	// Execute the query and return the new user's ID
	err = db.QueryRow(query, user.FirstName, user.LastName, user.Email, string(hashedPassword), user.PhoneNumber, user.Status, user.OTPCode, user.OTPExpiresAt, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password, status, created_at, updated_at FROM users WHERE email = $1`
	row := DB.QueryRow(query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err
	}

	return &user, nil
}

func CreateAdmin(db *sql.DB, admin *models.User) error {
	// Hash the password before saving it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Set the timestamps
	admin.CreatedAt = time.Now()
	admin.UpdatedAt = time.Now()

	// SQL query to insert the new user
	query := `
		INSERT INTO admins (first_name, last_name, email, password,  status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	// Execute the query and return the new user's ID
	err = db.QueryRow(query, admin.FirstName, admin.LastName, admin.Email, string(hashedPassword), admin.Status, admin.CreatedAt, admin.UpdatedAt).Scan(&admin.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetAdminByEmail(email string) (*models.User, error) {
	query := `SELECT id, email, password,  status, created_at, updated_at FROM admins WHERE email = $1`
	row := DB.QueryRow(query, email)

	var admin models.User
	err := row.Scan(&admin.ID, &admin.Email, &admin.Password, &admin.Status, &admin.CreatedAt, &admin.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No user found
		}
		return nil, err
	}

	return &admin, nil
}

// UpdateUser updates the user's status in the database
func UpdateUser(db *sql.DB, user *models.User) error {
	query := `
        UPDATE users 
        SET status = ?, otp_code = ?, otp_expires_at = ?, updated_at = NOW() 
        WHERE email = ?`

	_, err := db.Exec(query, user.Status, user.OTPCode, user.OTPExpiresAt, user.Email)
	return err

}
