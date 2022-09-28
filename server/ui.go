package server

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/comoyi/valheim-syncer-server/config"
	"github.com/comoyi/valheim-syncer-server/log"
	"github.com/comoyi/valheim-syncer-server/theme"
	"github.com/comoyi/valheim-syncer-server/utils/timeutil"
	"github.com/spf13/viper"
	"path/filepath"
)

var w fyne.Window
var c *fyne.Container
var myApp fyne.App

var msgContainer = widget.NewLabel("")
var announcementContent = ""

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
				}
			}, w).Show()
		}, w)
	})

	pathBox := container.NewVBox()
	pathBox.Add(pathLabel)
	pathBox.Add(pathInput)
	c2 := container.NewAdaptiveGrid(1)
	c2.Add(selectBtn)
	pathBox.Add(c2)
	c.Add(pathBox)
}

func initAnnouncement(c *fyne.Container) {
	announcementLabel := widget.NewLabel("公告")
	var announcementInput = widget.NewMultiLineEntry()
	announcementInput.SetMinRowsVisible(7)
	announcementBtn := widget.NewButton("发布公告", func() {
		announcementContent = announcementInput.Text
		addMsgWithTime("发布公告成功")
	})
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
