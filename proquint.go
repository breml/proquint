// Package proquint provides an implementation of the Proquint encoding scheme
// as described in https://www.ietf.org/archive/id/draft-rayner-proquint-07.html.
//
// The original behavior described in http://arXiv.org/html/0901.4016 can be
// enabled using the `Legacy...` options.
package proquint

var consonants = []byte{
	'b', 'd', 'f', 'g',
	'h', 'j', 'k', 'l',
	'm', 'n', 'p', 'r',
	's', 't', 'v', 'z',
}

var vowel = []byte{
	'a', 'i', 'o', 'u',
}
