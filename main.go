package main

import (
	"os"
	"fmt"
	"net/http"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func serve(port, dir string) {
	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.ListenAndServe(":" + port, nil)
}

func main() {
    app, _ := gtk.ApplicationNew("net.mikunonaka.fileserver", glib.APPLICATION_FLAGS_NONE)
    app.Connect("activate", func() { onActivate(app) })
    app.Run(os.Args)
}

func onActivate(app *gtk.Application) {
	win, _ := gtk.ApplicationWindowNew(app)
    win.SetTitle("HTTP File Server")
    win.SetDefaultSize(400, 400)

	grid, _ := gtk.GridNew()
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)

	dirLabel, _ := gtk.LabelNew("Directory to serve: ")
	grid.Attach(dirLabel, 0, 0, 1, 1)
    dirInput, _ := gtk.EntryNew()
	grid.Attach(dirInput, 1, 0, 2, 1)

	portLabel, _ := gtk.LabelNew("Port:")
	grid.Attach(portLabel, 0, 1, 1, 1)
    portInput, _ := gtk.EntryNew()
	grid.Attach(portInput, 1, 1, 2, 1)

	buttonSwitch, _ := gtk.ButtonNew()
	buttonSwitch.SetLabel("Start")
	grid.Attach(buttonSwitch, 1, 2, 2, 1)

	statusLabel, _ := gtk.LabelNew("")
	grid.Attach(statusLabel, 0, 3, 1, 4)

	buttonSwitch.Connect("clicked", func() {
		port, _ := portInput.GetText()
		dir, _ := dirInput.GetText()
		go serve(port, dir)
		statusText := fmt.Sprintf("Serving directory '%s' on port '%s'", dir, port)
		statusLabel.SetText(statusText)
	})

    win.Add(grid)
    win.ShowAll()
}
