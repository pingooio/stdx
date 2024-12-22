package chacha20blake3_test

import (
	"encoding/hex"
)

func toHex(bits []byte) string {
	return hex.EncodeToString(bits)
}

func fromHex(bits string) []byte {
	b, err := hex.DecodeString(bits)
	if err != nil {
		panic(err)
	}
	return b
}

// func TestVectorsChaCha20Blake3(t *testing.T) {
// 	for i, v := range chacha20Blake3Vectors {
// 		dst := make([]byte, len(v.plaintext)+chacha20blake3.TagSize)

// 		cipher, err := chacha20blake3.New(v.key)
// 		if err != nil {
// 			t.Errorf("plaintext: %s", toHex(v.plaintext))
// 			t.Errorf("nonce: %s", toHex(v.nonce))
// 			t.Fatal(err)
// 		}

// 		cipher.Seal(dst[:0], v.nonce, v.plaintext, v.additionalData)
// 		if !bytes.Equal(dst, v.ciphertext) {
// 			t.Errorf("Test %d: ciphertext mismatch:\ngot:  %s\nwant: %s", i, toHex(dst), toHex(v.ciphertext))
// 		}

// 		decryptedPlaintext, err := cipher.Open(nil, v.nonce, dst, v.additionalData)
// 		if err != nil {
// 			t.Errorf("Test %d: %v", i, err)
// 		}
// 		if !bytes.Equal(decryptedPlaintext, v.plaintext) {
// 			t.Errorf("Test %d: plaintext mismatch:\ngot:  %s\nwant: %s", i, toHex(decryptedPlaintext), toHex(v.plaintext))
// 		}
// 	}
// }

// func TestBasicX(t *testing.T) {
// 	var key [chacha20blake3.KeySize]byte
// 	var nonce [chacha20blake3.NonceSizeX]byte

// 	originalPlaintext := []byte("Hello World")
// 	additionalData := []byte("!")

// 	rand.Read(key[:])
// 	rand.Read(nonce[:])

// 	cipher, _ := chacha20blake3.NewX(key[:])
// 	ciphertext := cipher.Seal(nil, nonce[:], originalPlaintext, additionalData)

// 	decryptedPlaintext, err := cipher.Open(nil, nonce[:], ciphertext, additionalData)
// 	if err != nil {
// 		t.Errorf("decrypting message: %s", err)
// 		return
// 	}

// 	if !bytes.Equal(decryptedPlaintext, originalPlaintext) {
// 		t.Errorf("original message (%s) != decrypted message (%s)", string(originalPlaintext), string(decryptedPlaintext))
// 		return
// 	}
// }

// func TestAdditionalDataX(t *testing.T) {
// 	var key [chacha20blake3.KeySize]byte
// 	var nonce [chacha20blake3.NonceSizeX]byte

// 	originalPlaintext := []byte("Hello World")
// 	additionalData := []byte("!")

// 	rand.Read(key[:])
// 	rand.Read(nonce[:])

// 	cipher, _ := chacha20blake3.NewX(key[:])
// 	ciphertext := cipher.Seal(nil, nonce[:], originalPlaintext, additionalData)

// 	_, err := cipher.Open(nil, nonce[:], ciphertext, []byte{})
// 	if !errors.Is(err, chacha20blake3.ErrOpen) {
// 		t.Errorf("expected error (%s) | got (%s)", chacha20blake3.ErrOpen, err)
// 		return
// 	}
// }

// var chacha20Blake3Vectors = []struct {
// 	key            []byte
// 	nonce          []byte
// 	plaintext      []byte
// 	additionalData []byte
// 	ciphertext     []byte
// }{
// 	{
// 		fromHex("0000000000000000000000000000000000000000000000000000000000000000"),
// 		fromHex("0000000000000000"),
// 		fromHex("48656C6C6F20576F726C6421"), //  Hello World!
// 		nil,
// 		fromHex("d7628bd23a716f15ead6f35d7a6be6fd6ee37de0c569f11705e9c1ba5b576d0886d0f7f1b5fdb82ea07856d7"),
// 	},
// 	{
// 		fromHex("0100000000000000000000000000000000000000000000000000000000000010"),
// 		fromHex("0100000000000010"),
// 		fromHex("4368614368613230"), // ChaCha20
// 		fromHex("424C414B4533"),     // BLAKE3
// 		fromHex("9f7a1c78e49065c7627ebc71dfe55740dc0ff6164667b3192b928a189da4b6153b579262acda29bb"),
// 	},
// 	{
// 		fromHex("a9541ec64e971c19216360a28aebffdefdbc2f2b4f8d683a2c5c17c12e86059d"),
// 		fromHex("5722cf5d7efbc3a1"),
// 		fromHex("112103a99299c403eb92c29ee81f8faa2c4bab00ef4a92ddb3cf7d0c3ec63d19b81ff83defbfa34fb1ac5bf594306a541fb4ba3c18f700d6d38d2eed4f118760"),
// 		fromHex("bd3c8a9c2c9362c392dd9b9ae7e31552"),
// 		fromHex("f8742d7ec6862e53715e526f1b91c8c3d60b005d93ea924ca7377e81d8cad3f69d102d604c0688befb8c0fbdfc499bf10ab55021e87d66b2cdee57401a93c4ddcd2501adf933f45bb6309810c744205d0e9d69b9943afa2f01188e2b91be7f11"),
// 	},
// }
