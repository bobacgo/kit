package util

import "golang.org/x/crypto/bcrypt"

type Bcrypt struct {
	Salt string
}

// NewBcrypt
// salt 添加杂质
func NewBcrypt(salt string) Bcrypt {
	return Bcrypt{
		Salt: salt,
	}
}

// Hash 明文加密
func (b Bcrypt) Hash(passwd string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(passwd+b.Salt), bcrypt.DefaultCost)
	return string(bytes)
}

// Check 校验密文和明文
func (b Bcrypt) Check(hash, passwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwd+b.Salt))
	return err == nil
}
