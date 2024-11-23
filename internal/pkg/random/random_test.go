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
	//	assert.Contains(t, Letters, randomString)
	//	assert.NotContains(t, badLetters, randomString)
}
