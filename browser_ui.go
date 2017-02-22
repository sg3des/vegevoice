package main

import (
	"github.com/mattn/go-gtk/gtk"
)

var (
	width  int
	height int
)

type UserInterface struct {
	window *gtk.Window
	vbox   *gtk.VBox

	accelGroup  *gtk.AccelGroup
	actionGroup *gtk.ActionGroup

	menubar  *gtk.Widget
	notebook *gtk.Notebook
	tabs     []*Tab
}

func CreateUi() *UserInterface {
	ui := &UserInterface{}
	ui.window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	ui.window.SetSizeRequest(900, 600)
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

func (ui *UserInterface) createMenubar() *gtk.Widget {

	UserInterfacexml := `
<ui>
	<menubar name='MenuBar'>
		<menu action='File'>
			<menuitem action='NewTab' />
			<menuitem action='CloseTab' />
			<menuitem action='OpenUrl' />
			<menuitem action='Reload' />
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
	uiManager.AddUIFromString(UserInterfacexml)

	ui.accelGroup = uiManager.GetAccelGroup()
	ui.window.AddAccelGroup(ui.accelGroup)

	ui.actionGroup = gtk.NewActionGroup("my_group")
	uiManager.InsertActionGroup(ui.actionGroup, 0)

	// File
	ui.actionGroup.AddAction(gtk.NewAction("File", "File", "", ""))

	ui.newAction("NewTab", "New Tab", "<control>t", ui.newTab)
	ui.newAction("CloseTab", "Close Tab", "<control>w", ui.CloseCurrentTab)
	ui.newAction("OpenUrl", "Open URL", "<control>l", ui.focusurl)
	ui.newAction("Reload", "Reload", "<control>r", ui.reload)
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

func (ui *UserInterface) newAction(dst, label, accel string, f func()) {
	action := gtk.NewAction(dst, label, "", "")
	action.Connect("activate", f)
	ui.actionGroup.AddActionWithAccel(action, accel)
}

func (ui *UserInterface) newActionStock(dst, stock, accel string, f func()) {
	action := gtk.NewAction(dst, "", "", stock)
	action.Connect("activate", f)
	ui.actionGroup.AddActionWithAccel(action, accel)
}

func (ui *UserInterface) newToggleAction(dst, label, accel string, state bool, f func()) {
	action := gtk.NewToggleAction(dst, label, "", "")
	action.SetActive(state)
	action.Connect("activate", f)
	ui.actionGroup.AddActionWithAccel(&action.Action, accel)
}

// actions
func (ui *UserInterface) windowResize() {
	ui.window.GetSize(&width, &height)
	ui.notebook.SetSizeRequest(width, height)
	ui.homogenousTabs()
}

func (ui *UserInterface) homogenousTabs() {
	lenTabs := len(ui.tabs)
	if lenTabs == 0 {
		return
	}

	// var pinned int
	for _, t := range ui.tabs {
		if t.Pinned {
			lenTabs--
		}
	}

	if lenTabs == 0 {
		return
	}

	tabwidth := (width - lenTabs*6) / lenTabs
	leftwidth := (width - lenTabs*6) % lenTabs

	for _, t := range ui.tabs {
		if t.Pinned {
			continue
		}

		if leftwidth > 0 {
			t.tabbox.SetSizeRequest(tabwidth+1, 12)
			leftwidth--
		} else {
			t.tabbox.SetSizeRequest(tabwidth, 12)
		}
	}
}

func (ui *UserInterface) newTab() {
	ui.NewTab("")
}

func (ui *UserInterface) reload() {
	ui.GetCurrentTab().Reload()
}

// func (ui *UserInterface) closeTab() {
// 	ui.CloseCurrentTab()

// }
func (ui *UserInterface) focusurl() {
	ui.GetCurrentTab().urlbar.GrabFocus()
}
func (ui *UserInterface) back() {
	ui.GetCurrentTab().HistoryBack()
}

func (ui *UserInterface) next() {
	ui.GetCurrentTab().HistoryNext()
}

func (ui *UserInterface) Find() {
	// ui.GetCurrentTab().Find()
}
func (ui *UserInterface) FindNext() {
	// currentTab().FindNext(true)
}
func (ui *UserInterface) FindPrev() {
	// currentTab().FindNext(false)
}
func (ui *UserInterface) ReplaceOne() {
	// currentTab().Replace(false)
}
func (ui *UserInterface) ReplaceAll() {
	// currentTab().Replace(true)
}
func (ui *UserInterface) ToggleMenuBar() {
	// conf.UserInterface.MenuBarVisible = !conf.UserInterface.MenuBarVisible
	// ui.menubar.SetVisible(conf.UserInterface.MenuBarVisible)
}
func (ui *UserInterface) ShowFindbar() {
}
func (ui *UserInterface) ShowReplbar() {
}

func (ui *UserInterface) Quit() {
	gtk.MainQuit()
}
