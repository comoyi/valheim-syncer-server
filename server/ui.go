package server

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	theme2 "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/comoyi/valheim-syncer-server/config"
	"github.com/comoyi/valheim-syncer-server/log"
	"github.com/comoyi/valheim-syncer-server/theme"
	"github.com/comoyi/valheim-syncer-server/util/fsutil"
	"github.com/comoyi/valheim-syncer-server/util/timeutil"
	"github.com/spf13/viper"
	"image/color"
	"os"
	"path/filepath"
)

var w fyne.Window
var c *fyne.Container
var myApp fyne.App

var msgContainer = widget.NewLabel("")
var announcementContent = ""
var dirStatusLed = canvas.NewRectangle(color.RGBA{R: 255, G: 0, B: 0, A: 255})

func StartGUI() {
	initUI()

	w.ShowAndRun()
}

func initUI() {
	initMainWindow()

	initMenu()
}

func initMainWindow() {
	windowTitle := fmt.Sprintf("%s-v%s", appName, versionText)

	myApp = app.NewWithID("com.comoyi.valheim-syncer-server")
	myApp.Settings().SetTheme(theme.CustomTheme)
	w = myApp.NewWindow(windowTitle)
	w.SetMaster()
	w.Resize(fyne.NewSize(800, 600))
	c = container.NewVBox()
	w.SetContent(c)

	initDir(c)
	initAnnouncement(c)
	initMsgContainer(c)
}

func initMenu() {
	firstMenu := fyne.NewMenu("操作")
	helpMenuItem := fyne.NewMenuItem("关于", func() {
		content := container.NewVBox()
		appInfo := widget.NewLabel(appName)
		content.Add(appInfo)
		versionInfo := widget.NewLabel(fmt.Sprintf("Version %v", versionText))
		content.Add(versionInfo)

		h := container.NewHBox()

		authorInfo := widget.NewLabel("Copyright © 2022 清新池塘")
		h.Add(authorInfo)
		linkInfo := widget.NewHyperlink(" ", nil)
		_ = linkInfo.SetURLFromString("https://github.com/comoyi/valheim-syncer-server")
		h.Add(linkInfo)
		content.Add(h)
		dialog.NewCustom("关于", "关闭", content, w).Show()
	})
	helpMenu := fyne.NewMenu("帮助", helpMenuItem)
	mainMenu := fyne.NewMainMenu(firstMenu, helpMenu)
	w.SetMainMenu(mainMenu)
}

func initDir(c *fyne.Container) {
	pathLabel := widget.NewLabel("Valheim文件夹 / MOD文件夹")
	pathInput := widget.NewLabel("")
	pathInput.SetText(config.Conf.Dir)

	selectBtnText := "选择文件夹"
	selectBtn := widget.NewButton(selectBtnText, func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				log.Debugf("select folder failed, err: %v\n", err)
				return
			}
			if uri == nil {
				log.Debugf("select folder cancelled\n")
				return
			}
			path := uri.Path()
			path = filepath.Clean(path)
			dialog.NewCustomConfirm("提示", "确定", "取消", widget.NewLabel("选择这个文件夹吗？\n"+path), func(b bool) {
				if b {
					pathInput.SetText(path)
					baseDir = path
					err := saveDirConfig(path)
					if err != nil {
						return
					}
					addMsgWithTime("文件夹设置为：" + path)
				}
			}, w).Show()
		}, w)
	})
	selectBtn.SetIcon(theme2.FolderIcon())

	dirStatusLed.SetMinSize(fyne.NewSize(5, 5))

	pathBox := container.NewVBox()
	pathBox.Add(pathLabel)
	pathBox.Add(dirStatusLed)
	pathBox.Add(pathInput)
	c2 := container.NewAdaptiveGrid(2)
	initManualInputBtn(c2, pathInput)
	c2.Add(selectBtn)
	pathBox.Add(c2)
	c.Add(pathBox)
}

func initManualInputBtn(c *fyne.Container, pathInput *widget.Label) {
	var manualInputDialog dialog.Dialog
	inputBtnText := "手动输入文件夹地址"
	inputBtn := widget.NewButton(inputBtnText, func() {
		manualPathInput := widget.NewEntry()
		tipLabel := widget.NewLabel("")
		box := container.NewVBox(manualPathInput, tipLabel)
		manualInputDialog = dialog.NewCustomConfirm("请输入文件夹地址", "确定", "取消", box, func(b bool) {
			if b {
				if manualPathInput.Text == "" {
					tipLabel.SetText("请输入文件夹地址")
					manualInputDialog.Show()
					return
				}
				path := filepath.Clean(manualPathInput.Text)
				exists, err := fsutil.Exists(path)
				if err != nil {
					tipLabel.SetText("文件夹地址检测失败")
					manualInputDialog.Show()
					return
				}
				if !exists {
					tipLabel.SetText("该文件夹不存在")
					manualInputDialog.Show()
					return
				}
				f, err := os.Stat(path)
				if err != nil {
					tipLabel.SetText("文件夹地址检测失败[2]")
					manualInputDialog.Show()
					return
				}
				if !f.IsDir() {
					tipLabel.SetText("请输入正确的文件夹地址")
					manualInputDialog.Show()
					return
				}

				pathInput.SetText(path)
				baseDir = path
				err = saveDirConfig(path)
				if err != nil {
					return
				}
				addMsgWithTime("文件夹设置为：" + path)
			}
		}, w)
		manualInputDialog.Resize(fyne.NewSize(700, 100))
		manualInputDialog.Show()
	})
	inputBtn.SetIcon(theme2.DocumentCreateIcon())

	c.Add(inputBtn)
}

func initAnnouncement(c *fyne.Container) {
	announcementLabel := widget.NewLabel("公告")
	var announcementInput = widget.NewMultiLineEntry()
	announcementInput.SetMinRowsVisible(7)
	announcementBtn := widget.NewButton("发布公告", func() {
		announcementContent = announcementInput.Text
		addMsgWithTime("发布公告成功")
	})
	announcementBtn.SetIcon(theme2.VolumeUpIcon())
	c.Add(announcementLabel)
	c.Add(announcementInput)
	c.Add(announcementBtn)
}

func initMsgContainer(*fyne.Container) {
	msgContainerScroll := container.NewScroll(msgContainer)
	msgContainerScroll.SetMinSize(fyne.NewSize(800, 200))
	c.Add(msgContainerScroll)
}

func saveDirConfig(path string) error {
	viper.Set("dir", path)
	err := config.SaveConfig()
	if err != nil {
		log.Debugf("save config failed, err: %+v\n", err)
		return err
	}
	return nil
}

func addMsgWithTime(msg string) {
	msg = fmt.Sprintf("%s %s", timeutil.GetCurrentDateTime(), msg)
	addMsg(msg)
}

func addMsg(msg string) {
	msgContainer.SetText(msg + "\n" + msgContainer.Text)
}

func setDirStatusLedRed() {
	if config.Conf.Gui == "OFF" {
		return
	}
	dirStatusLed.FillColor = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	dirStatusLed.Refresh()
}

func setDirStatusLedGreen() {
	if config.Conf.Gui == "OFF" {
		return
	}
	dirStatusLed.FillColor = color.RGBA{R: 0, G: 255, B: 0, A: 255}
	dirStatusLed.Refresh()
}
