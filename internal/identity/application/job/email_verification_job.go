package job

type EmailVerificationJob struct {
	UserID          string `json:"user_id"`
	Email           string `json:"email"`
	FullName        string `json:"full_name"`
	Token           string `json:"token"`
	VerificationURL string `json:"verification_url"`
}

const EmailVerificationJobName = "identity.email_verification.send"
