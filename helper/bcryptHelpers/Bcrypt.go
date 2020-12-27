package bcryptHelpers

import "golang.org/x/crypto/bcrypt"

//GenerateBcrypt will generate bcrypt hash with 12 rounds
func GenerateBcrypt(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

//GenerateBcryptCustomCost will generate bcrypt hash with custom rounds
func GenerateBcryptCustomCost(plainPassword string, rounds int) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), rounds)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

//CompareBcrypt will compare hash password with plain password
func CompareBcrypt(hash string, plain string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		return err
	}
	return nil

}
