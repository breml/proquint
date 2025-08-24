package proquint

import (
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	maskConsonant uint16 = 0x000F
	maskVowel     uint16 = 0x0003

	shiftFirst  = 16 - 4
	shiftSecond = 16 - 6
	shiftThird  = 16 - 10
	shiftForth  = 16 - 12
)

// FromUint16 encodes proquint from the provided uint16 argument.
func FromUint16(in uint16) string {
	str := strings.Builder{}

	str.WriteByte(consonants[(in>>shiftFirst)&maskConsonant])
	str.WriteByte(vowel[(in>>shiftSecond)&maskVowel])
	str.WriteByte(consonants[(in>>shiftThird)&maskConsonant])
	str.WriteByte(vowel[(in>>shiftForth)&maskVowel])
	str.WriteByte(consonants[in&maskConsonant])

	return str.String()
}

// FromInt16 encodes proquint from the provided int16 argument.
func FromInt16(in int16) string {
	return FromUint16(uint16(in))
}

// FromUint32 encodes proquint from the provided uint32 argument.
func FromUint32(in uint32, opts ...EncodingOption) string {
	cfg := encodingConfig{}

	for _, opt := range opts {
		opt(&cfg)
	}

	hyphens := ""
	if cfg.hyphens {
		hyphens = "-"
	}

	return FromUint16(uint16(in>>16)) + hyphens + FromUint16(uint16(in))
}

// FromInt32 encodes proquint from the provided int32 argument.
func FromInt32(in int32, opts ...EncodingOption) string {
	return FromUint32(uint32(in), opts...)
}

// FromUint64 encodes proquint from the provided uint64 argument.
func FromUint64(in uint64, opts ...EncodingOption) string {
	cfg := encodingConfig{}

	for _, opt := range opts {
		opt(&cfg)
	}

	hyphens := ""
	if cfg.hyphens {
		hyphens = "-"
	}

	return FromUint16(uint16(in>>48)) + hyphens + FromUint16(uint16(in>>32)) + hyphens + FromUint16(uint16(in>>16)) + hyphens + FromUint16(uint16(in))
}

// FromInt64 encodes proquint from the provided int64 argument.
func FromInt64(in int64, opts ...EncodingOption) string {
	return FromUint64(uint64(in), opts...)
}

type encodingConfig struct {
	hyphens        bool
	disablePadding bool
	legacyPadding  bool
}

type EncodingOption func(*encodingConfig)

// WithHyphens adds a hyphen between each proquint syllable:
//
//	lusab-babad
func WithHyphens() EncodingOption {
	return func(cfg *encodingConfig) {
		cfg.hyphens = true
	}
}

// LegacyWithoutPadding allows to disable padding of odd number of bytes.
// Without padding, trying to encode odd number of bytes returns an error.
//
// Disabling of padding is not following draft-rayner-proquint and is
// therefore considered deprecated. This option allows for compatibility
// with the original specification by Daniel S. Wilkerson.
func LegacyWithoutPadding() EncodingOption {
	return func(cfg *encodingConfig) {
		cfg.disablePadding = true
		cfg.legacyPadding = true
	}
}

// LegacyWithZeroPadding allows encoding of odd number of bytes by adding a
// zero (0x00) padding byte to the input.
//
// Using a padding byte without the final hyphen as padding indicator is not
// following draft-rayner-proquint and is therefore considered deprecated.
// This option allows for compatibility with the original specification by
// Daniel S. Wilkerson.
func LegacyWithZeroPadding() EncodingOption {
	return func(cfg *encodingConfig) {
		cfg.legacyPadding = true
	}
}

func FromBytes(in []byte, opts ...EncodingOption) (string, error) {
	cfg := encodingConfig{}

	for _, opt := range opts {
		opt(&cfg)
	}

	padded := false
	if !cfg.disablePadding && len(in)%2 == 1 {
		// Odd number of bytes in input, compensate with 0x00 padding byte.
		in = append(in, 0x00)
		padded = true
	}

	if len(in)%2 != 0 {
		return "", fmt.Errorf("only arguments with even length are supported")
	}

	str := strings.Builder{}

	for i := 0; i < len(in); i += 2 {
		if cfg.hyphens && i > 0 {
			str.WriteByte('-')
		}

		str.WriteString(FromUint16(uint16(in[i])<<8 + uint16(in[i+1])))
	}

	if !cfg.legacyPadding && padded {
		str.WriteByte('-')
	}

	return str.String(), nil
}

func FromHexString(in string, opts ...EncodingOption) (string, error) {
	hexBytes, err := hex.DecodeString(in)
	if err != nil {
		return "", err
	}

	return FromBytes(hexBytes, opts...)
}
