package proquint

import (
	"fmt"
	"strings"
)

type decodingConfig struct {
	finalZeroBytePadding bool
	finalHyphenPadding   bool
}

type DecodingOption func(*decodingConfig)

// WithFinalZeroBytePadding treats a final 0x00 byte as padding
// and therefore removes it from the returned value.
func WithFinalZeroBytePadding() DecodingOption {
	return func(cfg *decodingConfig) {
		cfg.finalZeroBytePadding = true
	}
}

// WithFinalHyphenPadding treats a final hyphen as indicator, that
// a final 0x00 byte is a padding byte and therefore removes it from
// the returned value.
func WithFinalHyphenPadding() DecodingOption {
	return func(cfg *decodingConfig) {
		cfg.finalHyphenPadding = true
	}
}

// ToBytes decodes a proquint string to a slice of bytes.
func ToBytes(in string, opts ...DecodingOption) ([]byte, error) {
	cfg := decodingConfig{}

	for _, opt := range opts {
		opt(&cfg)
	}

	hasFinalHyphen := strings.HasSuffix(in, "-")
	in = strings.ToLower(strings.ReplaceAll(in, "-", ""))

	if len(in)%5 != 0 {
		return nil, fmt.Errorf("invalid proquint, length not multiple of 5")
	}

	res := make([]byte, 0, len(in)/5)

	for i := 0; i < len(in)/5; i++ {
		ui16, err := ToUint16(in[i*5 : (i+1)*5])
		if err != nil {
			return nil, err
		}

		res = append(res, byte(ui16>>8))
		res = append(res, byte(ui16))
	}

	finalByte := res[len(res)-1]
	isFinalHyphenPadding := cfg.finalHyphenPadding && hasFinalHyphen && finalByte == 0
	isZeroBytePadding := cfg.finalZeroBytePadding && finalByte == 0
	if isFinalHyphenPadding || isZeroBytePadding {
		// Strip final byte, since it is 0x00 and padding is enabled.
		res = res[:len(res)-1]
	}

	return res, nil
}

// ToUint16 decodes a proquint syllable to uint16.
func ToUint16(in string) (uint16, error) {
	if len(in) != 5 {
		return 0, fmt.Errorf("invalid quint %q does not have 5 characters", in)
	}

	var res uint16
	for i, letter := range []byte(in) {
		table := consonants
		if i%2 == 1 {
			table = vowel
		}

		ui16, err := indexOf(letter, table)
		if err != nil {
			return 0, err
		}

		if i != 0 {
			shift := 2 + 2*((i+1)%2)
			res <<= shift
		}

		res += ui16
	}

	return res, nil
}

func indexOf(letter byte, table []byte) (uint16, error) {
	for i, c := range table {
		if c == letter {
			return uint16(i), nil
		}
	}

	return 0, fmt.Errorf("invalid letter %q in quint", string([]byte{letter}))
}

// ToInt16 decodes a proquint syllable to int16.
func ToInt16(in string) (int16, error) {
	ui16, err := ToUint16(in)
	return int16(ui16), err
}

// ToUint32 decodes two proquint syllables to uint32.
func ToUint32(in string) (uint32, error) {
	quints := strings.Split(in, "-")
	if len(quints) != 2 {
		return 0, fmt.Errorf("invalid input, expect 2 quints, got %d", len(quints))
	}

	var res uint32

	for i, quint := range quints {
		ui16, err := ToUint16(quint)
		if err != nil {
			return 0, err
		}

		if i != 0 {
			res <<= 16
		}

		res += uint32(ui16)
	}

	return res, nil
}

// ToInt32 decodes two proquint syllables to int32.
func ToInt32(in string) (int32, error) {
	ui32, err := ToUint32(in)
	return int32(ui32), err
}

// ToUint64 decodes four proquint syllables to uint64.
func ToUint64(in string) (uint64, error) {
	quints := strings.Split(in, "-")
	if len(quints) != 4 {
		return 0, fmt.Errorf("invalid input, expect 4 quints, got %d", len(quints))
	}

	var res uint64

	for i, quint := range quints {
		ui16, err := ToUint16(quint)
		if err != nil {
			return 0, err
		}

		if i != 0 {
			res <<= 16
		}

		res += uint64(ui16)
	}

	return res, nil
}

// ToInt64 decodes four proquint syllables to int64.
func ToInt64(in string) (int64, error) {
	ui64, err := ToUint64(in)
	return int64(ui64), err
}
