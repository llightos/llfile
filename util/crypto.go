package util

//AES对称加解密，注意加密后的[]byte不适用于传输，因此调用后需要进行base64或者其他方式的编码
import (
	"crypto/aes"
	"fmt"
)

func Encrypt(in []byte, key []byte) (re []byte) {
	// 得到一个cipher.Block接口
	cipher, _ := aes.NewCipher(generateKey(key))
	//length表示要分成几块 aes.BlockSize=16byte,128bit
	length := (len(in) + aes.BlockSize) / aes.BlockSize
	fmt.Println(length)
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, in)
	//pad表示分块后plain有多少空byte
	pad := byte(len(plain) - len(in))
	//fmt.Println("pad=", pad)
	//对plain进行填充
	//fmt.Println("填充前为", plain)
	for i := len(in); i < len(plain); i++ {
		plain[i] = pad
	}
	//fmt.Println("填充后为", plain)
	re = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(in); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(re[bs:be], plain[bs:be])
	}

	return re
}

func Decrypt(re []byte, key []byte) (decrypted []byte) {
	//得到一个cipher.Block接口
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(re))
	//分组分块解密
	for bs, be := 0, cipher.BlockSize(); bs < len(re); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], re[bs:be])
	}
	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}
	return decrypted[:trim]
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	//fmt.Println("genKey", string(genKey))
	//fmt.Println("genKey", len(genKey))
	return genKey
}
