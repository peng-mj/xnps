package crypt

import (
	"golang.org/x/exp/slog"
	"testing"
)

func Test(t *testing.T) {

	for i := 0; i < 2000000; i++ {
		SnowID(int64(i % 1023))
	}

}

func Test_Check(t *testing.T) {

	r := RandStr().AddLetter().AddNum()
	for i := 0; i < 20; i++ {
		passwd := r.Generate(20)
		slog.Info("password check", passwd, CheckPassed(passwd))
	}
	emials := []string{"123@1.cn", "123asd@163.com", "123@google.com", "123asd@@.com", "123s.@cd.cc"}
	for i := 0; i < len(emials); i++ {
		slog.Info("email check", emials[i], CheckEmail(emials[i]))
	}

}
