package crypt

import (
	"testing"
)

func Test(t *testing.T) {

	for i := 0; i < 2000000; i++ {
		SnowID(int64(i % 1023))
	}

}
