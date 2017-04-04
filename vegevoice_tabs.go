package main

import (
	"log"
	"net/url"
	"strings"
	"unsafe"

	"github.com/mattn/go-gtk/gdk"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"

	"github.com/sg3des/vegevoice/urlstorage"
	"github.com/sg3des/vegevoice/webkit"
)

type Tab struct {
	tabbox       *gtk.VBox
	favicon      *gtk.Image
	label        *gtk.Label
	progressbar  *gtk.ProgressBar
	tabPopupMenu *gtk.Menu
	Pinned       bool

	idonChanged      int
	urlbar           *gtk.Entry
	urlbarCompletion *gtk.EntryCompletion
	urlbarHints      []string
	urlbarHint       string

	webview *webkit.WebView

	findbox     *gtk.Table
	findbar     *gtk.Entry
	findbtnCS   *gtk.ToggleButton
	findbtnNext *gtk.Button
	findbtnPrev *gtk.Button

	vbox *gtk.VBox
	swin *gtk.ScrolledWindow
}

func (ui *UserInterface) NewTab(reqURL string) *Tab {
	t := &Tab{}

	// tab
	t.favicon = gtk.NewImage()
	t.favicon.SetSizeRequest(16, 16)
	t.label = gtk.NewLabel(reqURL)

	tabtable := gtk.NewTable(1, 2, false)
	tabtable.Attach(t.favicon, 0, 1, 0, 1, gtk.FILL, gtk.FILL, 1, 1)
	tabtable.Attach(t.label, 1, 2, 0, 1, gtk.FILL, gtk.FILL, 0, 0)

	eventbox := gtk.NewEventBox()
	eventbox.Add(tabtable)
	eventbox.ShowAll()

	t.progressbar = gtk.NewProgressBar()
	t.progressbar.SetSizeRequest(4, 4)

	t.tabbox = gtk.NewVBox(false, 0)
	t.tabbox.Add(eventbox)
	t.tabbox.PackEnd(t.progressbar, false, true, 0)

	//urlbar
	t.urlbarCompletion = gtk.NewEntryCompletion()
	urlbarListStore := gtk.NewListStore(glib.G_TYPE_STRING)
	t.urlbarCompletion.SetModel(&urlbarListStore.TreeModel)
	t.urlbarCompletion.SetTextColumn(0)

	t.urlbar = gtk.NewEntry()
	t.urlbar.SetText(reqURL)
	t.urlbar.SetCompletion(t.urlbarCompletion)

	//webview
	t.webview = webkit.NewWebView()
	ApplySettings(t.webview)

	t.swin = gtk.NewScrolledWindow(nil, nil)
	t.swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	t.swin.SetShadowType(gtk.SHADOW_IN)
	t.swin.Add(t.webview)

	//findbar
	t.findbar = gtk.NewEntry()
	t.findbar.Connect("changed", func() { t.onSearch(true) })

	t.findbtnCS = gtk.NewToggleButtonWithLabel("Aa")
	t.findbtnCS.Clicked(func() { t.onSearch(true) })

	t.findbtnNext = gtk.NewButton()
	t.findbtnNext.SetImage(gtk.NewArrow(gtk.ARROW_RIGHT, gtk.SHADOW_NONE))
	t.findbtnNext.Clicked(func() { t.onSearch(true) })

	t.findbtnPrev = gtk.NewButton()
	t.findbtnPrev.SetImage(gtk.NewArrow(gtk.ARROW_LEFT, gtk.SHADOW_NONE))
	t.findbtnPrev.Clicked(func() { t.onSearch(false) })

	t.findbox = gtk.NewTable(1, 4, false)
	t.findbox.Attach(t.findbtnCS, 0, 1, 0, 1, gtk.FILL, gtk.FILL, 0, 0)
	t.findbox.Attach(t.findbar, 1, 2, 0, 1, gtk.EXPAND|gtk.FILL, gtk.FILL, 0, 0)
	t.findbox.Attach(t.findbtnPrev, 2, 3, 0, 1, gtk.FILL, gtk.FILL, 0, 0)
	t.findbox.Attach(t.findbtnNext, 3, 4, 0, 1, gtk.FILL, gtk.FILL, 0, 0)

	// t.findbox = gtk.NewHBox(false, 0)
	// t.findbox.Add(t.findbar)
	// t.findbox.PackEnd(t.findbtnCS, false, false, 0)
	// t.findbox.PackEnd(t.findbtnNext, false, false, 0)
	// t.findbox.PackEnd(t.findbtnPrev, false, false, 0)

	//main container
	t.vbox = gtk.NewVBox(false, 0)
	t.vbox.PackStart(t.urlbar, false, false, 0)
	t.vbox.PackStart(t.swin, true, true, 0)
	t.vbox.PackEnd(t.findbox, false, false, 0)

	//notebook
	ui.tabs = append(ui.tabs, t)
	n := ui.notebook.AppendPage(t.vbox, t.tabbox)
	ui.notebook.ShowAll()
	ui.notebook.SetCurrentPage(n)
	t.urlbar.GrabFocus()

	t.progressbar.SetVisible(false)
	t.findbox.SetVisible(false)

	t.urlbar.Connect("activate", t.onUrlbarActivate)
	t.webview.Connect("load-progress-changed", t.onLoadProgressChanged)
	t.webview.Connect("load-finished", t.onLoadFinished)
	t.webview.Connect("download-requested", func() { log.Println("download") })
	t.webview.Connect("create-web-view", t.onCreateWebView)
	t.webview.Connect("icon-loaded", t.onIconLoaded)
	t.tabbox.Connect("button-release-event", t.onLabelContextMenu)

	t.initTabPopupMenu()

	t.urlbarCompletion.Connect("action-activated", t.onUrlbarCompetionActivated)
	t.idonChanged = t.urlbar.Connect("changed", t.onUrlbarChanged)

	if len(reqURL) > 0 {
		t.urlbar.Emit("activate")
	} else {
		t.label.SetText("New Tab")
	}

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
	u := t.urlbar.GetText()
	toPage := len(urlstorage.GetPinnedTabs())

	if t.Pinned {
		t.Pinned = false
		t.label.SetVisible(true)
		urlstorage.DelPinnedTab(u)
	} else {
		t.Pinned = true
		t.label.SetVisible(false)
		urlstorage.AddPinnedTab(u)
	}

	t.Reorder(toPage)
	ui.homogenousTabs()
}

func (t *Tab) Reorder(to int) {
	n := ui.notebook.PageNum(t.vbox)
	ui.notebook.ReorderChild(t.vbox, to)
	ui.tabs[to], ui.tabs[n] = ui.tabs[n], ui.tabs[to]
}

func (t *Tab) Close() {
	n := ui.notebook.PageNum(t.vbox)
	ui.CloseTab(n)
}

func (t *Tab) CloseOtherTabs() {
	min := 1
	for {
		for n, _t := range ui.tabs {
			if _t.label == t.label {
				ui.notebook.SetCurrentPage(n)
				continue
			}

			if _t.Pinned {
				min++
				continue
			}

			ui.CloseTab(n)
			break
		}

		if ui.notebook.GetNPages() == min {
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
	newtab := ui.NewTab("")
	return newtab.webview.GetWebView()
}

//onUrlbarChanged signal changed on urlbar entry
func (t *Tab) onUrlbarChanged() {
	if !t.urlbar.HasFocus() {
		return
	}

	substr := t.urlbar.GetText()

	var right string
	l := t.urlbar.GetPosition()
	if l+1 < len(substr) {
		right = substr[l+1:]
		substr = substr[:l+1]
	}

	prevHints := t.urlbarHints
	t.urlbarHints = urlstorage.GetURLs(substr)
	if len(t.urlbarHints) == 0 {

		//clear inline tail
		if len(right) > 0 && len(t.urlbarHint) > 0 && right == t.urlbarHint {
			t.urlbar.HandlerDisconnect(t.idonChanged)
			t.urlbar.SetText(substr)
			t.urlbarHint = ""
			t.idonChanged = t.urlbar.Connect("changed", t.onUrlbarChanged)
		}

		//delete completaions
		for i := range prevHints {
			t.urlbarCompletion.DeleteAction(i)
		}

		return
	}

	for i, a := range t.urlbarHints {

		//inline completion
		if i == 0 && l > 0 && l < len(a) && a[:l+1] == substr {
			t.urlbar.HandlerDisconnect(t.idonChanged)

			t.urlbar.SetPosition(0)
			t.urlbar.SetText(a)
			t.urlbarHint = a[l+1:]
			t.urlbar.SetPosition(l)

			t.idonChanged = t.urlbar.Connect("changed", t.onUrlbarChanged)
			// continue
		}

		if i < len(prevHints) && prevHints[i] == a {
			continue
		}

		if i < len(prevHints) {
			t.urlbarCompletion.DeleteAction(i)
		}

		//popup completion
		t.urlbarCompletion.InsertActionText(i, a)
	}

}

func (t *Tab) onUrlbarCompetionActivated(ctx *glib.CallbackContext) {
	u := t.getUrlFromHints(int(ctx.Args(0)))

	t.urlbar.SetText(u)
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
	reqURL := t.urlbar.GetText()
	if splitted := strings.Split(reqURL, "."); len(splitted) < 2 || len(splitted[1]) == 0 {
		reqURL = "https://google.com/search?q=" + reqURL
	}
	u := t.parseURL(reqURL)
	t.OpenUrl(u)
}

func (t *Tab) parseURL(reqURL string) *url.URL {
	u, err := url.Parse(reqURL)
	if err != nil {
		log.Println(err)
		return u
	}

	if u.Scheme == "" {
		u.Scheme = "http"
	}

	return u
}

func (t *Tab) OpenUrl(u *url.URL) {
	log.Println("open URL:", u.String())
	t.label.SetText(u.String())
	t.urlbar.SetText(u.String())

	// client := &http.Client{}
	// req, err := http.NewRequest("GET", u.String(), nil)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// useragent, ok := conf.Webkit["user-agent"]
	// if !ok {
	// 	useragent = "Go-http-client/2.0"
	// }
	// req.Header.Set("User-Agent", useragent.(string))
	// resp, err := client.Do(req)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// // resp, err := http.Get(u.String())
	// // if err != nil {
	// // 	log.Println(err)
	// // 	return
	// // }
	// data, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// log.Println(len(data))

	t.webview.LoadUri(u.String())
	// t.webview.LoadString(string(data), "text/html", "utf-8", u.String())
	// t.webview.LoadHtmlString(string(data), u.String())
	t.webview.GrabFocus()
}

func (t *Tab) onLoadProgressChanged() {
	// log.Println("onLoadProgressChanged", t.webview.GetProgress())
	t.progressbar.SetVisible(true)
	t.progressbar.SetFraction(t.webview.GetProgress())

	if title := t.webview.GetTitle(); len(title) > 0 {
		t.label.SetText(title)
	}

	if !t.urlbar.HasFocus() {
		t.urlbar.SetText(t.webview.GetUri())
	}
}

func (t *Tab) onLoadFinished() {
	// log.Println("onLoadFinished")
	t.progressbar.SetVisible(false)

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

func (t *Tab) onIconLoaded(ctx *glib.CallbackContext) {
	u, err := url.Parse(t.webview.GetUri())
	if err != nil {
		log.Println(err)
		return
	}

	favicon := GetFavicon(u.Hostname(), t.webview.GetIconUri())
	if favicon != nil {
		t.favicon.SetFromPixbuf(favicon)
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

func (t *Tab) Reload() {
	t.webview.Reload()
}

func (t *Tab) onSearch(next bool) {
	text := t.findbar.GetText()
	if len(text) == 0 {
		return
	}

	// var next = true
	// if data := ctx.Data(); data != nil {
	// 	next = data.(bool)
	// }

	t.webview.UnmarkTextMatches()
	t.webview.SearchText(text, t.findbtnCS.GetActive(), next, true)

	t.webview.MarkTextMatches(text, t.findbtnCS.GetActive(), 128)
	t.webview.SetHighlightTextMatches(true)
}
