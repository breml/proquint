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
			name: "with zero padding",
			in:   "bahaf-basab",
			decodingOptions: []proquint.DecodingOption{
				proquint.WithFinalZeroBytePadding(),
			},

			assertErr: require.NoError,
			want:      []byte{0x1, 0x2, 0x3},
		},
		{
			name: "with final hyphen padding without final hyphen",
			in:   "bahaf-basab",
			decodingOptions: []proquint.DecodingOption{
				proquint.WithFinalHyphenPadding(),
			},

			assertErr: require.NoError,
			want:      []byte{0x1, 0x2, 0x3, 0x0},
		},
		{
			name: "with final hyphen padding with final hyphen",
			in:   "bahaf-basab-",
			decodingOptions: []proquint.DecodingOption{
				proquint.WithFinalHyphenPadding(),
			},

			assertErr: require.NoError,
			want:      []byte{0x1, 0x2, 0x3},
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := proquint.ToInt64(tc.in)
			tc.assertErr(t, err)

			require.Equal(t, tc.want, got)
		})
	}
}
