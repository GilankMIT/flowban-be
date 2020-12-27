package stringBuilder

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateRandom(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // E.g. "ExcbsVQs"

	return str
}

func GenerateRandomCustom(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("abcdefghjklmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ" +
		"123456789")

	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String() // E.g. "ExcbsVQs"

	return str
}
