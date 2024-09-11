package encrypt

import "golang.org/x/crypto/bcrypt"

// HashPassword 加密密码
func HashPassword(password string) (string, error) {
	// 使用 bcrypt 加密密码，第二个参数是成本参数，推荐值是 14
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPassword 验证密码
func CheckPassword(password string, hashedPassword string) bool {
	// 验证输入的密码是否与加密后的密码匹配
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
