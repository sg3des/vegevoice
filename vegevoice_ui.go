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

	findbar *findbar
}

func CreateUi() *UserInterface {
	ui := &UserInterface{}
	ui.window = gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	ui.window.SetSizeRequest(900, 600)
	ui.window.SetTitle("webkit")
	ui.window.Connect("destroy", ui.Quit)
	// ui.window.Connect("check-resize", ui.windowResize)
	ui.findbar = ui.createFindbar()

	ui.menubar = ui.createMenubar()
	ui.notebook = gtk.NewNotebook()
	ui.notebook.SetBorderWidth(0)
	ui.notebook.SetShowBorder(true)
	ui.notebook.SetTabBorder(1)

	ui.vbox = gtk.NewVBox(false, 0)
	ui.vbox.PackStart(ui.menubar, true, true, 0)
	ui.vbox.PackStart(ui.notebook, true, true, 0)
	ui.vbox.PackStart(ui.findbar.table, false, false, 0)

	ui.window.Add(ui.vbox)
	ui.window.ShowAll()

	ui.menubar.SetVisible(false)
	ui.findbar.SetVisible(false)

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

	ui.newActionStock("Find", gtk.STOCK_FIND, "", ui.showFindbar)
	ui.newAction("FindNext", "Find Next", "F3", ui.findbar.FindNext)
	ui.newAction("FindPrev", "Find Previous", "<shift>F3", ui.findbar.FindPrev)

	// View
	ui.actionGroup.AddAction(gtk.NewAction("View", "View", "", ""))
	// ui.actionGroup.AddAction(gtk.NewAction("Encoding", "Encoding", "", ""))

	ui.newToggleAction("Menubar", "Menubar", "<control>M", false, ui.toggleMenuBar)

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

func (ui *UserInterface) newTab() {
	ui.NewTab("")
}

func (ui *UserInterface) reload() {
	ui.GetCurrentTab().Reload()
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

func (ui *UserInterface) focusurl() {
	ui.GetCurrentTab().urlbar.GrabFocus()
}
func (ui *UserInterface) back() {
	ui.GetCurrentTab().HistoryBack()
}

func (ui *UserInterface) next() {
	ui.GetCurrentTab().HistoryNext()
}

func (ui *UserInterface) showFindbar() {
	ui.findbar.SetVisible(true)
}

func (ui *UserInterface) toggleMenuBar() {
	// conf.UserInterface.MenuBarVisible = !conf.UserInterface.MenuBarVisible
	// ui.menubar.SetVisible(conf.UserInterface.MenuBarVisible)
}

func (ui *UserInterface) Quit() {
	gtk.MainQuit()
}

type findbar struct {
	table     *gtk.Table
	entryFind *gtk.Entry
	btnCase   *gtk.ToggleButton
	btnNext   *gtk.Button
	btnPrev   *gtk.Button
}

func (ui *UserInterface) createFindbar() *findbar {
	var fb = new(findbar)

	//findbar
	fb.entryFind = gtk.NewEntry()
	fb.entryFind.Connect("changed", fb.FindNext)

	fb.btnCase = gtk.NewToggleButtonWithLabel("Aa")
	fb.btnCase.Clicked(fb.FindNext)

	fb.btnNext = gtk.NewButton()
	fb.btnNext.SetImage(gtk.NewArrow(gtk.ARROW_RIGHT, gtk.SHADOW_NONE))
	fb.btnNext.Clicked(fb.FindNext)

	fb.btnPrev = gtk.NewButton()
	fb.btnPrev.SetImage(gtk.NewArrow(gtk.ARROW_LEFT, gtk.SHADOW_NONE))
	fb.btnPrev.Clicked(fb.FindPrev)

	btnClose := gtk.NewButton()
	btnClose.SetImage(gtk.NewImageFromStock(gtk.STOCK_CLOSE, gtk.ICON_SIZE_MENU))
	btnClose.Clicked(func() { fb.SetVisible(false) })

	fb.table = gtk.NewTable(1, 5, false)
	fb.table.Attach(fb.btnCase, 0, 1, 0, 1, gtk.FILL, gtk.FILL, 0, 0)
	fb.table.Attach(fb.entryFind, 1, 2, 0, 1, gtk.EXPAND|gtk.FILL, gtk.FILL, 0, 0)
	fb.table.Attach(fb.btnPrev, 2, 3, 0, 1, gtk.FILL, gtk.FILL, 0, 0)
	fb.table.Attach(fb.btnNext, 3, 4, 0, 1, gtk.FILL, gtk.FILL, 0, 0)
	fb.table.Attach(btnClose, 4, 5, 0, 1, gtk.FILL, gtk.FILL, 0, 0)
	fb.table.ShowAll()

	return fb
}

func (fb *findbar) SetVisible(b bool) {
	fb.table.SetVisible(b)
	if b {
		fb.entryFind.GrabFocus()
	}
}

func (fb *findbar) FindNext() {
	text := fb.entryFind.GetText()
	if len(text) == 0 {
		return
	}
	ui.GetCurrentTab().Find(text, fb.btnCase.GetActive(), true)
}

func (fb *findbar) FindPrev() {
	text := fb.entryFind.GetText()
	if len(text) == 0 {
		return
	}
	ui.GetCurrentTab().Find(text, fb.btnCase.GetActive(), false)
}
