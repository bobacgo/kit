package security

import (
	"errors"
	"unicode"
)

// 定义哨兵错误
var (
	ErrPasswordTooShort  = errors.New("password must be at least the minimum required length")
	ErrPasswordNoUpper   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoDigit   = errors.New("password must contain at least one digit")
	ErrPasswordNoSpecial = errors.New("password must contain at least one special character")
)

// PasswordStrength 定义密码强度等级
type PasswordStrength int

const (
	Weak       PasswordStrength = iota // 弱
	Moderate                           // 中等
	Strong                             // 强
	VeryStrong                         // 非常强
)

// PasswordValidator 用于校验密码强度
type PasswordValidator struct {
	MinLength      int  // 最小长度
	RequireUpper   bool // 需要大写字母
	RequireLower   bool // 需要小写字母
	RequireDigit   bool // 需要数字
	RequireSpecial bool // 需要特殊字符
}

// NewPasswordValidator 创建一个新的密码校验器
func NewPasswordValidator(minLength int, requireUpper, requireLower, requireDigit, requireSpecial bool) *PasswordValidator {
	return &PasswordValidator{
		MinLength:      minLength,
		RequireUpper:   requireUpper,
		RequireLower:   requireLower,
		RequireDigit:   requireDigit,
		RequireSpecial: requireSpecial,
	}
}

// Validate 校验密码强度
func (v *PasswordValidator) Validate(password string) (PasswordStrength, error) {
	// 检查密码长度
	if len(password) < v.MinLength {
		return Weak, ErrPasswordTooShort
	}

	// 初始化校验标志
	var hasUpper, hasLower, hasDigit, hasSpecial bool

	// 遍历密码字符
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case isSpecialCharacter(char):
			hasSpecial = true
		}
	}

	// 检查规则
	if v.RequireUpper && !hasUpper {
		return Weak, ErrPasswordNoUpper
	}
	if v.RequireLower && !hasLower {
		return Weak, ErrPasswordNoLower
	}
	if v.RequireDigit && !hasDigit {
		return Weak, ErrPasswordNoDigit
	}
	if v.RequireSpecial && !hasSpecial {
		return Weak, ErrPasswordNoSpecial
	}

	// 根据满足的规则数量确定密码强度
	score := 0
	if hasUpper {
		score++
	}
	if hasLower {
		score++
	}
	if hasDigit {
		score++
	}
	if hasSpecial {
		score++
	}

	switch {
	case score >= 4:
		return VeryStrong, nil
	case score == 3:
		return Strong, nil
	case score == 2:
		return Moderate, nil
	default:
		return Weak, nil
	}
}

// isSpecialCharacter 检查字符是否为特殊字符
func isSpecialCharacter(char rune) bool {
	specialChars := `!@#$%^&*()_+-=[]{};':",./<>?\|`
	for _, c := range specialChars {
		if char == c {
			return true
		}
	}
	return false
}