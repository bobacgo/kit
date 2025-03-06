package tag

import (
	"fmt"
	"testing"
)

func TestMaskTag(t *testing.T) {
	type Address struct {
		City    string `json:"city"`
		ZipCode string `json:"zip_code" mask:""` // 使用默认规则
	}

	type Settings struct {
		NotificationsEnabled bool   `json:"notifications_enabled"`
		Theme                string `json:"theme" mask:""` // 使用默认规则
	}
	type User struct {
		ID       int                `json:"id"`
		Name     string             `json:"name"`
		Email    string             `json:"email" mask:""`        // 使用默认规则
		Password string             `json:"password" mask:"^.*$"` // 使用正则表达式
		Phone    string             `json:"phone"`                // 没有 mask 标签，不脱敏
		Address  map[string]Address `json:"address"`              // 嵌套结构体
		Settings Settings           `json:"settings"`             // 嵌套结构体
		Emails   []string           `json:"emails" mask:""`       // Slice 类型
		Aliases  [2]string          `json:"aliases" mask:""`      // Array 类型
		Metadata map[string]string  `json:"metadata" mask:""`     // Map 类型

	}

	user := User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john.doe@example.com",
		Password: "s3cr3tP@ssw0rd",
		Phone:    "13800138000", // 没有 mask 标签，不脱敏
		Address: map[string]Address{
			"home": {City: "New York", ZipCode: "10001"},
		},
		Settings: Settings{
			NotificationsEnabled: true,
			Theme:                "dark",
		},
		Emails: []string{
			"john.doe@example.com",
			"jane.doe@example.com",
		},
		Aliases: [2]string{
			"johnny",
			"johnd",
		},
		Metadata: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
	}

	fmt.Println("Before desensitization:")
	fmt.Printf("%+v\n", user)

	// 返回一个新的脱敏对象
	maskedUser := Desensitize(user)

	fmt.Println("After desensitization:")
	fmt.Printf("%+v\n", maskedUser)

	// 原始对象未被修改
	fmt.Println("Original object (unchanged):")
	fmt.Printf("%+v\n", user)
}