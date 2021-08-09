package main

import (
	"flag"
	"log"
	"os"

	cfm "github.com/0l1v3rr/cli-file-manager/pkg"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var path string

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	defaultPath, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	flag.StringVar(&path, "path", defaultPath, "The path of the folder.")
	flag.Parse()

	initWidgets(path)

}

func initWidgets(path string) {
	l := widgets.NewList()
	l.Title = path
	l.Rows = cfm.ReadFiles(path)
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 0, 40, 20)

	ui.Render(l)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Down>":
			l.ScrollDown()
		case "<Up>":
			l.ScrollUp()
		case "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		}

		ui.Render(l)
	}
}
