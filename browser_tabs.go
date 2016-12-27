package main

import (
	"log"
	"net/url"

	"github.com/mattn/go-gtk/gtk"
	"github.com/sg3des/vegevoice/webkit"
)

type Tab struct {
	label   *gtk.Label
	entry   *gtk.Entry
	webview *webkit.WebView

	vbox *gtk.VBox
	swin *gtk.ScrolledWindow
}

func (ui *UI) NewTab(addr string) *Tab {
	t := &Tab{}

	t.entry = gtk.NewEntry()
	t.entry.SetText(addr)

	t.webview = webkit.NewWebView()

	t.swin = gtk.NewScrolledWindow(nil, nil)
	t.swin.Add(t.webview)

	t.vbox = gtk.NewVBox(false, 0)
	t.vbox.PackStart(t.entry, false, false, 0)
	t.vbox.PackStart(t.swin, true, true, 0)

	t.label = gtk.NewLabel(addr)

	n := ui.notebook.AppendPage(t.vbox, t.label)
	ui.notebook.ShowAll()
	ui.notebook.SetCurrentPage(n)
	t.entry.GrabFocus()

	t.entry.Connect("activate", t.activate)

	if len(addr) > 0 {
		t.entry.Emit("activate")
	} else {
		t.label.SetText("New Tab")
	}

	ui.tabs = append(ui.tabs, t)
	return t
}

func (t *Tab) activate() {
	addr := t.entry.GetText()
	u, err := url.Parse(addr)
	if err != nil {
		log.Println(err)
		return
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	t.label.SetText(u.Path)
	t.entry.SetText(u.String())
	t.webview.LoadUri(u.String())
	t.webview.GrabFocus()
}

func (ui *UI) CloseCurrentTab() {
	n := ui.notebook.GetCurrentPage()
	ui.notebook.RemovePage(ui.tabs[n].swin, n)
	ui.tabs = append(ui.tabs[:n], ui.tabs[n+1:]...)
}

func (ui *UI) GetCurrentTab() *Tab {
	n := ui.notebook.GetCurrentPage()
	if n < 0 {
		return nil
	}
	return ui.tabs[n]
}
