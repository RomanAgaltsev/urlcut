package random

import "math/rand"

const Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// String формирует последовательность английских больших, маленьких букв и цифр заданной длины.
func String(lenght int) string {

	b := make([]byte, lenght)
	for i := range b {
		b[i] = Letters[rand.Int63()%int64(len(Letters))]
	}
	return string(b)
}
