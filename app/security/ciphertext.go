package security

import "github.com/bobacgo/kit/pkg/ucrypto"

// Ciphertext 密文
// use: 前端密码字段的传输
//
// 密码字段设计:
// 1.前端密码字段加密
// 2.后端解密出原文
// 3.后端密码强度校验
// 4.入库时hash不可逆编码(可以加盐)
type Ciphertext string

// Encrypt 加密 (可以反解密)
func (ct *Ciphertext) Encrypt(secret string) error {
	pt, err := ucrypto.AESEncrypt(string(*ct), secret)
	if err != nil {
		return err
	}
	*ct = Ciphertext(pt)
	return nil
}

// Decrypt 解密 (解出来的是原文)
func (ct *Ciphertext) Decrypt(secret string) error {
	pt, err := ucrypto.AESDecrypt(string(*ct), secret)
	if err != nil {
		return err
	}
	*ct = Ciphertext(pt)
	return nil
}

// BcryptHash 密码加密
func (ct *Ciphertext) BcryptHash() string {
	h := processPwd(string(*ct))
	*ct = Ciphertext(h)
	return h
}

// BcryptVerify 验证密码
func (ct *Ciphertext) BcryptVerify(hashPasswd string) bool {
	return verifyPwd(hashPasswd, string(*ct))
}
