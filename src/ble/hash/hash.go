package hash

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	hh "hash"
	"log"
	"time"
)

type HashEasy struct {
	hh.Hash
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

func EasyNonce(stuff ...interface{}) string {
	h := NewHashEasy()
	for _, v := range stuff {
		switch v.(type) {
		case string:
			h.WriteStrAnd(v.(string))
		case int:
			h.WriteIntAnd(v.(int))
		case fmt.Stringer:
			h.WriteStrAnd(v.(fmt.Stringer).String())
		default:
			log.Print("Unexpected value %#v", v)
		}
	}
	return h.String()
}
