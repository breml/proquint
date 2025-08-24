package proquint

import (
	"fmt"
	"strings"
)

type decodingConfig struct {
	legacyFinalZeroBytePadding bool
}

type DecodingOption func(*decodingConfig)

// LegacyWithFinalZeroBytePadding treats a final 0x00 byte as padding
// and therefore removes it from the returned value.
//
// Using a padding byte without the final hyphen as padding indicator is not
// following draft-rayner-proquint and is therefore considered deprecated.
// This option allows for compatibility with the original specification by
// Daniel S. Wilkerson.
func LegacyWithFinalZeroBytePadding() DecodingOption {
	return func(cfg *decodingConfig) {
		cfg.legacyFinalZeroBytePadding = true
	}
}

// ToBytes decodes a proquint string to a slice of bytes.
func ToBytes(in string, opts ...DecodingOption) ([]byte, error) {
	cfg := decodingConfig{}

	for _, opt := range opts {
		opt(&cfg)
	}

	if !cfg.legacyFinalZeroBytePadding && strings.Contains(in, "--") {
		return nil, fmt.Errorf("invalid proquint, consecutive hyphens are not allowed")
	}

	if !cfg.legacyFinalZeroBytePadding && strings.HasPrefix(in, "-") {
		return nil, fmt.Errorf("invalid proquint, leading hyphen is not allowed")
	}

	hasFinalHyphen := strings.HasSuffix(in, "-")
	if hasFinalHyphen {
		in = in[:len(in)-1]
	}

	if strings.Contains(in, "-") {
		for i := 0; i < len(in); i++ {
			if i%6 == 5 && in[i] != '-' {
				return nil, fmt.Errorf("invalid proquint, proquint with hyphens, but not between all syllables")
			}

			if i%6 != 5 && in[i] == '-' {
				return nil, fmt.Errorf("invalid proquint, hyphen in unexpected position detected")
			}
		}

		in = strings.ReplaceAll(in, "-", "")
	}

	in = strings.ToLower(in)

	if len(in) == 0 {
		return nil, fmt.Errorf("invalid proquint, length is 0")
	}

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

	isFinalHyphenPaddingInvalidFinalByte := !cfg.legacyFinalZeroBytePadding && hasFinalHyphen && finalByte != 0x00
	if isFinalHyphenPaddingInvalidFinalByte {
		return nil, fmt.Errorf("invalid proquint, final hyphen present, but last byte is not 0x00")
	}

	isFinalHyphenPadding := !cfg.legacyFinalZeroBytePadding && hasFinalHyphen && finalByte == 0x00
	isZeroBytePadding := cfg.legacyFinalZeroBytePadding && finalByte == 0x00
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
	if strings.Contains(in, "--") {
		return 0, fmt.Errorf("invalid propuint, consecutive hyphens are not allowed")
	}

	if strings.HasPrefix(in, "-") {
		return 0, fmt.Errorf("invalid proquint, leading hyphen is not allowed")
	}

	if strings.Contains(in, "-") {
		for i := 0; i < len(in); i++ {
			if i%6 == 5 && in[i] != '-' {
				return 0, fmt.Errorf("invalid proquint, proquint with hyphens, but not between all syllables")
			}

			if i%6 != 5 && in[i] == '-' {
				return 0, fmt.Errorf("invalid proquint, hyphen in unexpected position detected")
			}
		}

		in = strings.ReplaceAll(in, "-", "")
	}

	in = strings.ToLower(in)

	if len(in) != 10 {
		return 0, fmt.Errorf("invalid input, expect 10 characters (without hyphen), got %d", len(in))
	}

	var res uint32

	for i := range 2 {
		ui16, err := ToUint16(in[i*5 : (i+1)*5])
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
	if strings.Contains(in, "--") {
		return 0, fmt.Errorf("invalid propuint, consecutive hyphens are not allowed")
	}

	if strings.HasPrefix(in, "-") {
		return 0, fmt.Errorf("invalid proquint, leading hyphen is not allowed")
	}

	if strings.Contains(in, "-") {
		for i := 0; i < len(in); i++ {
			if i%6 == 5 && in[i] != '-' {
				return 0, fmt.Errorf("invalid proquint, proquint with hyphens, but not between all syllables")
			}

			if i%6 != 5 && in[i] == '-' {
				return 0, fmt.Errorf("invalid proquint, hyphen in unexpected position detected")
			}
		}

		in = strings.ReplaceAll(in, "-", "")
	}

	in = strings.ToLower(in)

	if len(in) != 20 {
		return 0, fmt.Errorf("invalid input, expect 20 characters (without hyphen), got %d", len(in))
	}

	var res uint64

	for i := range 4 {
		ui16, err := ToUint16(in[i*5 : (i+1)*5])
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
