package util

import (
	"reflect"
	"testing"
)

func TestExample(t *testing.T) {
	bcrypt := NewBcrypt("gogo")

	hash := bcrypt.Hash("Abc@123")
	check := bcrypt.Check(hash, "Abc@123")
	t.Logf("check: %v\n", check)
}

func TestBcrypt_Check(t *testing.T) {
	type fields struct {
		Salt string
	}
	type args struct {
		hash   string
		passwd string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{name: "Check-True", fields: fields{Salt: "gogo"}, args: args{hash: NewBcrypt("gogo").Hash("abc123"), passwd: "abc123"}, want: true},
		{name: "Check-False", fields: fields{Salt: "gogo"}, args: args{hash: NewBcrypt("gogo").Hash("abc123"), passwd: "abcd123"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bcrypt{
				Salt: tt.fields.Salt,
			}
			if got := b.Check(tt.args.hash, tt.args.passwd); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBcrypt(t *testing.T) {
	type args struct {
		salt string
	}
	tests := []struct {
		name string
		args args
		want Bcrypt
	}{
		{name: "创建Bcrypt", args: args{salt: "gogo"}, want: Bcrypt{Salt: "gogo"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBcrypt(tt.args.salt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBcrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
