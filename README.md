# pkcs7

[![GoDoc](https://godoc.org/go.mozilla.org/pkcs7?status.svg)](https://godoc.org/go.mozilla.org/pkcs7)
[![Build Status](https://github.com/mozilla-services/pkcs7/workflows/CI/badge.svg?branch=master&event=push)](https://github.com/mozilla-services/pkcs7/actions/workflows/ci.yml?query=branch%3Amaster+event%3Apush)

pkcs7 implements parsing and creating signed and enveloped messages.

```go
package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

    "go.mozilla.org/pkcs7"
)

func SignAndDetach(content []byte, cert *x509.Certificate, privkey *rsa.PrivateKey) (signed []byte, err error) {
	toBeSigned, err := NewSignedData(content)
	if err != nil {
		err = fmt.Errorf("Cannot initialize signed data: %s", err)
		return
	}
	if err = toBeSigned.AddSigner(cert, privkey, SignerInfoConfig{}); err != nil {
		err = fmt.Errorf("Cannot add signer: %s", err)
		return
	}

	// Detach signature, omit if you want an embedded signature
	toBeSigned.Detach()

	signed, err = toBeSigned.Finish()
	if err != nil {
		err = fmt.Errorf("Cannot finish signing data: %s", err)
		return
	}

	// Verify the signature
	pem.Encode(os.Stdout, &pem.Block{Type: "PKCS7", Bytes: signed})
	p7, err := pkcs7.Parse(signed)
	if err != nil {
		err = fmt.Errorf("Cannot parse our signed data: %s", err)
		return
	}

	// since the signature was detached, reattach the content here
	p7.Content = content

	if bytes.Compare(content, p7.Content) != 0 {
		err = fmt.Errorf("Our content was not in the parsed data:\n\tExpected: %s\n\tActual: %s", content, p7.Content)
		return
	}
	if err = p7.Verify(); err != nil {
		err = fmt.Errorf("Cannot verify our signed data: %s", err)
		return
	}

	return signed, nil
}
```



## Credits
This is a fork of [fullsailor/pkcs7](https://github.com/fullsailor/pkcs7)

## TODO
0. verify_test.go: refactor to use other digest algorithms than MD5 and SHA-1
1. Add support for ~~ECDSA and ~~ CMS with ECC

# Supported Algorithms
RSA, ECDSA
MD5, SHA1,
SHA2, SHA3
DES, AES, ChaCha20Poly1305(TODO: decrypt)

# Version
1.3.8
