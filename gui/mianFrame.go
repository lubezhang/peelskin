package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/lubezhang/peelskin/extract"
)

func Show() {
	a := app.New()
	w := a.NewWindow("extract")

	w.Resize(fyne.NewSize(800, 600))

	w.SetContent(layoutMain())
	w.ShowAndRun()
}

func layoutMain() *fyne.Container {
	inputUrl := widget.NewEntry()
	// inputUrl.SetText("https://www.iqiyi.com/v_19rrjac83w.html?vfrm=pcw_home&vfrmblk=712211_cainizaizhui&vfrmrst=712211_cainizaizhui_image1")
	inputUrl.Enable()
	inputUrl.SetPlaceHolder("input page url")

	input2 := widget.NewEntry()
	input2.Disable()
	input2.MultiLine = true
	input2.Wrapping = fyne.TextWrapWord

	btn2 := widget.NewButton("extract", func() {
		xmjx := extract.NewXmjx(inputUrl.Text)
		video, _ := xmjx.ExtractVideo()
		input2.SetText(video.Url)
	})

	// hcontainer1 := container.New(layout.NewHBoxLayout(), inputUrl)
	// hcontainer2 := container.New(layout.NewMaxLayout(), btn3, canvas.NewImageFromResource(theme.FyneLogo()))

	c1 := container.NewBorder(nil, nil, nil, btn2, inputUrl)

	return container.New(layout.NewVBoxLayout(), c1, input2)
}
