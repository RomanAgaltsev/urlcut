package random

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	const (
		stringLenght = 8
		badLetters   = "абвгдеёжзиклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ/*-+!№;%:?*()"
	)
	randomString := String(stringLenght)

	assert.Len(t, randomString, stringLenght)
	assert.True(t, strings.ContainsAny(randomString, Letters))
	assert.False(t, strings.ContainsAny(randomString, badLetters))
}

func BenchmarkString(b *testing.B) {
	const lenght = 8
	for range b.N {
		_ = String(lenght)
	}
}
