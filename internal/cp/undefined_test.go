package cp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCharsetToUTF8_UndefinedBytes(t *testing.T) {
	// A byte the code page leaves undefined decodes to what SQL Server returns for it (checked against SQL Server
	// 2022), not U+FFFD: a C1 control from 0x80 to 0x9F, a private use character above. Encoding the value back to
	// the code page then gives the same byte.
	tests := []struct {
		name     string
		c        Collation
		input    []byte
		expected string
	}{
		{"cp1252 0x81", Collation{SortId: 50}, []byte("Gonz\xe3\x81Lez"), "Gonzã\u0081Lez"},
		{"cp1252 0x8D 0x8F 0x90 0x9D", Collation{SortId: 50}, []byte{0x8D, 0x8F, 0x90, 0x9D}, "\u008D\u008F\u0090\u009D"},
		{"cp1253 0xAA", Collation{SortId: 112}, []byte{0xAA}, ""},
		{"cp1255 0xCA", Collation{SortId: 136}, []byte{0xCA}, "ֺ"},
		{"cp874 0xDB", Collation{SortId: 204}, []byte{0xDB}, ""},
		{"cp932 0xA0", Collation{SortId: 192}, []byte{0xA0}, ""},
		{"cp949 0x80", Collation{SortId: 194}, []byte{0x80}, "\u0080"},
	}
	for _, tc := range tests {
		assert.Equal(t, tc.expected, CharsetToUTF8(tc.c, tc.input), tc.name)
	}

	for _, cm := range []*charsetMap{
		getcp437(), getcp850(), getcp874(), getcp932(), getcp936(), getcp949(), getcp950(), getcp1250(),
		getcp1251(), getcp1252(), getcp1253(), getcp1254(), getcp1255(), getcp1256(), getcp1257(), getcp1258(),
	} {
		assert.NotContains(t, cm.sb[:], rune(0xFFFD), "a single byte still decodes as U+FFFD")
	}
}
