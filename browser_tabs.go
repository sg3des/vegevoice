package main

import (
	"log"
	"net/url"
	"strings"
	"unsafe"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/gdkpixbuf"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"github.com/sg3des/vegevoice/webkit"

	"github.com/sg3des/vegevoice/addrs"
)

type Tab struct {
	tabbox  *gtk.EventBox
	favicon *gtk.Image
	label   *gtk.Label
	Pinned  bool

	idonChanged      int
	urlbar           *gtk.Entry
	urlbarCompletion *gtk.EntryCompletion
	urlbarHints      []string

	tabPopupMenu *gtk.Menu

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

	t.favicon = gtk.NewImage()
	// t.favicon.SetSizeRequest(8, 8)
	// t.favicon.
	t.label = gtk.NewLabel(addr)

	htabbox := gtk.NewHBox(false, 0)
	htabbox.PackStart(t.favicon, false, false, 0)
	htabbox.PackStart(t.label, true, true, 0)

	t.tabbox = gtk.NewEventBox()
	t.tabbox.Add(htabbox)
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

	t.initTabPopupMenu()

	if len(addr) > 0 {
		t.urlbar.Emit("activate")
	} else {
		t.label.SetText("New Tab")
	}

	t.urlbarCompletion.Connect("action-activated", t.onUrlbarCompetionActivated)
	t.idonChanged = t.urlbar.Connect("changed", t.onUrlbarChanged)

	ui.tabs = append(ui.tabs, t)
	return t
}

func (t *Tab) initTabPopupMenu() {
	itemPin := gtk.NewMenuItemWithLabel("Pin Tab")
	itemClose := gtk.NewMenuItemWithLabel("Close Tab")
	itemCloseOther := gtk.NewMenuItemWithLabel("Close Other Tabs")

	itemPin.Connect("activate", t.Pin)
	itemClose.Connect("activate", t.Close)
	itemCloseOther.Connect("activate", t.CloseOtherTabs)

	t.tabPopupMenu = gtk.NewMenu()
	t.tabPopupMenu.Add(itemPin)
	t.tabPopupMenu.Add(itemClose)
	t.tabPopupMenu.Add(itemCloseOther)
	t.tabPopupMenu.ShowAll()
}

func (t *Tab) Pin() {
	t.Pinned = true
	t.tabbox.SetSizeRequest(12, 12)
}

func (t *Tab) Close() {
	n := UI.notebook.PageNum(t.vbox)
	UI.CloseTab(n)
}

func (t *Tab) CloseOtherTabs() {
	min := 1
	for {
		for n, _t := range UI.tabs {
			if _t.label == t.label {
				UI.notebook.SetCurrentPage(n)
				continue
			}

			if _t.Pinned {
				min++
				continue
			}

			UI.CloseTab(n)
			break
		}

		if UI.notebook.GetNPages() == min {
			break
		}
	}
}

func (t *Tab) onLabelContextMenu(ctx *glib.CallbackContext) {
	arg := ctx.Args(0)
	event := *(**gdk.EventButton)(unsafe.Pointer(&arg))

	if event.Button == 3 {
		t.tabPopupMenu.Popup(nil, nil, nil, t.label, uint(arg), uint32(ctx.Args(1)))
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

	l := t.urlbar.GetPosition()
	substr = substr[:l+1]

	prevHints := t.urlbarHints

	t.urlbarHints = addrs.GetAddrs(substr)
	if len(t.urlbarHints) == 0 {
		return
	}

	for i, a := range t.urlbarHints {
		if i == 0 && l > 0 && l < len(a) && a[:l+1] == substr {
			t.urlbar.HandlerDisconnect(t.idonChanged)

			t.urlbar.SetPosition(0)
			t.urlbar.SetText(a)
			t.urlbar.SetPosition(l)

			t.idonChanged = t.urlbar.Connect("changed", t.onUrlbarChanged)
			continue
		}

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

	if uri := t.webview.GetIconUri(); len(uri) > 0 {
		iconpath := downloadIcon(uri)
		pix := gtk.NewImageFromFile(iconpath).GetPixbuf()
		pix = pix.ScaleSimple(12, 12, gdkpixbuf.INTERP_BILINEAR)
		t.favicon.SetFromPixbuf(pix)
	}
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
}

func (t *Tab) HistoryNext() {
	t.webview.GoForward()
	t.label.SetText(t.webview.GetTitle())
	t.urlbar.SetText(t.webview.GetUri())
}

func (ui *UserInterface) CloseCurrentTab() {
	n := ui.notebook.GetCurrentPage()
	ui.CloseTab(n)
}

func (ui *UserInterface) CloseTab(n int) {
	if len(ui.tabs) > 1 {
		if n == 0 {
			ui.notebook.SetCurrentPage(n + 1)
		} else {
			ui.notebook.SetCurrentPage(n - 1)
		}
	}

	ui.notebook.RemovePage(ui.tabs[n].vbox, n)

	ui.tabs[n] = nil
	ui.tabs = append(ui.tabs[:n], ui.tabs[n+1:]...)

	if len(ui.tabs) == 0 {
		gtk.MainQuit()
	}
}

func (ui *UserInterface) GetCurrentTab() *Tab {
	n := ui.notebook.GetCurrentPage()
	if n < 0 {
		return nil
	}
	return ui.tabs[n]
}
