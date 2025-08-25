package proquint_test

import (
	"regexp"
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
			in:   "kivaf-damur",

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
			name: "some hyphens 1",
			in:   "kivaf-damurzabalhilup",

			assertErr: require.Error,
			want:      nil,
		},
		{
			name: "some hyphens 2",
			in:   "kivafdamur-zabalhilup",

			assertErr: require.Error,
			want:      nil,
		},
		{
			name: "some hyphens 3",
			in:   "kivaf-damur-zabalhilup",

			assertErr: require.Error,
			want:      nil,
		},
		{
			name: "misplaced hyphen",
			in:   "kiva-fdamur-zabal-hilup",

			assertErr: require.Error,
			want:      nil,
		},
		{
			name: "extreme hyphens",
			in:   "k-i-v-a-f-d-a-m-u-r-z-a-b-a-l-h-i-l-u-p",

			assertErr: require.Error,
			want:      nil,
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
			name: "error - double hyphen",
			in:   "bahaf--basab",

			assertErr: require.Error,
			want:      nil,
		},
		{
			name: "error - leading hyphen",
			in:   "-bahaf-basab",

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
			name: "error - double hyphen",
			in:   "kivaf--damur",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "error - leading hyphen",
			in:   "-kivaf-damur",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "misplaced hyphen",
			in:   "kivafd-amur",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "misplaced hyphen 2",
			in:   "kiva-fdamur",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "extreme hyphens",
			in:   "k-i-v-a-f-d-a-m-u-r",

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
			name: "error - double hyphen",
			in:   "kivaf-damur-zabal--hilup",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "error - leading hyphen",
			in:   "-kivaf-damur-zabal-hilup",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "some hyphens 1",
			in:   "kivaf-damurzabalhilup",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "some hyphens 2",
			in:   "kivafdamur-zabalhilup",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "some hyphens 3",
			in:   "kivaf-damur-zabalhilup",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "missplaced hyphen",
			in:   "kiva-fdamur-zabal-hilup",

			assertErr: require.Error,
			want:      0,
		},
		{
			name: "extreme hyphens",
			in:   "k-i-v-a-f-d-a-m-u-r-z-a-b-a-l-h-i-l-u-p",

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

var proquintRegex = regexp.MustCompile(`(?i)^(([bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz])*|([bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz]-)+)(([bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz])|([bdfghjklmnprstvz][aiou][bhms]ab-))$`)

var corpus = []string{
	`babab`,
	`zuzuz`,
	`damuh`,
	`zabat`,
	`ruroz`,
	`damuh-zabat`,
	`himug-lamuh-gajaz-lijuh-hubuh-lisab`,
	`himug-lamuh-gajaz-lijuh-hubuh-lisab`,
	`himug-lamuh-gajaz-lijuh-hubuh-lisab`,
	`himug-lamuh-gajaz-lijuh-hubuh-lisab-`,
	`kivafdamur`,
	`kivafdamur`,
	`kivaf-damur`,
	`KIVAFDAMUR`,
	`KIVAF-DAMUR`,
	`bahaf-basab`,
	`bahaf-basab`,
	`bahaf-basab-`,
}

func FuzzToBytes(f *testing.F) {
	for _, quint := range corpus {
		f.Add(quint)
	}

	f.Fuzz(func(t *testing.T, in string) {
		regRes := proquintRegex.MatchString(in)
		got, err := proquint.ToBytes(in)

		if (err == nil) != regRes {
			t.Errorf("proquint %q produced %v with FromBytes and %t with regex match: %t != %t, got: %q", string(in), err, regRes, err == nil, regRes, got)
		}
	})
}
