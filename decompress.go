package cms

import (
	"bytes"
	"compress/zlib"
	"encoding/asn1"
	"errors"
	"io"
)

var ErrNotCompressedContent = errors.New("pkcs7: content data is not a compressed data type")

func (p7 *PKCS7) Decompress() ([]byte, error) {
	// 1. convert to CompressedData
	compressedData, ok := p7.raw.(compressedData)
	if !ok {
		return nil, ErrNotCompressedContent
	}

	// 2. parse EncapsulatedContentInfo
	var encapsulatedContentInfo encapsulatedContentInfo
	_, err := asn1.Unmarshal(compressedData.EncapContentInfo.EContent, &encapsulatedContentInfo)
	if err != nil {
		return nil, err
	}

	// 3. decompress using zlib
	buf := bytes.NewBuffer(encapsulatedContentInfo.EContent)
	r, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}
