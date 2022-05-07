package extract

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractUrl(t *testing.T) {
	assetObj := assert.New(t)
	assetObj.Equal("e10adc3949ba59abbe56e057f20f883e", "e10adc3949ba59abbe56e057f20f883e")

	// xmjx := NewXmjx("https://v.qq.com/x/cover/mzc00200imi2b3v/l00344c9o6b.html")
	// xmjx := NewXmjx("https://www.iqiyi.com/v_230c1yg1uzo.html?vfrm=pcw_dianying&vfrmblk=E&vfrmrst=711219_dianying_top_video_play4")
	xmjx := NewXmjx("https://v.youku.com/v_show/id_XNTg1NjE3MDM0OA==.html?spm=a2ha1.14919748_WEBHOME_GRAY.drawer5.d_zj1_2&s=bddeff276f7f4fde94ae&scm=20140719.rcmd.7182.show_bddeff276f7f4fde94ae")
	xmjx.ExtractUrl()
}

func TestEncryptFVkey(t *testing.T) {
	assetObj := assert.New(t)

	// xmjx := NewXmjx("https://v.qq.com/x/cover/mzc00200imi2b3v/l00344c9o6b.html")
	xmjx := NewXmjx("https://www.iqiyi.com/v_230c1yg1uzo.html?vfrm=pcw_dianying&vfrmblk=E&vfrmrst=711219_dianying_top_video_play4")
	vfkey := xmjx.encryptFVkey("fd48169c34f4b6808fc038de21dde261")
	assetObj.Equal("qUUb1eItdEqtXYqHfxUQd7eJNLFct79w8hjiyq04CwY=", vfkey)
}
