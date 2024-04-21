package jwtTool

import (
	"testing"
	"time"
)

const MaxC = 100000

func Test(t *testing.T) {
	driver := NewToken()
	tokens := make(map[string]string)
	for i := 0; i < MaxC; i++ {
		uid := string(randByte(10))
		token := driver.Generate(uid, 3*time.Second)
		tokens[token] = uid
	}
	for k, _ := range tokens {
		_, err := driver.Parse(k)
		if err != nil {
			continue
		}
	}
}
