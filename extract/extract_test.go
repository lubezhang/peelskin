package extract

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractM3u8(t *testing.T) {
	assetObj := assert.New(t)
	assetObj.Equal("e10adc3949ba59abbe56e057f20f883e", "e10adc3949ba59abbe56e057f20f883e")
	ExtractM3u8("https://v.qq.com/x/cover/mzc00200imi2b3v/l00344c9o6b.html")
}

func TestExtractM3u81(t *testing.T) {
	assetObj := assert.New(t)
	assetObj.Equal("e10adc3949ba59abbe56e057f20f883e", "e10adc3949ba59abbe56e057f20f883e")
	// ExtractM3u81()
}

func TestGetMD5(t *testing.T) {
	assetObj := assert.New(t)
	assetObj.Equal(GetMD5("rXjWvXl6"), "dac2c08904d3011ff90564c9f77865cf")
}

func TestEncodeJSData(t *testing.T) {
	origData := "https://cache.m3u8.shenglinyiyang.cn/duoduo/20220426/ce5304064cce73f628fdce2f021d1695.m3u8?st=sZSyCujsbh4Q4ZjNoxW1Fw&e=1650956401"

	orgCryptoData := "DLJuGFol+m6NMEznRGwUnoZkAycE1n4Wni/ZlmaCDntBuTe8ZkBdpEciE79wS8CMydPkfBudhSk6o6cacRGEzVh2DUcpJ8Rl7TUyQnWHRyzG01iO/pxlo9jdzOhoTLGc+1Vh9QDrQmGz20CVXjA+1jkckf+KEdHrU8XivUN7soRna6cSAC8xVb+SAPKW8tD9"
	decodeData, _ := base64.StdEncoding.DecodeString(orgCryptoData)

	// hexKey, _ := hex.DecodeString("rXjWvXl6")
	// decodeData1, _ := base64.StdEncoding.DecodeString("rXjWvXl6")
	// decodeData1 := []byte("rXjWvXl6oPYHs¤Ç")

	iv := []byte("NXbHoWJbpsEOin8b")
	key := GetMD5("rXjWvXl6")
	got, err := Decrypt(decodeData, []byte(key), iv)
	if err != nil {
		t.Errorf("DecryptJs() error = %v", err)
		return
	}
	assetObj := assert.New(t)
	assetObj.Equal(string(got), origData)
}

func TestEncodeString(t *testing.T) {
	key := []byte("dac2c08904d3011ff90564c9f77865cf")
	ciphertext := []byte("https://cache.m3u8.shenglinyiyang.cn/duoduo/20220426/ce5304064cce73f628fdce2f021d1695.m3u8?st=sZSyCujsbh4Q4ZjNoxW1Fw&e=1650956401")

	data1, _ := Encrypt(ciphertext, key)
	fmt.Println(data1)
	// assetObj := assert.New(t)
	// assetObj.Equal(GetMD5("rXjWvXl6"), "dac2c08904d3011ff90564c9f77865cf")
}
