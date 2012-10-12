package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"hash"
	"time"
)

type HashEasy struct {
	hash.Hash
}

func NewHashEasy() HashEasy {
	e := HashEasy{sha1.New()}
	return e
}

func (h HashEasy) WriteStrAnd(s string) HashEasy {
	h.Write([]byte(s))
	return h
}

func (h HashEasy) Nonce() HashEasy {
	nanos := time.Now().UnixNano()
	binary.Write(h, binary.LittleEndian, nanos)
	return h
}

func (h HashEasy) WriteIntAnd(i int) HashEasy {
	binary.Write(h, binary.LittleEndian, int64(i))
	return h
}

func (h HashEasy) String() string {
	buffer := new(bytes.Buffer)
	base64.NewEncoder(base64.URLEncoding, buffer).Write(h.Sum(nil))
	return buffer.String()
}
