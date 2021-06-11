package types

import (
	"fmt"
	"testing"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func Test_signAndVerify(t *testing.T) {
	k, n := 3, 5
	acts := generateAccounts(n)
	var pkeySet []PublicKey
	for _, v := range acts {
		var pkey PublicKey
		copy(pkey[:], ethcmn.Hex2Bytes(v.PubKey))
		pkeySet = append(pkeySet, pkey)
	}

	// 多签账户
	pkey := NewPubKeyMultisigThreshold(k, pkeySet)

	tx := NewMultisigEvmTx(k, pkeySet)
	tx.From = pkey.Address()
	tx.Signature.PubKey = pkey

	if err := tx.Sign(acts[0].PrivKey); err != nil {
		t.Fatal(err)
	}

	if err := tx.Sign(acts[1].PrivKey); err != nil {
		t.Fatal(err)
	}

	if err := tx.Sign(acts[3].PrivKey); err != nil {
		t.Fatal(err)
	}

	b := tx.Verify()
	if !b {
		t.Fatal(b)
	}
}

func Test_Signers(t *testing.T) {
	k, n := 3, 5
	acts := generateAccounts(n)
	var pkeySet []PublicKey
	for _, v := range acts {
		var pkey PublicKey
		copy(pkey[:], ethcmn.Hex2Bytes(v.PubKey))
		pkeySet = append(pkeySet, pkey)
	}

	// 多签账户
	pkey := NewPubKeyMultisigThreshold(k, pkeySet)

	tx := NewMultisigEvmTx(k, pkeySet)
	tx.From = pkey.Address()
	tx.Signature.PubKey = pkey

	if err := tx.Sign(acts[0].PrivKey); err != nil {
		t.Fatal(err)
	}

	if err := tx.Sign(acts[3].PrivKey); err != nil {
		t.Fatal(err)
	}

	for _, v := range tx.Signer() {
		fmt.Println(v.Hex())
	}
}

type account struct {
	PrivKey string
	PubKey  string
	Address string
}

func generateAccounts(n int) []account {
	acts := make([]account, n)
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		buff := make([]byte, 32)
		copy(buff[32-len(key.D.Bytes()):], key.D.Bytes())
		acts[i].PrivKey = ethcmn.Bytes2Hex(buff)
		acts[i].PubKey = ethcmn.Bytes2Hex(crypto.CompressPubkey(&key.PublicKey))
		acts[i].Address = crypto.PubkeyToAddress(key.PublicKey).String()
		fmt.Println("generateAccounts ", acts[i].Address)
	}

	return acts
}
