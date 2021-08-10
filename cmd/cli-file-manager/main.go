package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	cfm "github.com/0l1v3rr/cli-file-manager/pkg"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var path string

var l = widgets.NewList()

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

	initWidgets()

}

func initWidgets() {
	l.Title = "CLI File Manager"
	l.Rows = cfm.ReadFiles(path)
	l.TextStyle = ui.NewStyle(ui.ColorWhite)
	l.WrapText = false
	l.SetRect(0, 0, 35, 20)
	l.BorderStyle.Fg = ui.ColorBlue
	l.TitleStyle.Modifier = ui.ModifierBold
	l.SelectedRowStyle.Fg = ui.ColorBlue
	l.SelectedRowStyle.Modifier = ui.ModifierBold

	p := widgets.NewParagraph()
	p.Title = "Help Menu"
	p.Text = "[↑](fg:green) - Scroll Up\n[↓](fg:green) - Scroll Down\n[q](fg:green) - Quit\n[Enter](fg:green) - Open\n"
	p.SetRect(35, 0, 70, 20)
	p.BorderStyle.Fg = ui.ColorBlue
	p.TitleStyle.Modifier = ui.ModifierBold

	p2 := widgets.NewParagraph()
	p2.Title = "File Information"
	p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
	p2.SetRect(0, 30, 70, 20)
	p2.BorderStyle.Fg = ui.ColorBlue
	p2.TitleStyle.Modifier = ui.ModifierBold

	ui.Render(l, p, p2)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Down>":
			l.ScrollDown()
			p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
		case "<Up>":
			l.ScrollUp()
			p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
		case "<Home>":
			l.ScrollTop()
			p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
		case "<End>":
			l.ScrollBottom()
			p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
		case "<Enter>":
			selected := getFileName(l.SelectedRow)
			if selected[len(selected)-1] == '/' {
				if selected == "../" {
					splitted := strings.Split(path, "/")
					if len(splitted) > 0 {
						splitted = splitted[:len(splitted)-1]
					}
					path = strings.Join(splitted, "/")
				} else {
					if path[len(path)-1] == '/' || selected[0] == '/' {
						path = fmt.Sprintf("%v%v", path, selected)
					} else {
						path = fmt.Sprintf("%v/%v", path, selected)
					}
				}
				l.Rows = cfm.ReadFiles(path)

				l.SelectedRow = 0
				l.SelectedRowStyle.Fg = ui.ColorBlue
				l.SelectedRowStyle.Modifier = ui.ModifierBold
				p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
			}
		}

		ui.Render(l, p2)
	}
}

func getFileName(n int) string {
	row := l.Rows[n]
	sliced := strings.Split(strings.Replace(row, "](fg:green)", "", 1), " ")
	sliced = sliced[1:]
	result := strings.Join(sliced, " ")

	return result
}
