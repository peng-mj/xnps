package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"regexp"
	"sort"
	"time"
)

const SnowStartTime = 1700000000000

var entropy *ulid.MonotonicEntropy

func init() {
	snowflake.Epoch = SnowStartTime
	entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

}

// en
func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// de
func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	err, origData = PKCS5UnPadding(origData)
	return origData, err
}

// Completion when the length is insufficient
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// Remove excess
func PKCS5UnPadding(origData []byte) (error, []byte) {
	length := len(origData)
	unPadding := int(origData[length-1])
	if (length - unPadding) < 0 {
		return errors.New("len error"), nil
	}
	return nil, origData[:(length - unPadding)]
}

// Generate 256-bit sha256 strings
func Sha256(s string) string {
	sha := sha256.New()
	sha.Write([]byte(s))
	return hex.EncodeToString(sha.Sum(nil))
}
func Sha1(s string) string {
	sha := sha1.New()
	sha.Write([]byte(s))
	return hex.EncodeToString(sha.Sum(nil))
}

func Md5(s []byte) string {
	h := md5.New()
	h.Write(s)
	return hex.EncodeToString(h.Sum(nil))
}

type RandString struct {
	base []byte
}

func (r *RandString) AddAll() *RandString {
	r.base = []byte("abcdefghijklmnopqrstuvwxyzQWERTYUIOPASDFGHJKLZXCVBNM-_=#$|?/.,;:+*@~")
	return r
}

func (r *RandString) AddNum() *RandString {
	r.base = append(r.base, []byte("0123456789")...)

	return r
}
func (r *RandString) AddLetter() *RandString {
	r.base = append(r.base, []byte("abcdefghijklmnopqrstuvwxyzQWERTYUIOPASDFGHJKLZXCVBNM")...)
	return r
}
func (r *RandString) AddSymbol() *RandString {
	r.base = append(r.base, []byte("-_=#$|?/.,;:+*@~")...)
	return r
}
func (r *RandString) SetBase(base []byte) *RandString {
	r.base = base
	return r
}
func (r *RandString) Generate(l int) string {
	if len(r.base) == 0 {
		r.AddNum().AddLetter()
	}
	n := len(r.base)
	res := make([]byte, l)
	l--
	for l >= 0 {
		res[l] = r.base[rand.Intn(n)]
		l--
	}
	return string(res)
}
func (r *RandString) GenerateList(length, num int) []string {
	res := make([]string, num)
	for i := range res {
		res[i] = r.Generate(length)
	}
	return res
}

func RandStr() *RandString {
	return &RandString{}
}

func RandVKey() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String() + RandStr().Generate(6)
}
func SnowID(workerID int64) int64 {
	node, _ := snowflake.NewNode(workerID)
	return node.Generate().Int64()
}
func Ulid() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
func Uuid() string {
	return uuid.New().String()
}

func CheckPassed(passwd string) int {
	bt := []byte(passwd)
	var num, letLow, letUp, other int
	for i := range bt {
		if bt[i] >= '0' && bt[i] <= '9' {
			num++
		} else if bt[i] >= 'a' && bt[i] <= 'z' {
			letLow++
		} else if bt[i] >= 'A' && bt[i] <= 'Z' {
			letUp++
		} else {
			other++
		}
	}
	res := 0
	if num != 0 {
		res++
	}
	if letLow != 0 {
		res++
	}
	if letUp != 0 {
		res++
	}
	if other != 0 {
		res++
	}
	return res
}

func CheckEmail(email string) bool {
	match, err := regexp.Match(`(\w+)@(\w{2,}).([.a-z]+)`, []byte(email))
	if err != nil {
		return false
	}
	return match
}

func MergeRanges(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return [][]int{}
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	var merged [][]int
	cur := intervals[0]
	for i := 0; i < len(intervals); i++ {
		if cur[1]+1 >= intervals[i][0] {
			cur[1] = intervals[i][1]
		} else {
			merged = append(merged, cur)
			cur = intervals[i]
		}
	}
	// 将最后一个范围加入结果集合
	merged = append(merged, cur)

	return merged
}
