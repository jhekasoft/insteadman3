package primary

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/jhekasoft/insteadman3/cmd/insteadman-fyne/data"
	"github.com/jhekasoft/insteadman3/core/configurator"
	"github.com/jhekasoft/insteadman3/core/manager"
)

type GameInfoScreen struct {
	win    fyne.Window
	m      *manager.Manager
	c      *configurator.Configurator
	Screen fyne.CanvasObject
	game   *manager.Game

	// Widgets
	title         *widget.Label
	desc          *widget.Label
	version       *widget.Label
	lang          *widget.Label
	repository    *widget.Label
	size          *widget.Label
	container     *widget.SplitContainer
	image         *widget.Icon
	hyperlink     *widget.Hyperlink
	installButton *widget.Button
	runButton     *widget.Button
	deleteButton  *widget.Button
}

func (scr *GameInfoScreen) UpdateInfo(g *manager.Game) {
	scr.game = g

	scr.title.SetText(g.Title)
	scr.desc.SetText(g.Description)

	// Labels
	scr.version.SetText(g.Version)
	scr.lang.SetText(strings.Join(g.Languages, ", "))
	scr.repository.SetText(g.RepositoryName)

	// URL
	if g.Descurl != "" {
		scr.hyperlink.SetURLFromString(g.Descurl)
		scr.hyperlink.Show()
	}

	scr.size.SetText(g.HumanSize())
	scr.size.Show()

	// Buttons
	// TODO: add Update button
	if g.Installed {
		scr.installButton.Hide()
		scr.runButton.Show()
		scr.deleteButton.Show()
	} else {
		scr.installButton.Show()
		scr.runButton.Hide()
		scr.deleteButton.Hide()
	}

	var icon fyne.Resource = data.InsteadManLogo
	var b []byte = nil

	scr.image.SetResource(icon)

	fileName, e := scr.m.GetGameImage(g)
	if e == nil {
		iconFile, e := os.Open(scr.c.DataResourcePath(fileName))
		if e == nil {
			r := bufio.NewReader(iconFile)

			b, e = ioutil.ReadAll(r)
		}

		if e != nil {
			// dialog.ShowError(e, scr.Window)
			fmt.Printf("Error: %v\n", e)
		} else {
			icon = fyne.NewStaticResource("game_"+g.Name, b)
			scr.image.SetResource(icon)
		}
	}

	// if scr.Container != nil {
	// 	scr.Container.Refresh()
	// }
}

func NewGameInfoScreen(
	win fyne.Window,
	m *manager.Manager,
	c *configurator.Configurator,
	onRefresh func()) *GameInfoScreen {
	scr := GameInfoScreen{win: win, m: m, c: c}

	scr.image = widget.NewIcon(data.InsteadManLogo)
	scr.title = widget.NewLabelWithStyle("InsteadMan", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	scr.desc = widget.NewLabel("Выберите игру слева в списке")
	scr.desc.Wrapping = fyne.TextWrapWord
	scr.size = widget.NewLabel("")
	scr.size.Hide()

	descScroll := widget.NewVScrollContainer(
		scr.desc,
	)
	// descScroll.SetMinSize(fyne.NewSize(0, 100))

	scr.hyperlink = widget.NewHyperlink("Website", nil)
	scr.hyperlink.Hide()
	scr.installButton = widget.NewButtonWithIcon("Install", theme.ContentAddIcon(), func() {
		progDialog := dialog.NewProgress(scr.game.Title, "Installing...", scr.win)
		progDialog.Show()
		err := scr.m.InstallGame(scr.game, func(size uint64) {
			percents := float64(size) / float64(scr.game.Size)
			progDialog.SetValue(percents)
			if float64(size) >= float64(scr.game.Size) {
				progDialog.SetValue(1)
				progDialog.Hide()
			}
		})

		if err != nil {
			progDialog.Hide()
			dialog.ShowError(err, scr.win)
			return
		}

		scr.game.Installed = true
		scr.UpdateInfo(scr.game)

		if onRefresh != nil {
			onRefresh()
		}
	})
	scr.installButton.Style = widget.PrimaryButton
	scr.installButton.Hide()
	scr.runButton = widget.NewButtonWithIcon("Run", theme.MediaPlayIcon(), func() {
		scr.m.RunGame(scr.game)
	})
	scr.runButton.Style = widget.PrimaryButton
	scr.runButton.Hide()
	scr.deleteButton = widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
		scr.m.RemoveGame(scr.game)

		// TODO: Check error
		scr.game.Installed = false
		scr.UpdateInfo(scr.game)
		if onRefresh != nil {
			onRefresh()
		}
	})
	scr.deleteButton.Hide()
	scr.version = widget.NewLabel("")
	scr.lang = widget.NewLabel("")
	scr.repository = widget.NewLabel("")
	// scr.repository.Wrapping = fyne.TextWrapWord
	bottomInfoScroll := widget.NewHScrollContainer(container.NewHBox(
		scr.hyperlink,
		scr.version,
		scr.lang,
		scr.size,
		scr.repository,
	))
	buttonsContainer := container.NewHBox(scr.installButton, scr.runButton, scr.deleteButton)
	bottomContainer := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, buttonsContainer, nil),
		buttonsContainer,
		bottomInfoScroll,
	)

	contentContainer := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(scr.title, bottomContainer, nil, nil),
		descScroll,
		scr.title,
		bottomContainer,
	)

	scr.container = widget.NewVSplitContainer(scr.image, contentContainer)

	scr.Screen = fyne.NewContainerWithLayout(
		layout.NewBorderLayout(nil, nil, nil, nil),
		scr.container,
	)

	return &scr
}