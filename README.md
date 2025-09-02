# proquint

[![Go Reference](https://pkg.go.dev/badge/github.com/breml/proquint.svg)](https://pkg.go.dev/github.com/breml/proquint)
[![Test Status](https://github.com/breml/proquint/workflows/Go%20Matrix/badge.svg)](https://github.com/breml/proquint/actions?query=workflow%3AGo%20Matrix)
[![Go Report Card](https://goreportcard.com/badge/github.com/breml/proquint)](https://goreportcard.com/report/github.com/breml/proquint)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Implementation of the Proquint encoding scheme in Go.

Proquints are identifiers, that are readable, spellable, and pronounceable.

This package provides functions for encoding and decoding Proquints from `[]byte`
(also supporting `netip.Addr` or `github.com/google/uuid.UUID`), `u?int(16|32|64)`
and hex encoded strings.

This implementation follows the draft [draft-rayner-proquint-07](https://www.ietf.org/archive/id/draft-rayner-proquint-07.html),
which is a more strict version of the original proposal by Daniel S. Wilkerson.

The behavior from the original proposal can be enabled by using the `Legacy...` options.
See the [documentation](https://pkg.go.dev/github.com/breml/proquint) for details.

## Links

* [Proquint independent submission by Thomas Rayner - draft-rayner-proquint-07](https://www.ietf.org/archive/id/draft-rayner-proquint-07.html), [datatracker][https://datatracker.ietf.org/doc/draft-rayner-proquint/)
* [Proquint original proposal by Daniel S. Wilkerson](http://arXiv.org/html/0901.4016)
* [dsw/proquint - Daniel S. Wilkerson's reference implementations and links to other implementations](https://github.com/dsw/proquint)
* [proquint - alternative Go implementation](https://github.com/icco/proquint)

## Author

Copyright 2025 by Lucas Bremgartner ([breml](https://github.com/breml))

## License

[MIT License](LICENSE)
