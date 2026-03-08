package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/webfraggle/mbd-cli/internal/api"
	"github.com/webfraggle/mbd-cli/internal/config"
)

var commandList = []string{"--next", "--prev", "--setTime", "--setTrain1", "--setTrain2", "--setTrain3", "--setAllTrains", "--image"}

// splitCommandLine splits a command string into tokens, respecting double-quoted arguments.
func splitCommandLine(s string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == '"':
			inQuote = !inQuote
		case c == ' ' && !inQuote:
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteByte(c)
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}

func Run() {
	a := app.New()
	w := a.NewWindow("mbd-cli Konfiguration")
	w.Resize(fyne.NewSize(960, 540))

	// ── Exe path variants ─────────────────────────────────────────────────────
	exeFull, err := os.Executable()
	if err != nil {
		exeFull = "mbd-cli"
	}
	exeBase := filepath.Base(exeFull)
	var exeRel string
	switch runtime.GOOS {
	case "windows":
		exeRel = `.\` + exeBase
	default:
		exeRel = "./" + exeBase
	}

	// ── Shared state ─────────────────────────────────────────────────────────
	configs, _ := config.List()
	var selectedConfig string

	// ── Right panel widgets (declared early, needed by loadConfig) ────────────
	cmdSelect := widget.NewSelect(commandList, nil)
	cmdSelect.SetSelected("--next")

	gleisGroup := widget.NewRadioGroup([]string{"A", "B"}, nil)
	gleisGroup.SetSelected("A")
	gleisGroup.Horizontal = true

	timeEntry := widget.NewEntry()
	timeEntry.SetPlaceHolder("HH:MM")
	timeEntry.SetText("12:30")

	newTrainEntries := func(nr, tim, dest, via string) (eNr, eTim, eDest, eVia, eDelay, eInfo *widget.Entry) {
		eNr = widget.NewEntry(); eNr.SetText(nr)
		eTim = widget.NewEntry(); eTim.SetText(tim)
		eDest = widget.NewEntry(); eDest.SetText(dest)
		eVia = widget.NewEntry(); eVia.SetText(via)
		eDelay = widget.NewEntry(); eDelay.SetText("0")
		eInfo = widget.NewEntry(); eInfo.SetPlaceHolder("Sonderinfo")
		return
	}
	train1Nr, train1Time, train1Dest, train1Via, train1Delay, train1Info := newTrainEntries("ICE123", "12:30", "Berlin", "Hannover")
	train2Nr, train2Time, train2Dest, train2Via, train2Delay, train2Info := newTrainEntries("RE50", "14:00", "Frankfurt", "Kassel")
	train3Nr, train3Time, train3Dest, train3Via, train3Delay, train3Info := newTrainEntries("IC2", "16:30", "Hamburg", "Hannover")

	imageFile := widget.NewEntry()
	imageFile.SetPlaceHolder("00logo.png")
	imageFile.SetText("00logo.png")

	cmdOutput := widget.NewEntry()
	fullPathCheck := widget.NewCheck("Vollständiger Pfad", nil)

	execStatus := widget.NewLabel("")
	dynamicArea := container.NewVBox()

	// ── Build command string ──────────────────────────────────────────────────
	buildCommand := func() string {
		var parts []string
		if fullPathCheck.Checked {
			parts = append(parts, exeFull)
		} else {
			parts = append(parts, exeRel)
		}

		switch cmdSelect.Selected {
		case "--next":
			parts = append(parts, "--next")
		case "--prev":
			parts = append(parts, "--prev")
		case "--setTime":
			parts = append(parts, "--setTime", fmt.Sprintf(`"%s"`, timeEntry.Text))
		case "--setTrain1", "--setTrain2", "--setTrain3":
			n := string(cmdSelect.Selected[len(cmdSelect.Selected)-1])
			var nr, tim, dest, via, delay, info string
			switch n {
			case "1":
				nr, tim, dest, via, delay, info = train1Nr.Text, train1Time.Text, train1Dest.Text, train1Via.Text, train1Delay.Text, train1Info.Text
			case "2":
				nr, tim, dest, via, delay, info = train2Nr.Text, train2Time.Text, train2Dest.Text, train2Via.Text, train2Delay.Text, train2Info.Text
			case "3":
				nr, tim, dest, via, delay, info = train3Nr.Text, train3Time.Text, train3Dest.Text, train3Via.Text, train3Delay.Text, train3Info.Text
			}
			parts = append(parts, "--setTrain"+n, fmt.Sprintf(`"%s"`, strings.Join([]string{nr, tim, dest, via, delay, info}, "|")))
		case "--setAllTrains":
			for i, f := range [3][6]string{
				{train1Nr.Text, train1Time.Text, train1Dest.Text, train1Via.Text, train1Delay.Text, train1Info.Text},
				{train2Nr.Text, train2Time.Text, train2Dest.Text, train2Via.Text, train2Delay.Text, train2Info.Text},
				{train3Nr.Text, train3Time.Text, train3Dest.Text, train3Via.Text, train3Delay.Text, train3Info.Text},
			} {
				parts = append(parts, fmt.Sprintf("--setTrain%d", i+1), fmt.Sprintf(`"%s"`, strings.Join(f[:], "|")))
			}
		case "--image":
			parts = append(parts, "--image", imageFile.Text)
		}

		if gleisGroup.Selected == "B" {
			parts = append(parts, "--gleis B")
		}
		if selectedConfig != "" && selectedConfig != "default" {
			parts = append(parts, "--conf", selectedConfig)
		}
		return strings.Join(parts, " ")
	}

	refreshCmd := func() {
		cmdOutput.SetText(buildCommand())
	}

	// ── Update dynamic fields based on selected command ───────────────────────
	trainForm := func(nr, tim, dest, via, delay, info *widget.Entry) fyne.CanvasObject {
		return widget.NewForm(
			widget.NewFormItem("Zug-Nr", nr),
			widget.NewFormItem("Zeit", tim),
			widget.NewFormItem("Ziel", dest),
			widget.NewFormItem("Via", via),
			widget.NewFormItem("Verspätung", delay),
			widget.NewFormItem("Hinweis", info),
		)
	}
	clearEntries := func(entries ...*widget.Entry) {
		for _, e := range entries {
			e.SetText("")
		}
	}

	updateDynamic := func() {
		var items []fyne.CanvasObject
		switch cmdSelect.Selected {
		case "--setTime":
			items = []fyne.CanvasObject{
				widget.NewForm(widget.NewFormItem("Zeit", timeEntry)),
			}
		case "--setTrain1":
			items = []fyne.CanvasObject{
				trainForm(train1Nr, train1Time, train1Dest, train1Via, train1Delay, train1Info),
				widget.NewButton("Felder leeren", func() {
					clearEntries(train1Nr, train1Time, train1Dest, train1Via, train1Delay, train1Info)
				}),
			}
		case "--setTrain2":
			items = []fyne.CanvasObject{
				trainForm(train2Nr, train2Time, train2Dest, train2Via, train2Delay, train2Info),
				widget.NewButton("Felder leeren", func() {
					clearEntries(train2Nr, train2Time, train2Dest, train2Via, train2Delay, train2Info)
				}),
			}
		case "--setTrain3":
			items = []fyne.CanvasObject{
				trainForm(train3Nr, train3Time, train3Dest, train3Via, train3Delay, train3Info),
				widget.NewButton("Felder leeren", func() {
					clearEntries(train3Nr, train3Time, train3Dest, train3Via, train3Delay, train3Info)
				}),
			}
		case "--setAllTrains":
			makeCol := func(title string, nr, tim, dest, via, delay, info *widget.Entry) fyne.CanvasObject {
				return container.NewVBox(
					widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
					trainForm(nr, tim, dest, via, delay, info),
				)
			}
			items = []fyne.CanvasObject{
				container.NewGridWithColumns(3,
					makeCol("Zug 1", train1Nr, train1Time, train1Dest, train1Via, train1Delay, train1Info),
					makeCol("Zug 2", train2Nr, train2Time, train2Dest, train2Via, train2Delay, train2Info),
					makeCol("Zug 3", train3Nr, train3Time, train3Dest, train3Via, train3Delay, train3Info),
				),
				widget.NewButton("Alle Felder leeren", func() {
					clearEntries(
						train1Nr, train1Time, train1Dest, train1Via, train1Delay, train1Info,
						train2Nr, train2Time, train2Dest, train2Via, train2Delay, train2Info,
						train3Nr, train3Time, train3Dest, train3Via, train3Delay, train3Info,
					)
				}),
			}
		case "--image":
			items = []fyne.CanvasObject{
				widget.NewForm(widget.NewFormItem("Dateiname", imageFile)),
			}
		}
		dynamicArea.Objects = items
		dynamicArea.Refresh()
		refreshCmd()
	}

	// ── Left panel: config list ───────────────────────────────────────────────
	configStatusLabel := widget.NewLabel("")
	endpointEntry := widget.NewEntry()
	endpointEntry.SetPlaceHolder("http://192.168.178.xxx")

	loadConfig := func(name string) {
		selectedConfig = name
		cfg, err := config.Load(name)
		if err == nil {
			endpointEntry.SetText(cfg.Endpoint)
		}
		configStatusLabel.SetText("")
		refreshCmd()
	}

	configList := widget.NewList(
		func() int { return len(configs) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(configs[i])
		},
	)

	refreshList := func() {
		configs, _ = config.List()
		configList.Refresh()
	}

	configList.OnSelected = func(id widget.ListItemID) {
		if id < len(configs) {
			loadConfig(configs[id])
		}
	}

	saveBtn := widget.NewButton("Speichern", func() {
		if err := config.Save(selectedConfig, &config.Config{Endpoint: endpointEntry.Text}); err != nil {
			configStatusLabel.SetText("Fehler: " + err.Error())
		} else {
			configStatusLabel.SetText("Gespeichert.")
		}
	})

	testConnBtn := widget.NewButton("Verbindung testen", func() {
		configStatusLabel.SetText("Teste...")
		go func() {
			client := api.NewClient(endpointEntry.Text, 5000)
			if err := client.Ping(); err != nil {
				configStatusLabel.SetText("Fehler: " + err.Error())
			} else {
				configStatusLabel.SetText("Verbindung OK")
			}
		}()
	})

	newBtn := widget.NewButton("+ Neu", func() {
		nameEntry := widget.NewEntry()
		nameEntry.SetPlaceHolder("z.B. gleis1")
		dialog.ShowCustomConfirm("Neue Konfiguration", "Erstellen", "Abbrechen",
			container.NewVBox(widget.NewLabel("Name:"), nameEntry),
			func(ok bool) {
				if !ok || nameEntry.Text == "" {
					return
				}
				if err := config.Save(nameEntry.Text, &config.Config{Endpoint: "http://"}); err != nil {
					configStatusLabel.SetText("Fehler: " + err.Error())
					return
				}
				refreshList()
				for i, c := range configs {
					if c == nameEntry.Text {
						configList.Select(i)
						loadConfig(nameEntry.Text)
						break
					}
				}
			}, w)
	})

	deleteBtn := widget.NewButton("Löschen", func() {
		if selectedConfig == "default" {
			dialog.ShowInformation("Nicht möglich", "Die default-Konfiguration kann nicht gelöscht werden.", w)
			return
		}
		dialog.ShowConfirm("Löschen", "Konfiguration '"+selectedConfig+"' wirklich löschen?",
			func(ok bool) {
				if !ok {
					return
				}
				if err := config.Delete(selectedConfig); err != nil {
					configStatusLabel.SetText("Fehler: " + err.Error())
					return
				}
				refreshList()
				if len(configs) > 0 {
					configList.Select(0)
					loadConfig(configs[0])
				}
			}, w)
	})

	// Left panel layout: list fills space, buttons + endpoint editor at bottom
	leftBottom := container.NewVBox(
		container.NewHBox(newBtn, deleteBtn),
		widget.NewSeparator(),
		widget.NewForm(widget.NewFormItem("Endpoint",
			container.NewBorder(nil, nil, nil, widget.NewLabel("  "), endpointEntry),
		)),
		container.NewHBox(saveBtn, testConnBtn),
		configStatusLabel,
	)
	leftPanel := container.NewBorder(nil, leftBottom, nil, nil, configList)

	// ── Right panel: CLI builder ──────────────────────────────────────────────
	cmdSelect.OnChanged = func(_ string) { updateDynamic() }
	gleisGroup.OnChanged = func(_ string) { refreshCmd() }
	fullPathCheck.OnChanged = func(_ bool) { refreshCmd() }
	timeEntry.OnChanged = func(_ string) { refreshCmd() }
	for _, e := range []*widget.Entry{
		train1Nr, train1Time, train1Dest, train1Via, train1Delay, train1Info,
		train2Nr, train2Time, train2Dest, train2Via, train2Delay, train2Info,
		train3Nr, train3Time, train3Dest, train3Via, train3Delay, train3Info,
	} {
		e.OnChanged = func(_ string) { refreshCmd() }
	}
	imageFile.OnChanged = func(_ string) { refreshCmd() }

	copyBtn := widget.NewButtonWithIcon("Kopieren", theme.ContentCopyIcon(), func() {
		w.Clipboard().SetContent(cmdOutput.Text)
	})

	executeBtn := widget.NewButtonWithIcon("Ausführen", theme.MediaPlayIcon(), func() {
		execStatus.SetText("Gestartet...")
		go func() {
			exePath, err := os.Executable()
			if err != nil {
				execStatus.SetText("Fehler: " + err.Error())
				return
			}
			// Parse command text (skip first token = executable name)
			tokens := splitCommandLine(cmdOutput.Text)
			if len(tokens) < 2 {
				execStatus.SetText("Kein Befehl angegeben")
				return
			}
			args := tokens[1:]
			cmd := exec.Command(exePath, args...)
			if err := cmd.Start(); err != nil {
				execStatus.SetText("Fehler: " + err.Error())
				return
			}
			// The foreground process spawns --bg and exits quickly
			cmd.Wait()
			execStatus.SetText("Gestartet — siehe debug.log")
		}()
	})

	cmdRow := container.NewBorder(nil, nil, nil, copyBtn, cmdOutput)

	rightTop := container.NewVBox(
		widget.NewRichTextFromMarkdown("### CLI Befehlsgenerator"),
		widget.NewSeparator(),
		widget.NewForm(
			widget.NewFormItem("Befehl", cmdSelect),
			widget.NewFormItem("Gleis", gleisGroup),
		),
	)

	actionRow := container.NewBorder(nil, nil, fullPathCheck, container.NewHBox(execStatus, executeBtn), nil)

	rightBottom := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("Befehl zum Kopieren:"),
		cmdRow,
		actionRow,
	)

	// Border layout: top and bottom sections are fixed, dynamic fields fill the middle
	rightPanel := container.NewBorder(rightTop, rightBottom, nil, nil,
		container.NewScroll(dynamicArea),
	)

	// ── Initialize ────────────────────────────────────────────────────────────
	if len(configs) > 0 {
		loadConfig(configs[0])
		configList.Select(0)
	}
	updateDynamic()

	split := container.NewHSplit(leftPanel, rightPanel)
	split.SetOffset(0.25)

	w.SetContent(split)
	w.ShowAndRun()
}
