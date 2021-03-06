package idutil

import (
	"crypto/rand"
	"strings"

	"github.com/sony/sonyflake"
	"github.com/speps/go-hashids/v2"

	"github.com/eachinchung/component-base/utils/iputil"
	"github.com/eachinchung/component-base/utils/stringutil"
)

// Alphabet62 字母表
//goland:noinspection SpellCheckingInspection
const (
	Alphabet36 = "abcdefghijklmnopqrstuvwxyz1234567890"
	Alphabet62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings
	st.MachineID = func() (uint16, error) {
		ip := iputil.GetLocalIP()

		return uint16([]byte(ip)[2])<<8 + uint16([]byte(ip)[3]), nil
	}

	sf = sonyflake.NewSonyflake(st)
}

// GenUint64ID 返回 uint64 的唯一 ID。
func GenUint64ID() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}

	return id
}

// GetInstanceID returns id format like: test-z8mv3z4nqw57
func GetInstanceID(uid uint64, prefix string) string {
	hd := hashids.NewData()
	hd.Alphabet = Alphabet36
	hd.MinLength = 6
	hd.Salt = "w11y16w"

	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}

	i, err := h.Encode([]int{int(uid)})
	if err != nil {
		panic(err)
	}

	var builder strings.Builder
	if prefix != "" {
		builder.WriteString(prefix)
		builder.WriteString("-")
	}
	builder.WriteString(stringutil.Reverse(i))
	return builder.String()
}

// GenSecretID 返回SecretID。
func GenSecretID() string {
	return randString(Alphabet62, 32)
}

// GenSecretKey 返回SecretKey或密码。
func GenSecretKey() string {
	return randString(Alphabet62, 36)
}

func randString(letters string, n int) string {
	output := make([]byte, n)

	// 我们将取n个字节，每个输出字符一个字节。
	randomness := make([]byte, n)

	// 随机读取所有
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}

	l := len(letters)
	// 填充输出
	for pos := range output {
		// 获取随机 item
		random := randomness[pos]

		// random % 64
		randomPos := random % uint8(l)

		// put into output
		output[pos] = letters[randomPos]
	}

	return string(output)
}
