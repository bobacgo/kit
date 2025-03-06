package security

import (
	"fmt"
	"testing"
)

func TestPwd(t *testing.T) {
	pv := new(PasswdVerifier)
	hash := pv.BcryptHash("admin123")
	t.Log(hash)
	if pv.BcryptVerify(hash, "admin123") {
		t.Log("success")
	} else {
		t.Log("fail")
	}
}

func TestPwdStrength(t *testing.T) {
	// 创建一个密码校验器
	validator := NewPasswordValidator(8, true, true, true, true)

	// 测试密码
	passwords := []string{
		"weak",           // 太短
		"weak123",        // 缺少大写字母
		"Strong123",      // 缺少特殊字符
		"VeryStrong123!", // 符合所有规则
	}

	// strengthToString 将密码强度转换为字符串
	var strengthToString = func(strength PasswordStrength) string {
		switch strength {
		case Weak:
			return "Weak"
		case Moderate:
			return "Moderate"
		case Strong:
			return "Strong"
		case VeryStrong:
			return "Very Strong"
		default:
			return "Unknown"
		}
	}

	for _, password := range passwords {
		strength, err := validator.Validate(password)
		if err != nil {
			switch err {
			case ErrPasswordTooShort:
				fmt.Printf("Password: %s, Error: %v\n", password, err)
			case ErrPasswordNoUpper:
				fmt.Printf("Password: %s, Error: %v\n", password, err)
			case ErrPasswordNoLower:
				fmt.Printf("Password: %s, Error: %v\n", password, err)
			case ErrPasswordNoDigit:
				fmt.Printf("Password: %s, Error: %v\n", password, err)
			case ErrPasswordNoSpecial:
				fmt.Printf("Password: %s, Error: %v\n", password, err)
			default:
				fmt.Printf("Password: %s, Unknown Error: %v\n", password, err)
			}
		} else {
			fmt.Printf("Password: %s, Strength: %v\n", password, strengthToString(strength))
		}
	}
}