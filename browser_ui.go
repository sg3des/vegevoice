package main

import "github.com/mattn/go-gtk/gtk"

var (
	width  int
	height int
)

type UI struct {
	window *gtk.Window
	vbox   *gtk.VBox

	accelGroup  *gtk.AccelGroup
	actionGroup *gtk.ActionGroup

	menubar  *gtk.Widget
	notebook *gtk.Notebook
	tabs     []*Tab
}

func CreateUi() *UI {
	ui := &UI{}
	ui.window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	ui.window.SetSizeRequest(600, 600)
	ui.window.SetTitle("webkit")
	ui.window.Connect("destroy", ui.Quit)
	ui.window.Connect("check-resize", ui.windowResize)

	ui.menubar = ui.createMenubar()
	ui.notebook = gtk.NewNotebook()

	ui.vbox = gtk.NewVBox(false, 0)
	ui.vbox.PackStart(ui.menubar, true, true, 0)
	ui.vbox.PackStart(ui.notebook, true, true, 0)

	ui.window.Add(ui.vbox)
	ui.window.ShowAll()

	ui.menubar.SetVisible(false)

	return ui
}

func (ui *UI) createMenubar() *gtk.Widget {

	UIxml := `
<ui>
	<menubar name='MenuBar'>
		<menu action='File'>
			<menuitem action='NewTab' />
			<menuitem action='CloseTab' />
			<menuitem action='OpenUrl' />
			<menuitem action='Back' />
			<menuitem action='Next' />
			<separator />
			<menuitem action='Quit' />
		</menu>

		<menu action='Edit'>
			<menuitem action='Find'/>
			<menuitem action='FindNext'/>
			<menuitem action='FindPrev'/>
			<separator />
			<menuitem action='Replace'/>
			<menuitem action='ReplaceOne'/>
			<menuitem action='ReplaceAll'/>
		</menu>

		<menu name='View' action='View'>
			<menuitem action='Menubar'/>
		</menu>

	</menubar>
</ui>
`
	uiManager := gtk.NewUIManager()
	uiManager.AddUIFromString(UIxml)

	ui.accelGroup = uiManager.GetAccelGroup()
	ui.window.AddAccelGroup(ui.accelGroup)

	ui.actionGroup = gtk.NewActionGroup("my_group")
	uiManager.InsertActionGroup(ui.actionGroup, 0)

	// File
	ui.actionGroup.AddAction(gtk.NewAction("File", "File", "", ""))

	ui.newAction("NewTab", "New Tab", "<control>t", ui.newTab)
	ui.newAction("CloseTab", "Close Tab", "<control>w", ui.closeTab)
	ui.newAction("OpenUrl", "Open URL", "<control>l", ui.focusurl)
	ui.newAction("Back", "Back", "<Alt>Left", ui.back)
	ui.newAction("Next", "Next", "<Alt>Right", ui.next)
	ui.newActionStock("Quit", gtk.STOCK_QUIT, "", ui.Quit)

	// Edit
	ui.actionGroup.AddAction(gtk.NewAction("Edit", "Edit", "", ""))

	ui.newActionStock("Find", gtk.STOCK_FIND, "", ui.ShowFindbar)
	ui.newAction("FindNext", "Find Next", "F3", ui.FindNext)
	ui.newAction("FindPrev", "Find Previous", "<shift>F3", ui.FindPrev)

	ui.newActionStock("Replace", gtk.STOCK_FIND_AND_REPLACE, "<control>h", ui.ShowReplbar)
	ui.newAction("ReplaceOne", "Replace One", "<control><shift>h", ui.ReplaceOne)
	ui.newAction("ReplaceAll", "Replace All", "<control><alt>Return", ui.ReplaceAll)

	// View
	ui.actionGroup.AddAction(gtk.NewAction("View", "View", "", ""))
	// ui.actionGroup.AddAction(gtk.NewAction("Encoding", "Encoding", "", ""))

	ui.newToggleAction("Menubar", "Menubar", "<control>M", false, ui.ToggleMenuBar)

	return uiManager.GetWidget("/MenuBar")
}

func (ui *UI) newAction(dst, label, accel string, f func()) {
	action := gtk.NewAction(dst, label, "", "")
	action.Connect("activate", f)
	ui.actionGroup.AddActionWithAccel(action, accel)
}

func (ui *UI) newActionStock(dst, stock, accel string, f func()) {
	action := gtk.NewAction(dst, "", "", stock)
	action.Connect("activate", f)
	ui.actionGroup.AddActionWithAccel(action, accel)
}

func (ui *UI) newToggleAction(dst, label, accel string, state bool, f func()) {
	action := gtk.NewToggleAction(dst, label, "", "")
	action.SetActive(state)
	action.Connect("activate", f)
	ui.actionGroup.AddActionWithAccel(&action.Action, accel)
}

// actions
func (ui *UI) windowResize() {
	ui.window.GetSize(&width, &height)
	ui.notebook.SetSizeRequest(width, height)
	ui.homogenousTabs()
}

func (ui *UI) homogenousTabs() {
	if len(ui.tabs) == 0 {
		return
	}

	tabwidth := (width - len(ui.tabs)*6) / len(ui.tabs)
	leftwidth := (width - len(ui.tabs)*6) % len(ui.tabs)

	for _, t := range ui.tabs {
		if leftwidth > 0 {
			t.label.SetSizeRequest(tabwidth+1, 12)
			leftwidth--
		} else {
			t.label.SetSizeRequest(tabwidth, 12)
		}
	}
}

func (ui *UI) newTab() {
	ui.NewTab("")
}
func (ui *UI) closeTab() {
	ui.CloseCurrentTab()

	if len(ui.tabs) == 0 {
		gtk.MainQuit()
	}
}
func (ui *UI) focusurl() {
	ui.GetCurrentTab().urlbar.GrabFocus()
}
func (ui *UI) back() {
	ui.GetCurrentTab().HistoryBack()
}

func (ui *UI) next() {
	ui.GetCurrentTab().HistoryNext()
}

func (ui *UI) Find() {
	// ui.GetCurrentTab().Find()
}
func (ui *UI) FindNext() {
	// currentTab().FindNext(true)
}
func (ui *UI) FindPrev() {
	// currentTab().FindNext(false)
}
func (ui *UI) ReplaceOne() {
	// currentTab().Replace(false)
}
func (ui *UI) ReplaceAll() {
	// currentTab().Replace(true)
}
func (ui *UI) ToggleMenuBar() {
	// conf.UI.MenuBarVisible = !conf.UI.MenuBarVisible
	// ui.menubar.SetVisible(conf.UI.MenuBarVisible)
}
func (ui *UI) ShowFindbar() {
}
func (ui *UI) ShowReplbar() {
}

func (ui *UI) Quit() {
	gtk.MainQuit()
}
