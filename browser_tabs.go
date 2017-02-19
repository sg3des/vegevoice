package main

import (
	"log"
	"net/url"
	"strings"
	"unsafe"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/sg3des/vegevoice/webkit"

	"github.com/sg3des/vegevoice/addrs"
)

type Tab struct {
	tabbox *gtk.EventBox
	label  *gtk.Label

	urlbar           *gtk.Entry
	urlbarCompletion *gtk.EntryCompletion
	urlbarHints      []string

	webview *webkit.WebView

	vbox *gtk.VBox
	swin *gtk.ScrolledWindow
}

func (ui *UserInterface) NewTab(addr string) *Tab {
	t := &Tab{}

	t.urlbarCompletion = gtk.NewEntryCompletion()
	urlbarListStore := gtk.NewListStore(glib.G_TYPE_STRING)
	t.urlbarCompletion.SetModel(&urlbarListStore.TreeModel)
	t.urlbarCompletion.SetTextColumn(0)

	t.urlbar = gtk.NewEntry()
	t.urlbar.SetCompletion(t.urlbarCompletion)

	t.webview = webkit.NewWebView()
	ApplySettings(t.webview)

	t.swin = gtk.NewScrolledWindow(nil, nil)
	t.swin.Add(t.webview)

	t.vbox = gtk.NewVBox(false, 0)
	t.vbox.PackStart(t.urlbar, false, false, 0)
	t.vbox.PackStart(t.swin, true, true, 0)

	t.label = gtk.NewLabel(addr)
	t.tabbox = gtk.NewEventBox()
	t.tabbox.Add(t.label)
	t.tabbox.ShowAll()

	n := ui.notebook.AppendPage(t.vbox, t.tabbox)
	ui.notebook.ShowAll()
	ui.notebook.SetCurrentPage(n)
	t.urlbar.GrabFocus()

	t.urlbar.Connect("activate", t.onUrlbarActivate)
	t.webview.Connect("load-progress-changed", t.onLoadProgressChanged)
	t.webview.Connect("load-finished", t.onLoadFinished)
	t.webview.Connect("create-web-view", t.onCreateWebView)
	t.webview.Connect("web-view-ready", t.onWebViewReady)
	t.tabbox.Connect("button-release-event", t.onLabelContextMenu)

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

func (t *Tab) onLabelContextMenu(ctx *glib.CallbackContext) {
	arg := ctx.Args(0)
	event := *(**gdk.EventButton)(unsafe.Pointer(&arg))

	if event.Button == 3 {
		m := gtk.NewMenu()
		m.Add(gtk.NewMenuItemWithLabel("label"))
		m.Add(gtk.NewMenuItemWithLabel("label - 2"))
		m.ShowAll()
		m.Popup(nil, nil, nil, t.label, uint(ctx.Args(0)), uint32(ctx.Args(1)))
	}
}

func (t *Tab) onCreateWebView() interface{} {
	newtab := UI.NewTab("")
	return newtab.webview.GetWebView()
}

func (t *Tab) onWebViewReady(ctx *glib.CallbackContext) {
	log.Println("onWebViewReady")
	log.Println(ctx.Args(0))
	// log.Println(ctx.Data())
}

//onUrlbarChanged signal changed on urlbar entry
func (t *Tab) onUrlbarChanged() {
	substr := t.urlbar.GetText()
	if len(substr) == 0 {
		return
	}

	prevHints := t.urlbarHints

	t.urlbarHints = addrs.GetAddrs(substr)
	if len(t.urlbarHints) == 0 {
		return
	}

	for i, a := range t.urlbarHints {
		if i < len(prevHints) && prevHints[i] == a {
			continue
		}

		if i < len(prevHints) {
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
	if i >= len(t.urlbarHints) {
		log.Println("WARNING! iter more then hint urls")
		return ""
	}

	return t.urlbarHints[i]
}

func (t *Tab) onUrlbarActivate() {
	saddr := t.urlbar.GetText()
	if splitted := strings.Split(saddr, "."); len(splitted) < 2 || len(splitted[1]) == 0 {
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
	title := t.webview.GetTitle()
	uri := t.webview.GetUri()
	if len(title) == 0 || len(uri) == 0 {
		return
	}

	t.label.SetText(title)
	if !t.urlbar.HasFocus() {
		t.urlbar.SetText(uri)
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

func (ui *UserInterface) CloseCurrentTab() {
	n := ui.notebook.GetCurrentPage()
	if len(ui.tabs) > 1 {
		if n == 0 {
			ui.notebook.SetCurrentPage(n + 1)
		} else {
			ui.notebook.SetCurrentPage(n - 1)
		}
	}

	ui.notebook.RemovePage(ui.tabs[n].swin, n)

	ui.tabs[n] = nil
	ui.tabs = append(ui.tabs[:n], ui.tabs[n+1:]...)
}

func (ui *UserInterface) GetCurrentTab() *Tab {
	n := ui.notebook.GetCurrentPage()
	if n < 0 {
		return nil
	}
	return ui.tabs[n]
}
