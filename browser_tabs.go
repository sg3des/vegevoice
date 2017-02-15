package main

import (
	"log"
	"net/url"
	"strings"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/mattn/go-webkit/webkit"

	"github.com/sg3des/vegevoice/addrs"
)

type Tab struct {
	label *gtk.Label

	urlbar           *gtk.Entry
	urlbarCompletion *gtk.EntryCompletion
	urlbarHints      *gtk.ListStore
	urlbarUrlHints   []string

	webview *webkit.WebView

	vbox *gtk.VBox
	swin *gtk.ScrolledWindow
}

func (ui *UI) NewTab(addr string) *Tab {
	t := &Tab{}

	t.urlbarHints = gtk.NewListStore(glib.G_TYPE_STRING)
	t.urlbarCompletion = gtk.NewEntryCompletion()
	t.urlbarCompletion.SetModel(&t.urlbarHints.TreeModel)
	t.urlbarCompletion.SetTextColumn(0)

	t.urlbar = gtk.NewEntry()
	t.urlbar.SetCompletion(t.urlbarCompletion)

	t.webview = webkit.NewWebView()

	t.swin = gtk.NewScrolledWindow(nil, nil)
	t.swin.Add(t.webview)

	t.vbox = gtk.NewVBox(false, 0)
	t.vbox.PackStart(t.urlbar, false, false, 0)
	t.vbox.PackStart(t.swin, true, true, 0)

	t.label = gtk.NewLabel(addr)

	n := ui.notebook.AppendPage(t.vbox, t.label)
	ui.notebook.ShowAll()
	ui.notebook.SetCurrentPage(n)
	t.urlbar.GrabFocus()

	t.urlbar.Connect("activate", t.onUrlbarActivate)
	t.webview.Connect("load-progress-changed", t.onLoadProgressChanged)
	t.webview.Connect("load-finished", t.onLoadFinished)

	if len(addr) > 0 {
		t.urlbar.Emit("activate")
	} else {
		t.label.SetText("New Tab")
	}

	t.urlbarCompletion.Connect("action-activated", t.onUrlbarCompetionActivated)
	t.urlbar.Connect("changed", t.onUrlbarChanged)

	ui.tabs = append(ui.tabs, t)
	return t
}

//onUrlbarChanged signal changed on urlbar entry
func (t *Tab) onUrlbarChanged() {
	substr := t.urlbar.GetText()
	if len(substr) == 0 {
		return
	}

	prevcount := len(t.urlbarUrlHints)

	t.urlbarUrlHints = addrs.GetAddrs(substr)
	if len(t.urlbarUrlHints) == 0 {
		return
	}

	for i, a := range t.urlbarUrlHints {
		if i <= prevcount {
			t.urlbarCompletion.DeleteAction(i)
		}

		t.urlbarCompletion.InsertActionText(i, a)
	}

}

func (t *Tab) onUrlbarCompetionActivated(ctx *glib.CallbackContext) {
	addr := t.getUrlFromHints(int(ctx.Args(0)))

	t.urlbar.SetText(addr)
	t.urlbar.Emit("activate")
}

func (t *Tab) getUrlFromHints(i int) string {
	if i >= len(t.urlbarUrlHints) {
		log.Println("WARNING! iter more then hint urls")
		return ""
	}

	return t.urlbarUrlHints[i]
}

func (t *Tab) onUrlbarActivate() {
	saddr := t.urlbar.GetText()
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
	t.urlbar.SetText(addr.String())
	t.webview.LoadUri(addr.String())
	t.webview.GrabFocus()
}

func (t *Tab) onLoadProgressChanged() {
	t.label.SetText(t.webview.GetTitle())
	if !t.urlbar.HasFocus() {
		t.urlbar.SetText(t.webview.GetUri())
	}
	// log.Println(t.webview.GetProgress())
}

func (t *Tab) onLoadFinished() {
	t.label.SetText(t.webview.GetTitle())
	if !t.urlbar.HasFocus() {
		t.urlbar.SetText(t.webview.GetUri())
	}
}

func (t *Tab) HistoryBack() {
	t.webview.GoBack()
	t.label.SetText(t.webview.GetTitle())
	t.urlbar.SetText(t.webview.GetUri())
	// if t.historyN <= 0 {
	// 	return
	// }
	// t.historyN--
	// t.OpenUrl(t.history[t.historyN])
}

func (t *Tab) HistoryNext() {
	t.webview.GoForward()
	t.label.SetText(t.webview.GetTitle())
	t.urlbar.SetText(t.webview.GetUri())
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
