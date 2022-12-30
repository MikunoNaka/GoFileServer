package main

import (
	"os"
	"net/http"
	"log"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"context"
	"sync"
)

func serve(port, dir string, wg *sync.WaitGroup) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(dir)))
    server := &http.Server{Addr: ":" + port, Handler: mux}

    go func() {
        defer wg.Done()

        if err := server.ListenAndServe(); err != http.ErrServerClosed {
			// TODO: show error in GTK error window
            log.Fatalf("ListenAndServe(): %v", err)
        }
    }()

    return server
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

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

	dirBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	dirLabel, _ := gtk.LabelNew("Directory to serve: ")
	dirBox.PackStart(dirLabel, false, false, 5)
	dirInput, _ := gtk.EntryNew()
	dirBox.PackEnd(dirInput, false, false, 5)
	box.PackStart(dirBox, false, false, 5)

	portBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	portLabel, _ := gtk.LabelNew("Port:")
	portBox.PackStart(portLabel, false, false, 5)
	portInput, _ := gtk.EntryNew()
	portBox.PackEnd(portInput, false, false, 5)
	box.PackStart(portBox, false, false, 5)

	buttonBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	buttonSwitch, _ := gtk.ButtonNew()
	buttonSwitch.SetLabel("Start")
	buttonBox.PackStart(buttonSwitch, true, true, 0)
	box.PackStart(buttonBox, false, false, 5)

	statusLabel, _ := gtk.LabelNew("")
	box.PackStart(statusLabel, false, true, 10)

	var on = false
	var server *http.Server
	buttonSwitch.Connect("clicked", func() {
		port, _ := portInput.GetText()
		dir, _ := dirInput.GetText()

		if on {
			server.Shutdown(context.TODO())
		} else {
			go func() {
				killServerDone := &sync.WaitGroup{}
				server = serve(port, dir, killServerDone)
				killServerDone.Add(1)
				on = true
				buttonSwitch.SetLabel("Stop")
				killServerDone.Wait()
				on = false
				buttonSwitch.SetLabel("Start")
			}()
		}
	})

	win.Add(box)
	win.ShowAll()
}
