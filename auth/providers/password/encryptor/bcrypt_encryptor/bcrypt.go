package bcrypt_encryptor

import "golang.org/x/crypto/bcrypt"

type Config struct {
	Cost int //加密强度
}

type BcryptEncryptor struct {
	Config *Config
}

func New(config *Config) *BcryptEncryptor{
	if config == nil {
		config = &Config{}
	}

	if config.Cost == 0  {
		config.Cost = bcrypt.DefaultCost //10
	}

	return&BcryptEncryptor{
		config,
	}
}

//Digest generate encrypted password
func (bcryptEncryptor *BcryptEncryptor) Digest(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcryptEncryptor.Config.Cost)
	return string(hashedPassword), err
}

//Compare check hashed password
func(bcryptEncryptor *BcryptEncryptor) Compare(hashedPassword string, password string) error{
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}