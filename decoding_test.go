package proquint_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/breml/proquint"
)

func TestToBytes(t *testing.T) {
	tests := []struct {
		name            string
		in              string
		decodingOptions []proquint.DecodingOption

		assertErr require.ErrorAssertionFunc
		want      []byte
	}{
		{
			name: "regular - without dash",
			in:   "kivafdamur",

			assertErr: require.NoError,
			want:      []byte{0x67, 0x82, 0x12, 0x3b},
		},
		{
			name: "regular - with dash",
			in:   "kivafdamur",

			assertErr: require.NoError,
			want:      []byte{0x67, 0x82, 0x12, 0x3b},
		},
		{
			name: "uppercase - without dash",
			in:   "KIVAFDAMUR",

			assertErr: require.NoError,
			want:      []byte{0x67, 0x82, 0x12, 0x3b},
		},
		{
			name: "uppercase - with dash",
			in:   "KIVAF-DAMUR",

			assertErr: require.NoError,
			want:      []byte{0x67, 0x82, 0x12, 0x3b},
		},
		{
			name: "error - empty input",
			in:   "",

			assertErr: require.Error,
			want:      nil,
		},
		{
			name: "error - only final hyphen without syllables",
			in:   "-",

			assertErr: require.Error,
			want:      nil,
		},
		{
			name: "error - invalid length",
			in:   "bahaf-basa",

			assertErr: require.Error,
			want:      nil,
		},
		{
			name: "error - invalid character",
			in:   "bahaf-baXsa",

			assertErr: require.Error,
			want:      nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := proquint.ToBytes(tc.in, tc.decodingOptions...)
			tc.assertErr(t, err)

			require.Equal(t, tc.want, got)
		})
	}
}

func TestVectorsFromDraftRaynerToBytes(t *testing.T) {
	// Test vectors form Draft Rayner https://datatracker.ietf.org/doc/draft-rayner-proquint/04/

	tests := []struct {
		name            string
		in              string
		decodingOptions []proquint.DecodingOption

		assertErr require.ErrorAssertionFunc
		want      []byte
	}{
		{
			name: "babab",
			in:   "babab",

			assertErr: require.NoError,
			want:      []byte{0x00, 0x00},
		},
		{
			name: "zuzuz",
			in:   "zuzuz",

			assertErr: require.NoError,
			want:      []byte{0xFF, 0xFF},
		},
		{
			name: "damuh",
			in:   "damuh",

			assertErr: require.NoError,
			want:      []byte{0x12, 0x34},
		},
		{
			name: "zabat",
			in:   "zabat",

			assertErr: require.NoError,
			want:      []byte{0xF0, 0x0D},
		},
		{
			name: "ruroz",
			in:   "ruroz",

			assertErr: require.NoError,
			want:      []byte{0xBE, 0xEF},
		},
		{
			name: "damuh-zabat",
			in:   "damuh-zabat",

			assertErr: require.NoError,
			want:      []byte{0x12, 0x34, 0xF0, 0x0D},
		},
		{
			name: "himug-lamuh-gajaz-lijuh-hubuh-lisab",
			in:   "himug-lamuh-gajaz-lijuh-hubuh-lisab",

			assertErr: require.NoError,
			want:      []byte{0x46, 0x33, 0x72, 0x34, 0x31, 0x4F, 0x75, 0x74, 0x4C, 0x34, 0x77, 0x00},
		},
		{
			name: "himug-lamuh-gajaz-lijuh-hubuh-lisab with zero byte padding",
			in:   "himug-lamuh-gajaz-lijuh-hubuh-lisab",
			decodingOptions: []proquint.DecodingOption{
				proquint.LegacyWithFinalZeroBytePadding(),
			},

			assertErr: require.NoError,
			want:      []byte{0x46, 0x33, 0x72, 0x34, 0x31, 0x4F, 0x75, 0x74, 0x4C, 0x34, 0x77},
		},
		{
			name: "himug-lamuh-gajaz-lijuh-hubuh-lisab with final hyphen padding",
			in:   "himug-lamuh-gajaz-lijuh-hubuh-lisab",

			assertErr: require.NoError,
			want:      []byte{0x46, 0x33, 0x72, 0x34, 0x31, 0x4F, 0x75, 0x74, 0x4C, 0x34, 0x77, 0x00},
		},
		{
			name: "himug-lamuh-gajaz-lijuh-hubuh-lisab- with final hyphen padding",
			in:   "himug-lamuh-gajaz-lijuh-hubuh-lisab-",

			assertErr: require.NoError,
			want:      []byte{0x46, 0x33, 0x72, 0x34, 0x31, 0x4F, 0x75, 0x74, 0x4C, 0x34, 0x77},
		},

		// Padding examples.
		{
			name: "with final zero byte without padding hyphen",
			in:   "bahaf-basab",

			assertErr: require.NoError,
			want:      []byte{0x1, 0x2, 0x3, 0x0},
		},
		{
			name: "with final zero byte with padding hyphen",
			in:   "bahaf-basab-",

			assertErr: require.NoError,
			want:      []byte{0x1, 0x2, 0x3},
		},
		{
			name: "legacy with final zero padding",
			in:   "bahaf-basab",
			decodingOptions: []proquint.DecodingOption{
				proquint.LegacyWithFinalZeroBytePadding(),
			},

			assertErr: require.NoError,
			want:      []byte{0x1, 0x2, 0x3},
		},
		{
			name: "error - with padding hyphen but non zero final byte",
			in:   "bahaf-basad-",

			assertErr: require.Error,
			want:      nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := proquint.ToBytes(tc.in, tc.decodingOptions...)
			tc.assertErr(t, err)

			require.Equal(t, tc.want, got)
		})
	}
}

func TestToInt16(t *testing.T) {
	tests := []struct {
		name string
		in   string

		assertErr require.ErrorAssertionFunc
		want      int16
	}{
		{
			name: "success",
			in:   "kivaf",

			assertErr: require.NoError,
			want:      26498,
		},
		{
			name: "error - to big for int16",
			in:   "kivaf-kivaf",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "error - invalid length",
			in:   "kiva",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "error - invalid character",
			in:   "kiXaf",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "error - zero length",
			in:   "",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "error - only final hyphen for padding",
			in:   "-",

			assertErr: require.Error,
			want:      0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := proquint.ToInt16(tc.in)
			tc.assertErr(t, err)

			require.Equal(t, tc.want, got)
		})
	}
}

func TestToInt32(t *testing.T) {
	tests := []struct {
		name string
		in   string

		assertErr require.ErrorAssertionFunc
		want      int32
	}{
		{
			name: "success",
			in:   "kivaf-damur",

			assertErr: require.NoError,
			want:      1736577595,
		},
		{
			name: "error - invalid length",
			in:   "kivaf",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "error - invalid character",
			in:   "kivaf-daXur",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "error - zero length",
			in:   "",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "error - only final hyphen for padding",
			in:   "-",

			assertErr: require.Error,
			want:      0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := proquint.ToInt32(tc.in)
			tc.assertErr(t, err)

			require.Equal(t, tc.want, got)
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name string
		in   string

		assertErr require.ErrorAssertionFunc
		want      int64
	}{
		{
			name: "success",
			in:   "kivaf-damur-zabal-hilup",

			assertErr: require.NoError,
			want:      7458543981518341626,
		},
		{
			name: "error - invalid length",
			in:   "kivaf-damur-zabal",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "error - invalid character",
			in:   "kivaf-damur-zabal-hiXup",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "error - zero length",
			in:   "",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "error - only final hyphen for padding",
			in:   "-",

			assertErr: require.Error,
			want:      0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := proquint.ToInt64(tc.in)
			tc.assertErr(t, err)

			require.Equal(t, tc.want, got)
		})
	}
}
