package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/atotto/clipboard"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {
	var progress *walk.ProgressBar
	var main *walk.MainWindow

	mainWindow := MainWindow{
		Title:   "Test",
		MinSize: Size{Width: 600, Height: 50},
		Layout:  VBox{},
		Children: []Widget{
			ProgressBar{
				MinValue: 0,
				MaxValue: 100,
				AssignTo: &progress,
			},
		},
		AssignTo: &main,
	}

	go func() {
		for {
			if main != nil || progress != nil {
				main.SetBounds(walk.Rectangle{X: 0, Y: 0, Width: 500, Height: 50})

				args := os.Args[1:]
				if len(args) != 1 {
					walk.MsgBox(main, "Error", "Invalid arguments!", walk.MsgBoxOK)
					os.Exit(1)
				}
				fileName := args[0]
				file, err := os.Open(fileName)
				if err != nil {
					walk.MsgBox(main, "Error", err.Error(), walk.MsgBoxOK)
					os.Exit(1)
				}
				defer file.Close()

				fileStat, err := file.Stat()
				if err != nil {
					walk.MsgBox(main, "Error", err.Error(), walk.MsgBoxOK)
					os.Exit(1)
				}

				pr := &progressReader{
					r:   file,
					max: int(fileStat.Size()),
					report: func(sent int, max int) {
						if sent <= max {
							p := int(100 * (float64(sent) / float64(max)))
							progress.SetValue(p)
						}
					},
				}
				req, err := http.NewRequest("PUT", "https://transfer.sh/"+filepath.Base(fileName), pr)
				if err != nil {
					walk.MsgBox(main, "Error", err.Error(), walk.MsgBoxOK)
					os.Exit(1)
				}

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					walk.MsgBox(main, "Error", err.Error(), walk.MsgBoxOK)
					os.Exit(1)
				}

				defer resp.Body.Close()

				responseData, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					walk.MsgBox(main, "Error", err.Error(), walk.MsgBoxOK)
					os.Exit(1)
				}
				uploadedURL := string(responseData)
				clipboard.WriteAll(uploadedURL)
				os.Exit(1)
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	mainWindow.Run()

}

type progressReader struct {
	r      io.Reader
	max    int
	sent   int
	report func(int, int)
	atEOF  bool
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	pr.sent += n
	if err == io.EOF {
		pr.atEOF = true
	}
	pr.report(pr.sent, pr.max)
	return n, err
}
