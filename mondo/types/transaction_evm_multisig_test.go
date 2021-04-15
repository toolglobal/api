package types

import (
	"encoding/hex"
	"fmt"
	"math/big"
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

/*
{
	"deadline": 1608370617,
	"gasLimit": 21000,
	"gasPrice": "100",
	"from": "64FB3D7E6e338B80fBd127052Dc75A2A7901Aa44",
	"nonce": 0,
	"to": "8c5015F85B993243A80478D53bc951f84F553dBB",
	"value": "1200000000",
	"load": "",
	"memo": "",
	"signature": "{\"signatures\":[\"0xe53a4f71c4b308966d7110b26dfd785939b59b0802a05189fee5a23f5e3627a105d92009aa9a0975fa634460728a26331ae3be6f556ab6ea41d4cba4cb125a5e01\"],
\"pubKey\":{\"pubkeys\":[\"02331663b3ef21a37c7270185284617528f7dc9d6ca69c005547942a780ee9cc1b\",\"02abff05e8b4e7d6d2299a9d191ec97e5aa0e099e7e31c005a07940e5978b555d9\",\"03815a906de2017c7351be33644cd60a6fff9407ce04896b2328944bc4e628abd8\"],\"threshold\":2}}"
}
*/
func Test_CheckSign(t *testing.T) {
	var pkeys []PublicKey
	pkey, _ := HexToPubkey("02331663b3ef21a37c7270185284617528f7dc9d6ca69c005547942a780ee9cc1b")
	pkeys = append(pkeys, pkey)

	pkey, _ = HexToPubkey("02abff05e8b4e7d6d2299a9d191ec97e5aa0e099e7e31c005a07940e5978b555d9")
	pkeys = append(pkeys, pkey)

	pkey, _ = HexToPubkey("03815a906de2017c7351be33644cd60a6fff9407ce04896b2328944bc4e628abd8")
	pkeys = append(pkeys, pkey)

	pubkey := NewPubKeyMultisigThreshold(2, pkeys)

	tx := NewMultisigEvmTx(2, pkeys)
	tx.Deadline = 1608370617
	tx.GasPrice = new(big.Int).SetInt64(100)
	tx.GasLimit = 21000
	tx.From = pubkey.Address()
	tx.To = ethcmn.HexToAddress("8c5015F85B993243A80478D53bc951f84F553dBB")
	tx.Value = new(big.Int).SetInt64(1200000000)
	tx.Signature.PubKey = pubkey
	tx.AddSign("0xe53a4f71c4b308966d7110b26dfd785939b59b0802a05189fee5a23f5e3627a105d92009aa9a0975fa634460728a26331ae3be6f556ab6ea41d4cba4cb125a5e01")

	{
		privBytes, _ := hex.DecodeString("e53a4f71c4b308966d7110b26dfd785939b59b0802a05189fee5a23f5e3627a105d92009aa9a0975fa634460728a26331ae3be6f556ab6ea41d4cba4cb125a5e01")
		pubkey, err := crypto.SigToPub(tx.SigHash().Bytes(), privBytes)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(hex.EncodeToString(crypto.CompressPubkey(pubkey)))
	}

	//if !tx.Verify() {
	//	t.Fatal("check failed")
	//}
}
