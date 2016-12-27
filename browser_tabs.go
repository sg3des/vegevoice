package main

import (
	"log"
	"net/url"
	"strings"

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
	t.entry.Connect("changed", t.changed)
	t.webview.Connect("load-progress-changed", t.loadProgressChanged)
	t.webview.Connect("load-finished", t.loadFinished)

	if len(addr) > 0 {
		t.entry.Emit("activate")
	} else {
		t.label.SetText("New Tab")
	}

	ui.tabs = append(ui.tabs, t)
	return t
}

func (t *Tab) activate() {
	saddr := t.entry.GetText()
	if !strings.Contains(saddr, ".") {
		saddr = "http://google.com/search?q=" + saddr
	}
	addr := t.parseAddr(saddr)
	t.OpenUrl(addr)
}

func (t *Tab) parseAddr(reqaddr string) *url.URL {
	u, err := url.Parse(reqaddr)
	if err != nil {
		log.Println(err)
		return u
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	return u
}

func (t *Tab) OpenUrl(addr *url.URL) {
	t.label.SetText(addr.Path)
	t.entry.SetText(addr.String())
	t.webview.LoadUri(addr.String())
	t.webview.GrabFocus()
}

func (t *Tab) changed() {

}

func (t *Tab) loadProgressChanged() {
	t.label.SetText(t.webview.GetTitle())
	if !t.entry.HasFocus() {
		t.entry.SetText(t.webview.GetUri())
	}
	// log.Println(t.webview.GetProgress())
}

func (t *Tab) loadFinished() {
	t.label.SetText(t.webview.GetTitle())
	if !t.entry.HasFocus() {
		t.entry.SetText(t.webview.GetUri())
	}
}

func (t *Tab) HistoryBack() {
	t.webview.GoBack()
	t.label.SetText(t.webview.GetTitle())
	t.entry.SetText(t.webview.GetUri())
	// if t.historyN <= 0 {
	// 	return
	// }
	// t.historyN--
	// t.OpenUrl(t.history[t.historyN])
}

func (t *Tab) HistoryNext() {
	t.webview.GoForward()
	t.label.SetText(t.webview.GetTitle())
	t.entry.SetText(t.webview.GetUri())
	// if t.historyN >= len(t.history)-1 {
	// 	return
	// }
	// t.historyN++
	// t.OpenUrl(t.history[t.historyN])
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
