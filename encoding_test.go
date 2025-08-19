package proquint_test

import (
	"fmt"
	"net/netip"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/breml/proquint"
)

func TestExamplesFromSpec(t *testing.T) {
	// Test Cases form specification: https://arxiv.org/html/0901.4016
	tests := []struct {
		name string
		in   []byte

		want string
	}{
		{
			name: "127.0.0.1",
			in:   []byte{127, 0, 0, 1},

			want: "lusab-babad",
		},
		{
			name: "63.84.220.193",
			in:   []byte{63, 84, 220, 193},

			want: "gutih-tugad",
		},
		{
			name: "63.118.7.35",
			in:   []byte{63, 118, 7, 35},

			want: "gutuk-bisog",
		},
		{
			name: "140.98.193.141",
			in:   []byte{140, 98, 193, 141},

			want: "mudof-sakat",
		},
		{
			name: "64.255.6.200",
			in:   []byte{64, 255, 6, 200},

			want: "haguz-biram",
		},
		{
			name: "128.30.52.45",
			in:   []byte{128, 30, 52, 45},

			want: "mabiv-gibot",
		},
		{
			name: "147.67.119.2",
			in:   []byte{147, 67, 119, 2},

			want: "natag-lisaf",
		},
		{
			name: "212.58.253.68",
			in:   []byte{212, 58, 253, 68},

			want: "tibup-zujah",
		},
		{
			name: "216.35.68.215",
			in:   []byte{216, 35, 68, 215},

			want: "tobog-higil",
		},
		{
			name: "216.68.232.21",
			in:   []byte{216, 68, 232, 21},

			want: "todah-vobij",
		},
		{
			name: "198.81.129.136",
			in:   []byte{198, 81, 129, 136},

			want: "sinid-makam",
		},
		{
			name: "12.110.110.204",
			in:   []byte{12, 110, 110, 204},

			want: "budov-kuras",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			quint, err := proquint.FromBytes(tc.in, proquint.WithHyphens())
			require.NoError(t, err)

			require.Equal(t, tc.want, quint)

			quint = proquint.FromInt32(int32(tc.in[0])<<24+int32(tc.in[1])<<16+int32(tc.in[2])<<8+int32(tc.in[3]), proquint.WithHyphens())
			require.NoError(t, err)

			require.Equal(t, tc.want, quint)

			bytes, err := proquint.ToBytes(quint)
			require.NoError(t, err)

			require.Equal(t, tc.in, bytes)
		})
	}
}

func TestFromBytes(t *testing.T) {
	// Test Cases form specification: https://arxiv.org/html/0901.4016
	tests := []struct {
		name            string
		in              []byte
		encodingOptions []proquint.EncodingOption

		assertErr require.ErrorAssertionFunc
		want      string
	}{
		{
			name: "127.0.0.1 - standard",
			in:   []byte{127, 0, 0, 1},

			assertErr: require.NoError,
			want:      "lusabbabad",
		},
		{
			name: "127.0.0.1 - with hyphens",
			in:   []byte{127, 0, 0, 1},
			encodingOptions: []proquint.EncodingOption{
				proquint.WithHyphens(),
			},

			assertErr: require.NoError,
			want:      "lusab-babad",
		},
		{
			name: "zero padding",
			in:   []byte{1, 2, 3},
			encodingOptions: []proquint.EncodingOption{
				proquint.WithPadding(),
				proquint.WithHyphens(),
			},

			assertErr: require.NoError,
			want:      "bahaf-basab",
		},
		{
			name: "hyphen padding regular final zero byte",
			in:   []byte{1, 2, 3, 0},
			encodingOptions: []proquint.EncodingOption{
				proquint.WithPaddingFinalHyphen(),
			},

			assertErr: require.NoError,
			want:      "bahaf-basab",
		},
		{
			name: "hyphen padding odd number of bytes",
			in:   []byte{1, 2, 3},
			encodingOptions: []proquint.EncodingOption{
				proquint.WithPaddingFinalHyphen(),
			},

			assertErr: require.NoError,
			want:      "bahaf-basab-",
		},
		{
			name: "error - odd number of bytes without padding",
			in:   []byte{1, 2, 3},

			assertErr: require.Error,
			want:      "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			quint, err := proquint.FromBytes(tc.in, tc.encodingOptions...)
			tc.assertErr(t, err)

			require.Equal(t, tc.want, quint)
		})
	}
}

func TestIntx(t *testing.T) {
	var in uint64 = 0xCE

	require.Equal(t, "bagav", proquint.FromUint16(uint16(in)))
	require.Equal(t, "bagav", proquint.FromInt16(int16(in)))
	require.Equal(t, "bababbagav", proquint.FromUint32(uint32(in)))
	require.Equal(t, "bababbagav", proquint.FromInt32(int32(in)))
	require.Equal(t, "babab-bagav", proquint.FromUint32(uint32(in), proquint.WithHyphens()))
	require.Equal(t, "babab-bagav", proquint.FromInt32(int32(in), proquint.WithHyphens()))
	require.Equal(t, "bababbababbababbagav", proquint.FromUint64(uint64(in)))
	require.Equal(t, "bababbababbababbagav", proquint.FromInt64(int64(in)))
	require.Equal(t, "babab-babab-babab-bagav", proquint.FromUint64(uint64(in), proquint.WithHyphens()))
	require.Equal(t, "babab-babab-babab-bagav", proquint.FromInt64(int64(in), proquint.WithHyphens()))
}

func TestHexToProquint(t *testing.T) {
	tests := []struct {
		name string
		in   string

		assertErr require.ErrorAssertionFunc
		want      string
	}{
		{
			name: "success",
			in:   "6782123b",

			assertErr: require.NoError,
			want:      "kivaf-damur",
		},
		{
			name: "error - invalid hex",
			in:   "invalid", // invalid hex string

			assertErr: require.Error,
			want:      "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			quint, err := proquint.FromHexString(tc.in, proquint.WithHyphens())
			tc.assertErr(t, err)

			require.Equal(t, tc.want, quint)
		})
	}
}

func ExampleFromBytes() {
	ipv4 := netip.MustParseAddr("127.0.0.1")
	ipv4bytes := ipv4.As4()
	quint, _ := proquint.FromBytes(ipv4bytes[:], proquint.WithHyphens())

	fmt.Println(quint)
	// Output: lusab-babad
}

func ExampleFromBytes_uuid() {
	id := uuid.MustParse(`6782123b-f007-45fa-a9b8-24d0136facd4`)
	quint, _ := proquint.FromBytes(id[:], proquint.WithHyphens())

	fmt.Println(quint)
	// Output: kivaf-damur-zabal-hilup-pokum-figib-datoz-pugih
}

func ExampleFromHexString() {
	quint, _ := proquint.FromHexString("6782123b", proquint.WithHyphens())

	fmt.Println(quint)
	// Output: kivaf-damur
}
