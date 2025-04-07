package models

type User struct {
	ID                      string `json:"id"`
	Email                   string `json:"email"`
	Phone                   string `json:"phone"`
	Name                    string `json:"name"`
	PasswordHash            string `json:"password_hash"`
	Role                    string `json:"role"`
	EmailVerified           bool   `json:"email_verified"`
	EmailVerificationToken  string `json:"email_verification_token"`
	EmailVerificationSentAt string `json:"email_verification_sent_at"`
	OtpEnabled              bool   `json:"otp_enabled"`
	OtpSecret               string `json:"otp_secret"`
	OtpLastUsed             string `json:"otp_last_used"`
	LastLogin               string `json:"last_login"`
	OauthProvider           string `json:"oauth_provider"`
	OauthID                 string `json:"oauth_id"`
	FcmToken                string `json:"fcm_token"`
	CreatedAt               string `json:"created_at"`
	UpdatedAt               string `json:"updated_at"`
	DeletedAt               string `json:"deleted_at"`
}
