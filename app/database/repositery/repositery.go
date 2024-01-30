package repository

type ResetCodeRepository interface {
	CreateResetCode(code *ResetCode) error
	GetResetCodeByEmail(email string) (*ResetCode, error)
	DeleteResetCodeByEmail(email string) error
}

type ResetCode struct {
	Email string
	Code  string
}
