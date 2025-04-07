package db

import (
	"TinyBase/utils"
	"database/sql"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func Connect(host string, port int, username, password, database string) (*sql.DB, error) {
	const dsn = "user=%s password=%s host=%s port=%d dbname=%s sslmode=disable"
	db, err := sql.Open("postgres", fmt.Sprintf(dsn, username, password, host, port, database))
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Connected to database successfully")
	return db, nil
}

func Bootstrap(conn *sql.DB, su string, salt string) error {
	query := `
		SET TIME ZONE 'UTC';
		CREATE TABLE IF NOT EXISTS users (
		  	id UUID PRIMARY KEY,
		  	email TEXT UNIQUE NOT NULL,
		  	phone TEXT,
		  	name TEXT,
		  	password_hash TEXT,
		  	role TEXT,
		  	-- Email verification
		  	email_verified BOOLEAN DEFAULT FALSE,
		  	email_verification_token TEXT,
		  	email_verification_sent_at TIMESTAMP WITH TIME ZONE,
		  	-- OTP
		  	otp_enabled BOOLEAN DEFAULT FALSE,
		  	otp_secret TEXT,
		  	otp_last_used  TIMESTAMP WITH TIME ZONE,
		  	
		  	-- JWT
		  	last_login  TIMESTAMP WITH TIME ZONE,
		  	-- OAuth
		  	oauth_provider TEXT,
		  	oauth_id TEXT,
		  	
		  	-- FCM
		  	fcm_token TEXT,
		  	
		  	-- Metadata
		  	created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		    deleted_at  TIMESTAMP WITH TIME ZONE DEFAULT NULL
		);
	`

	_, err := conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// If users table is empty, insert a superuser
	var count int
	err = conn.QueryRow("SELECT COUNT(id) FROM users").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}
	if count == 0 {
		randomPassword := utils.GenerateHash()
		hashed, err := utils.HashPassword(randomPassword, salt)
		red := color.New(color.FgHiRed).SprintFunc()
		fmt.Printf("⚠ Super User Password: %s\n", red(randomPassword))

		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		_, err = conn.Exec(`
			INSERT INTO users (id, email, name, role, password_hash, email_verified, email_verification_sent_at) VALUES ($1, $2, 'TinyBase Admin' , 'SuperUser', $3, true, now())
		`, uuid.NewString(), su, hashed)
		if err != nil {
			return fmt.Errorf("failed to insert superuser: %w", err)
		}
	}
	return nil
}
