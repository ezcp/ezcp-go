package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"hash"
	"io"
)

func crypt(pass string, reader io.Reader) (io.Reader, error) {
	var r io.Reader
	r = reader
	if pass != "" {
		hashKey, iv := genIvAndKey([]byte{}, []byte(pass), md5.New(), 32, 1)
		aesBlock, err := aes.NewCipher(hashKey)
		if err != nil {
			panic(err)
		}
		stream := cipher.NewCTR(aesBlock, iv)

		r = &cipher.StreamReader{S: stream, R: reader}
	}
	return r, nil
}

// genIvAndKey was found here : https://play.golang.org/p/zsPMd8QN0b
// the thing is node.js creates its aes key in its own way
func genIvAndKey(salt, data []byte, h hash.Hash, keyLen, blockLen int) (key []byte, iv []byte) {
	res := make([]byte, 0, keyLen+blockLen)
	p := append(data, salt...)
	var dlast []byte

	for ; len(res) < keyLen+blockLen; h.Reset() {
		h.Write(append(dlast, p...))
		resNew := h.Sum(res)
		dlast = resNew[len(res):]
		res = resNew
	}

	return res[:keyLen], res[keyLen:]
}
