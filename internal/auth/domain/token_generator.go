package domain

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int // in seconds
}

// TokenClaims represents the payload embedded in the token
type TokenClaims struct {
	UserID string
	Role   string
	Type   string // "access" or "refresh"
}

// TokenGenerator defines the contract for generating and validating tokens
type TokenGenerator interface {
	GenerateTokenPair(userID string, role string) (*TokenPair, error)
	ValidateToken(tokenString string, expectedType string) (*TokenClaims, error)
}
