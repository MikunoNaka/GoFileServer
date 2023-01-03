package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"log"
	"strconv"
	"sync"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

var (
	DEFAULT_PORT, DEFAULT_DIR, APP_VERSION string = "8080", ".", "v0.2.0"
)

func serve(port, dir string, wg *sync.WaitGroup) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(dir)))
    server := &http.Server{Addr: ":" + port, Handler: mux}

	if port == "" { port = DEFAULT_PORT }

    go func() {
		wg.Add(1)
        defer wg.Done()

        if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Error while running the server: %v\n", err)
        }
    }()

    return server
}


func main() {
	app, _ := gtk.ApplicationNew("net.mikunonaka.gofileserver", glib.APPLICATION_FLAGS_NONE)
	app.Connect("activate", func() { onActivate(app) })
	app.Run(os.Args)
}

func onActivate(app *gtk.Application) {
	win, _ := gtk.ApplicationWindowNew(app)
	win.SetTitle("HTTP File Server")
	win.SetDefaultSize(330, 230)
	win.SetResizable(false)
	win.SetPosition(gtk.WIN_POS_CENTER)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

	dirBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	dirLabel, _ := gtk.LabelNew("Directory to serve: ")
	dirBox.PackStart(dirLabel, false, false, 5)
	browseButton, _ := gtk.FileChooserButtonNew("Browse Directory To Serve", gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER)
	dirBox.PackEnd(browseButton, false, false, 5)
	box.PackStart(dirBox, false, false, 5)

	portBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	portLabel, _ := gtk.LabelNew("Port:")
	portBox.PackStart(portLabel, false, false, 5)
	portInput, _ := gtk.EntryNew()
	portInput.SetText(DEFAULT_PORT)
	portBox.PackEnd(portInput, false, false, 5)
	box.PackStart(portBox, false, false, 5)

	buttonBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	aboutButton, _ := gtk.ButtonNewWithLabel("About")
	buttonBox.PackStart(aboutButton, false, true, 5)
	buttonSwitch, _ := gtk.ButtonNewWithLabel("Start")
	buttonBox.PackEnd(buttonSwitch, true, true, 5)
	box.PackStart(buttonBox, false, false, 5)

	statusLabel, _ := gtk.LabelNew("")
	statusLabel.SetLineWrap(true)
	statusLabel.SetLineWrapMode(pango.WrapMode(gtk.ALIGN_START))
	statusLabel.SetJustify(gtk.JUSTIFY_CENTER)
	box.PackStart(statusLabel, false, true, 10)

	var on = false
	var server *http.Server
	buttonSwitch.Connect("clicked", func() {
		// clicking the button too fast too many times crashes the app
		// this disables the button as long as the button is running
		// user likely won't notice the button greying out or something
		buttonSwitch.SetSensitive(false)
		defer buttonSwitch.SetSensitive(true)

		port, _ := portInput.GetText()
		dir := browseButton.GetFilename()
		if dir == "" { dir = DEFAULT_DIR }

		if on {
			server.Shutdown(context.TODO())
			statusLabel.SetMarkup("<span foreground='#ffcc00'>Server was terminated by user.</span>")
		} else {
			go func() {
				killServerDone := &sync.WaitGroup{}
				server = serve(port, dir, killServerDone)

				//killServerDone.Add(1)
				// do this after server starts
				on = true
				buttonSwitch.SetLabel("Stop")
				statusLabel.SetMarkup(fmt.Sprintf("Serving\n%s\nOn <a href=\"http://localhost:%s\">http://localhost:%s</a>", dir, port, port))
				browseButton.SetCanFocus(false)
				portInput.SetEditable(false)
				portInput.SetCanFocus(false)

				killServerDone.Wait()
				// do this after server shuts down
				on = false
				buttonSwitch.SetLabel("Start")
				browseButton.SetCanFocus(true)
				portInput.SetEditable(true)
				portInput.SetCanFocus(true)
			}()
		}
	})

	portInput.Connect("insert-text", func(_ *gtk.Entry, input string) {
		if _, err := strconv.Atoi(input); err != nil {
			// if input is not numeric, don't insert
			portInput.StopEmission("insert-text")
		} else {
			// stop if port number is invalid / out of range
			currentInput, _ := portInput.GetText()
			nextInput := currentInput + input
			n, _ := strconv.Atoi(nextInput)
			if n < 0 || n > 65535 {
				portInput.StopEmission("insert-text")
			}
		}
	})

	browseButton.Connect("selection-changed", func() {
		/* when different directory is selected,
		 * shut down the server and show warning message */
		if on {
			server.Shutdown(context.TODO())
			statusLabel.SetMarkup("<span foreground='#ffcc00'>Server was terminated because root directory changed.</span>")
		}
	})

	aboutButton.Connect("clicked", func() {
		aboutWindow, _ := gtk.ApplicationWindowNew(app)
		aboutWindow.SetTitle("About - GoFileServer")
		aboutWindow.SetDefaultSize(420, 180)
		aboutWindow.SetResizable(false)
		aboutWindow.SetPosition(gtk.WIN_POS_MOUSE)

		box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)

		titleLabel, _ := gtk.LabelNew("")
		titleLabel.SetMarkup(fmt.Sprintf("<b>GoFileServer</b> <span color='lightgray'>%s</span>", APP_VERSION))
		box.PackStart(titleLabel, true, true, 5)

		copyrightLabel, _ := gtk.LabelNew("Copyright (c) 2022 Vidhu Kant Sharma")
		box.PackStart(copyrightLabel, true, true, 5)

		urlLabel, _ := gtk.LabelNew("")
		urlLabel.SetMarkup("<a href='https://github.com/MikunoNaka/GoFileServer'>https://github.com/MikunoNaka/GoFileServer</a>")
		box.PackStart(urlLabel, true, true, 5)

		gplLabel, _ := gtk.LabelNew("")
		gplLabel.SetMarkup("<span>This program comes with absolutely no warranty.\nSee the <a href='https://www.gnu.org/licenses/gpl-3.0.en.html'>GNU General Public License Version 3 or later</a> for details.</span>")
		gplLabel.SetLineWrap(true)
		gplLabel.SetLineWrapMode(pango.WrapMode(gtk.ALIGN_START))
		gplLabel.SetJustify(gtk.JUSTIFY_CENTER)
		box.PackStart(gplLabel, true, true, 5)

		buttonBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
		closeButton, _ := gtk.ButtonNewWithLabel("Close")
		buttonBox.PackEnd(closeButton, false, true, 5)
		box.PackEnd(buttonBox, true, true, 5)

		closeButton.Connect("clicked", func() { aboutWindow.Close() })

		aboutWindow.Add(box)
		aboutWindow.ShowAll()
	})

	win.Add(box)
	win.ShowAll()
}
