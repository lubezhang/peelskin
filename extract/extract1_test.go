package extract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractM3u8(t *testing.T) {
	assetObj := assert.New(t)
	assetObj.Equal("e10adc3949ba59abbe56e057f20f883e", "e10adc3949ba59abbe56e057f20f883e")
	// ExtractM3u8("https://v.qq.com/x/cover/mzc00200imi2b3v/l00344c9o6b.html")
	ExtractM3u8("https://www.iqiyi.com/v_230c1yg1uzo.html?vfrm=pcw_dianying&vfrmblk=E&vfrmrst=711219_dianying_top_video_play4")
}

func TestGetMD5(t *testing.T) {
	assetObj := assert.New(t)
	assetObj.Equal(GetMD5("rXjWvXl6"), "dac2c08904d3011ff90564c9f77865cf")
}
