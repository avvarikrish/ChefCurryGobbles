package hash

import "golang.org/x/crypto/bcrypt"

// CreateHash takes a string and returns a cryptographic hashed password
func CreateHash(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

// CompareHash takes two byte array passwords and compares them
func CompareHash(stored []byte, req []byte) bool {
	if err := bcrypt.CompareHashAndPassword(stored, req); err != nil {
		return false
	}

	return true
}
