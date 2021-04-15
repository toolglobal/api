package types

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"testing"
)

func TestTxParams_rlp(t *testing.T) {
	tx := NewTxParams()
	tx.CreatedAt = 1233
	pub, _ := base64.StdEncoding.DecodeString("EuhO0uxddP2NmF+AHouuYvV9CTQZIT1xm4GlKm13vFI=")
	tx.Sender = pub
	tx.Nonce = 123
	tx.Key = []byte("hello")
	tx.Value = []byte("world")
	priv, _ := base64.StdEncoding.DecodeString("k8C5HSZHX1+Ck1l1l9eSxsmsCiPgQciuoXXczy8M0BAS6E7S7F10/Y2YX4Aei65i9X0JNBkhPXGbgaUqbXe8Ug==")
	tx.Signature = tx.Sign(priv)

	fmt.Println(len(pub), len(priv))

	if !tx.Verify() {
		t.Fatal("verify failed")
	}

	bs := tx.ToBytes()
	tx2 := NewTxParams()
	err := tx2.FromBytes(bs)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(bs, tx2.ToBytes()) {
		t.Fatal("not equal")
	}
	fmt.Println(string(tx2.Key))
	fmt.Println(string(tx2.Value))
}
