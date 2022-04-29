package extract

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// Decrypt Golang解密
// ciphertext  important important 上面js的生成的密文进行了 hex.encoding 在这之前必须要进行 hex.Decoding
// 上面js代码最后返回的是16进制
// 所以收到的数据hexText还需要用hex.DecodeString(hexText)转一下,这里略了
func Decrypt(ciphertext, key []byte, iv []byte) ([]byte, error) {
	// pkey := PaddingLeft(key, ' ', 16) //和js的key补码方法一致
	pkey := key

	block, err := aes.NewCipher(pkey) //选择加密算法
	if err != nil {
		return nil, fmt.Errorf("key 长度必须 16/24/32长度: %s", err)
	}
	blockModel := cipher.NewCBCDecrypter(block, iv[:block.BlockSize()]) //和前端代码对应:   mode: CryptoJS.mode.CBC,// CBC算法
	plantText := make([]byte, len(ciphertext))
	blockModel.CryptBlocks(plantText, ciphertext)
	// plantText = PKCS7UnPadding(plantText) //和前端代码对应:  padding: CryptoJS.pad.Pkcs7
	plantText = NullUnPadding(plantText) //和前端代码对应:  padding: CryptoJS.pad.ZeroPadding
	return plantText, nil
}

func Encrypt(ciphertext, key []byte) (string, error) {
	// key = PaddingLeft(key, '0', 16) //和js的key补码方法一致
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	// iv := key[:aes.BlockSize]
	// ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	// if len(ciphertext)%aes.BlockSize != 0 {
	// 	panic("ciphertext is not a multiple of the block size")
	// }
	blockSize := block.BlockSize()
	ciphertext = PKCS7Padding(ciphertext, blockSize)
	mode := cipher.NewCBCEncrypter(block, key[:blockSize])

	crypted := make([]byte, len(ciphertext))
	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(crypted, ciphertext)
	fmt.Println("===:", base64.StdEncoding.EncodeToString(crypted))
	fmt.Println("===:", hex.EncodeToString(crypted))
	return "", nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//这个方案必须和js的方法是一样的
func PaddingLeft(ori []byte, pad byte, length int) []byte {
	if len(ori) >= length {
		return ori[:length]
	}
	pads := bytes.Repeat([]byte{pad}, length-len(ori))
	return append(pads, ori...)
}

func GetMD5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

func NullUnPadding(in []byte) []byte {
	return bytes.TrimRight(in, string([]byte{0}))
}
