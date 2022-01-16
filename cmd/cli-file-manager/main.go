package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	cfm "github.com/0l1v3rr/cli-file-manager/pkg"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/skratchdot/open-golang/open"
)

const VERSION = "v1.1"

var (
	path       string
	l               = widgets.NewList()
	showHidden bool = true
)

func main() {

	defaultPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("%v", err)
	}

	if len(os.Args) > 1 && os.Args[1] != "" {
		path = os.Args[1]
	} else {
		path = defaultPath
	}

	err2 := ui.Init()
	if err2 != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	initWidgets()
}

func initWidgets() {
	l.Title = fmt.Sprintf("CLI File Manager - %s", VERSION)
	l.Rows = cfm.ReadFiles(path, showHidden)
	l.TextStyle = ui.NewStyle(ui.ColorWhite)
	l.WrapText = false
	l.SetRect(0, 0, cfm.GetCliWidth()/2, int(float64(cfm.GetCliHeight())*0.73))
	l.BorderStyle.Fg = ui.ColorBlue
	l.TitleStyle.Modifier = ui.ModifierBold
	l.SelectedRowStyle.Fg = ui.ColorBlue
	l.SelectedRowStyle.Modifier = ui.ModifierBold

	p := widgets.NewParagraph()
	p.Title = "Help Menu"
	pText := "[↑](fg:green) - Scroll Up\n[↓](fg:green) - Scroll Down\n[q](fg:green) - Quit\n[Enter](fg:green) - Open\n[m](fg:green) - Memory Usage\n[f](fg:green) - Disk Information\n[^D (2 times)](fg:green) - Remove file\n[^F](fg:green) - Create file\n[^N](fg:green) - Create folder\n[^R](fg:green) - Rename file\n[^V](fg:green) - Launch VS Code\n[C](fg:green) - Copy file\n[h](fg:green) - Hide hidden files"
	p.Text = pText
	p.SetRect(cfm.GetCliWidth()/2, 0, cfm.GetCliWidth(), int(float64(cfm.GetCliHeight())*0.58))
	p.BorderStyle.Fg = ui.ColorBlue
	p.TitleStyle.Modifier = ui.ModifierBold

	disk := cfm.DiskUsage("/")

	p3 := widgets.NewParagraph()
	json, err := cfm.ReadJson()
	if json == "memory" || err != nil {
		p3.Title = "Memory Usage"
		p3.Text = cfm.ReadMemStats()
	} else {
		p3.Title = "Disk Information"
		p3.Text = fmt.Sprintf("[All: ](fg:green) - %.2f GB\n[Used:](fg:green) - %.2f GB\n[Free:](fg:green) - %.2f GB", float64(disk.All)/float64(1024*1024*1024), float64(disk.Used)/float64(1024*1024*1024), float64(disk.Free)/float64(1024*1024*1024))
	}
	p3.SetRect(cfm.GetCliWidth()/2, int(float64(cfm.GetCliHeight())*0.73), cfm.GetCliWidth(), int(float64(cfm.GetCliHeight())*0.58))
	p3.BorderStyle.Fg = ui.ColorBlue
	p3.TitleStyle.Modifier = ui.ModifierBold

	p2 := widgets.NewParagraph()
	p2.Title = "File Information"
	p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
	p2.SetRect(0, cfm.GetCliHeight(), cfm.GetCliWidth(), int(float64(cfm.GetCliHeight())*0.73))
	p2.BorderStyle.Fg = ui.ColorBlue
	p2.WrapText = false
	p2.TitleStyle.Modifier = ui.ModifierBold

	ui.Render(l, p, p2, p3)

	copyPath := ""
	previousKey := ""
	inputField := ""
	originalName := ""
	fileCreatingInProgress := false
	dirCreatingInProgress := false
	renameInProgress := false
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Down>":
			if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress {
				l.ScrollDown()
				p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
			}
		case "<Up>":
			if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress {
				l.ScrollUp()
				p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
			}
		case "<Home>":
			if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress {
				l.ScrollTop()
				p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
			}
		case "<End>":
			if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress {
				l.ScrollBottom()
				p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
			}
		case "<C-d>":
			if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress {
				if previousKey == "<C-d>" {
					selected := getFileName(l.SelectedRow)
					if selected != ".." && selected != "../" {
						filePath := ""
						if path[len(path)-1] == '/' || selected[0] == '/' {
							filePath = fmt.Sprintf("%v%v", path, selected)
						} else {
							filePath = fmt.Sprintf("%v/%v", path, selected)
						}
						err := os.Remove(filePath)
						if err == nil {
							l.Rows = cfm.ReadFiles(path, showHidden)
							l.SelectedRow = 0
							p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
						} else {
							err2 := os.RemoveAll(filePath)
							if err2 == nil {
								l.Rows = cfm.ReadFiles(path, showHidden)
								l.SelectedRow = 0
								p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
							}
						}
					}
				}
			}
		case "<C-f>":
			if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress {
				fileCreatingInProgress = true
				l.Rows = append(l.Rows, fmt.Sprintf("[?]: %v", inputField))
				l.SelectedRow = len(l.Rows) - 1
				textFieldStyle()
				p2.Text = cfm.EmptyFileInfo()
				p.Text = "[Esc](fg:green) - Cancel\n[Enter](fg:green) - Apply Changes\n"
			}
		case "<C-n>":
			if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress {
				dirCreatingInProgress = true
				l.Rows = append(l.Rows, fmt.Sprintf("[$]: %v", inputField))
				l.SelectedRow = len(l.Rows) - 1
				textFieldStyle()
				p2.Text = cfm.EmptyFileInfo()
				p.Text = "[Esc](fg:green) - Cancel\n[Enter](fg:green) - Apply Changes\n"
			}
		case "<C-r>":
			if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress {
				renameInProgress = true
				originalName = l.Rows[l.SelectedRow]
				inputField = getFileName(l.SelectedRow)
				if inputField[len(inputField)-1] == '/' {
					inputField = inputField[:len(inputField)-1]
				}
				l.Rows[l.SelectedRow] = fmt.Sprintf("[#]: %v", inputField)
				textFieldStyle()
				p.Text = "[Esc](fg:green) - Cancel\n[Enter](fg:green) - Apply Changes\n"
			}
		case "<Escape>":
			resetColors()
			if copyPath != "" {
				copyPath = ""
				p.Text = pText
			} else {
				if fileCreatingInProgress {
					fileCreatingInProgress = false
					inputField = ""
					l.SelectedRow = 0
					l.Rows = l.Rows[:len(l.Rows)-1]
					p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
					p.Text = pText
				} else if dirCreatingInProgress {
					dirCreatingInProgress = false
					inputField = ""
					l.SelectedRow = 0
					l.Rows = l.Rows[:len(l.Rows)-1]
					p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
					p.Text = pText
				} else if renameInProgress {
					renameInProgress = false
					inputField = ""
					l.Rows[l.SelectedRow] = originalName
					originalName = ""
					p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
					p.Text = pText
				}
			}
		case "m":
			if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress {
				p3.Title = "Memory Usage"
				p3.Text = cfm.ReadMemStats()
			}
		case "f":
			if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress {
				p3.Title = "Disk Information"
				p3.Text = fmt.Sprintf("[All: ](fg:green) - %.2f GB\n[Used:](fg:green) - %.2f GB\n[Free:](fg:green) - %.2f GB", float64(disk.All)/float64(1024*1024*1024), float64(disk.Used)/float64(1024*1024*1024), float64(disk.Free)/float64(1024*1024*1024))
			}
		case "c":
			if getFileName(l.SelectedRow)[len(getFileName(l.SelectedRow))-1] != '/' {
				copyPath = fmt.Sprintf("%s/%s", path, getFileName(l.SelectedRow))
				p.Text = fmt.Sprintf("[↑](fg:green) - Scroll Up\n[↓](fg:green) - Scroll Down\n[q](fg:green) - Quit\n[Enter](fg:green) - Open\n[m](fg:green) - Memory Usage\n[f](fg:green) - Disk Information\n[^D (2 times)](fg:green) - Remove file\n[^F](fg:green) - Create file\n[^N](fg:green) - Create folder\n[^R](fg:green) - Rename file\n[^V](fg:green) - Launch VS Code\n[C](fg:cyan) - Copied to clipboard ([%s](fg:cyan))\n[V](fg:green) - Paste", getFileName(l.SelectedRow))
			}
		case "v":
			if copyPath != "" {
				p.Text = "[↑](fg:green) - Scroll Up\n[↓](fg:green) - Scroll Down\n[q](fg:green) - Quit\n[Enter](fg:green) - Open\n[m](fg:green) - Memory Usage\n[f](fg:green) - Disk Information\n[^D (2 times)](fg:green) - Remove file\n[^F](fg:green) - Create file\n[^N](fg:green) - Create folder\n[^R](fg:green) - Rename file\n[^V](fg:green) - Launch VS Code\n[C](fg:cyan) - Copying..."
				cfm.Copy(copyPath, path)
				p.Text = pText
				copyPath = ""
				l.Rows = cfm.ReadFiles(path, showHidden)
				ui.Render(l, p, p2, p3)
			} else {
				p.Text = pText
			}
		case "<Resize>":
			l.SetRect(0, 0, cfm.GetCliWidth()/2, int(float64(cfm.GetCliHeight())*0.73))
			p.SetRect(cfm.GetCliWidth()/2, 0, cfm.GetCliWidth(), int(float64(cfm.GetCliHeight())*0.58))
			p2.SetRect(0, cfm.GetCliHeight(), cfm.GetCliWidth(), int(float64(cfm.GetCliHeight())*0.73))
			p3.SetRect(cfm.GetCliWidth()/2, int(float64(cfm.GetCliHeight())*0.73), cfm.GetCliWidth(), int(float64(cfm.GetCliHeight())*0.58))
			ui.Render(l, p, p2, p3)
		case "<C-v>":
			cmd := exec.Command("code", path)
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		case "h":
			if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress && copyPath == "" {
				if showHidden {
					p.Text = "[↑](fg:green) - Scroll Up\n[↓](fg:green) - Scroll Down\n[q](fg:green) - Quit\n[Enter](fg:green) - Open\n[m](fg:green) - Memory Usage\n[f](fg:green) - Disk Information\n[^D (2 times)](fg:green) - Remove file\n[^F](fg:green) - Create file\n[^N](fg:green) - Create folder\n[^R](fg:green) - Rename file\n[^V](fg:green) - Launch VS Code\n[C](fg:green) - Copy file\n[h](fg:green) - Show hidden files"
				} else {
					p.Text = pText
				}
				showHidden = !showHidden
				l.Rows = cfm.ReadFiles(path, showHidden)
			}
		case "<Enter>":
			if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress {
				selected := getFileName(l.SelectedRow)
				if selected[len(selected)-1] == '/' {
					if selected == "../" {
						splitted := strings.Split(path, "/")
						if len(splitted) > 0 {
							if len(splitted) == 2 {
								path = "/"
							} else {
								path = strings.Join(splitted[:len(splitted)-2], "/")
							}
						} else {
							path = "/"
						}
					} else {
						if path[len(path)-1] == '/' || selected[0] == '/' {
							path = fmt.Sprintf("%v%v", path, selected)
						} else {
							path = fmt.Sprintf("%v/%v", path, selected)
						}
					}
					l.Rows = cfm.ReadFiles(path, showHidden)

					l.SelectedRow = 0
					l.SelectedRowStyle.Fg = ui.ColorBlue
					l.SelectedRowStyle.Modifier = ui.ModifierBold
					p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
				} else {
					var filePath string
					if path[len(path)-1] == '/' || selected[0] == '/' {
						filePath = fmt.Sprintf("%v%v", path, selected)
					} else {
						filePath = fmt.Sprintf("%v/%v", path, selected)
					}
					open.Start(filePath)
				}
			} else if fileCreatingInProgress {
				if len(inputField) >= 3 {
					err := ioutil.WriteFile(fmt.Sprintf("%v/%v", path, inputField), []byte(""), 0755)
					if err == nil {
						l.Rows = cfm.ReadFiles(path, showHidden)
						l.SelectedRow = 0
						inputField = ""
						fileCreatingInProgress = false
						resetColors()
						p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
						p.Text = pText
					}
				}
			} else if dirCreatingInProgress {
				if len(inputField) >= 3 {
					err := os.Mkdir(fmt.Sprintf("%v/%v", path, inputField), 0755)
					if err == nil {
						l.Rows = cfm.ReadFiles(path, showHidden)
						l.SelectedRow = 0
						inputField = ""
						dirCreatingInProgress = false
						resetColors()
						p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
						p.Text = pText
					}
				}
			} else if renameInProgress {
				if len(inputField) >= 3 {
					original := getFileNameByFullName(originalName)
					err := os.Rename(fmt.Sprintf("%v/%v", path, original), fmt.Sprintf("%v/%v", path, inputField))
					if err == nil {
						l.Rows = cfm.ReadFiles(path, showHidden)
						inputField = ""
						originalName = ""
						original = ""
						renameInProgress = false
						resetColors()
						p2.Text = cfm.GetFileInformations(fmt.Sprintf("%v/%v", path, getFileName(l.SelectedRow)))
						p.Text = pText
					}
				}
			}
		}

		if fileCreatingInProgress {
			if e.ID[0] != '<' {
				inputField = inputField + e.ID
				l.Rows[len(l.Rows)-1] = fmt.Sprintf("[?]: %v", inputField)
			} else if e.ID == "<Backspace>" {
				le := len(inputField)
				if le > 0 {
					inputField = inputField[:le-1]
				}
				l.Rows[len(l.Rows)-1] = fmt.Sprintf("[?]: %v", inputField)
			}
		} else if dirCreatingInProgress {
			if e.ID[0] != '<' {
				inputField = inputField + e.ID
				l.Rows[len(l.Rows)-1] = fmt.Sprintf("[$]: %v", inputField)
			} else if e.ID == "<Backspace>" {
				le := len(inputField)
				if le > 0 {
					inputField = inputField[:le-1]
				}
				l.Rows[len(l.Rows)-1] = fmt.Sprintf("[$]: %v", inputField)
			}
		} else if renameInProgress {
			if e.ID[0] != '<' {
				inputField = inputField + e.ID
				l.Rows[l.SelectedRow] = fmt.Sprintf("[#]: %v", inputField)
			} else if e.ID == "<Backspace>" {
				le := len(inputField)
				if le > 0 {
					inputField = inputField[:le-1]
				}
				l.Rows[l.SelectedRow] = fmt.Sprintf("[#]: %v", inputField)
			}
		}

		if !fileCreatingInProgress && !dirCreatingInProgress && !renameInProgress {
			if previousKey == "<C-d>" {
				previousKey = ""
			} else {
				previousKey = e.ID
			}
		}

		ui.Render(l, p, p2, p3)
	}
}

func getFileName(n int) string {
	row := l.Rows[n]
	sliced := strings.Split(strings.Replace(row, "](fg:green)", "", 1), " ")
	sliced = sliced[1:]
	result := strings.Join(sliced, " ")

	return result
}

func getFileNameByFullName(s string) string {
	sliced := strings.Split(strings.Replace(s, "](fg:green)", "", 1), " ")
	sliced = sliced[1:]
	result := strings.Join(sliced, " ")

	return result
}

func textFieldStyle() {
	l.SelectedRowStyle.Bg = ui.ColorWhite
	l.SelectedRowStyle.Fg = ui.ColorBlack
}

func resetColors() {
	l.SelectedRowStyle.Bg = ui.ColorClear
	l.SelectedRowStyle.Fg = ui.ColorBlue
}
