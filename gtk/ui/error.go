package ui

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
)

func ShowErrorDlgFatal(txt string, parent *gtk.Window) {
	showErrorDlg(txt, true, parent)
}

func ShowErrorDlg(txt string, parent *gtk.Window) {
	showErrorDlg(txt, false, parent)
}

func showErrorDlg(txt string, fatal bool, parent *gtk.Window) {
	log.Printf("Error: %v", txt)

	dlg, _ := gtk.DialogNew()
	dlg.SetTitle("InsteadMan error")
	dlg.AddButton("Close", gtk.RESPONSE_ACCEPT)
	dlgBox, _ := dlg.GetContentArea()
	dlgBox.SetSpacing(6)

	lbl, _ := gtk.LabelNew(txt)
	lbl.SetMarginStart(6)
	lbl.SetMarginEnd(6)
	lbl.SetLineWrap(true)
	dlgBox.Add(lbl)
	lbl.Show()

	dlg.SetModal(true)
	dlg.SetPosition(gtk.WIN_POS_CENTER)
	dlg.SetResizable(false)

	if parent != nil {
		dlg.SetTransientFor(parent)
	}

	dlg.SetKeepAbove(true)
	dlg.Run()
	dlg.Destroy()
	if fatal {
		os.Exit(1)
	}
}