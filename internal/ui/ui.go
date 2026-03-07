package ui

import (
	"fmt"
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

const exeName = "mbd-cli"

var commandList = []string{"--next", "--prev", "--setTime", "--setTrain1", "--setTrain2", "--setTrain3", "--image"}

func Run() {
	a := app.New()
	w := a.NewWindow("mbd-cli Konfiguration")
	w.Resize(fyne.NewSize(960, 540))

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

	trainNr := widget.NewEntry()
	trainNr.SetPlaceHolder("ICE123")
	trainTime := widget.NewEntry()
	trainTime.SetPlaceHolder("12:30")
	trainDest := widget.NewEntry()
	trainDest.SetPlaceHolder("Berlin")
	trainVia := widget.NewEntry()
	trainVia.SetPlaceHolder("Hannover")
	trainDelay := widget.NewEntry()
	trainDelay.SetPlaceHolder("0")
	trainInfo := widget.NewEntry()
	trainInfo.SetPlaceHolder("Sonderinfo")

	imageFile := widget.NewEntry()
	imageFile.SetPlaceHolder("00Logo.png")

	cmdOutput := widget.NewLabel("")
	cmdOutput.Wrapping = fyne.TextWrapBreak

	execStatus := widget.NewLabel("")
	dynamicArea := container.NewVBox()

	// ── Build command string ──────────────────────────────────────────────────
	buildCommand := func() string {
		var parts []string
		parts = append(parts, exeName)

		switch cmdSelect.Selected {
		case "--next":
			parts = append(parts, "--next")
		case "--prev":
			parts = append(parts, "--prev")
		case "--setTime":
			t := timeEntry.Text
			if t == "" {
				t = "HH:MM"
			}
			parts = append(parts, "--setTime", fmt.Sprintf(`"%s"`, t))
		case "--setTrain1", "--setTrain2", "--setTrain3":
			n := string(cmdSelect.Selected[len(cmdSelect.Selected)-1])
			s := strings.Join([]string{
				trainNr.Text, trainTime.Text, trainDest.Text,
				trainVia.Text, trainDelay.Text, trainInfo.Text,
			}, "|")
			parts = append(parts, "--setTrain"+n, fmt.Sprintf(`"%s"`, s))
		case "--image":
			f := imageFile.Text
			if f == "" {
				f = "filename.png"
			}
			parts = append(parts, "--image", f)
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
		cmdOutput.Refresh()
	}

	// ── Update dynamic fields based on selected command ───────────────────────
	updateDynamic := func() {
		var items []fyne.CanvasObject
		switch cmdSelect.Selected {
		case "--setTime":
			items = []fyne.CanvasObject{
				widget.NewForm(widget.NewFormItem("Zeit", timeEntry)),
			}
		case "--setTrain1", "--setTrain2", "--setTrain3":
			items = []fyne.CanvasObject{
				widget.NewForm(
					widget.NewFormItem("Zug-Nr", trainNr),
					widget.NewFormItem("Zeit", trainTime),
					widget.NewFormItem("Ziel", trainDest),
					widget.NewFormItem("Via", trainVia),
					widget.NewFormItem("Verspätung", trainDelay),
					widget.NewFormItem("Hinweis", trainInfo),
				),
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
		widget.NewForm(widget.NewFormItem("Endpoint", endpointEntry)),
		container.NewHBox(saveBtn, testConnBtn),
		configStatusLabel,
	)
	leftPanel := container.NewBorder(nil, leftBottom, nil, nil, configList)

	// ── Right panel: CLI builder ──────────────────────────────────────────────
	cmdSelect.OnChanged = func(_ string) { updateDynamic() }
	gleisGroup.OnChanged = func(_ string) { refreshCmd() }
	timeEntry.OnChanged = func(_ string) { refreshCmd() }
	trainNr.OnChanged = func(_ string) { refreshCmd() }
	trainTime.OnChanged = func(_ string) { refreshCmd() }
	trainDest.OnChanged = func(_ string) { refreshCmd() }
	trainVia.OnChanged = func(_ string) { refreshCmd() }
	trainDelay.OnChanged = func(_ string) { refreshCmd() }
	trainInfo.OnChanged = func(_ string) { refreshCmd() }
	imageFile.OnChanged = func(_ string) { refreshCmd() }

	copyBtn := widget.NewButtonWithIcon("Kopieren", theme.ContentCopyIcon(), func() {
		w.Clipboard().SetContent(buildCommand())
	})

	executeBtn := widget.NewButtonWithIcon("Ausführen", theme.MediaPlayIcon(), func() {
		execStatus.SetText("Läuft...")
		go func() {
			cfg, err := config.Load(selectedConfig)
			if err != nil {
				execStatus.SetText("Fehler: " + err.Error())
				return
			}
			client := api.NewClient(cfg.Endpoint, 10000)
			gleis := "GleisA"
			if gleisGroup.Selected == "B" {
				gleis = "GleisB"
			}
			switch cmdSelect.Selected {
			case "--next":
				err = client.SkipNext(gleis)
			case "--prev":
				err = client.SkipPrev(gleis)
			case "--setTime":
				err = client.SetTime(gleis, timeEntry.Text)
			case "--setTrain1", "--setTrain2", "--setTrain3":
				n := int(cmdSelect.Selected[len(cmdSelect.Selected)-1] - '0')
				slot := api.ParseTrain(n, strings.Join([]string{
					trainNr.Text, trainTime.Text, trainDest.Text,
					trainVia.Text, trainDelay.Text, trainInfo.Text,
				}, "|"), gleis)
				err = client.SetTrains([]api.TrainSlot{slot})
			case "--image":
				err = client.ShowImage(gleis, imageFile.Text)
			}
			if err != nil {
				execStatus.SetText("Fehler: " + err.Error())
			} else {
				execStatus.SetText("Erfolgreich")
			}
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

	rightBottom := container.NewVBox(
		widget.NewSeparator(),
		widget.NewLabel("Befehl zum Kopieren:"),
		cmdRow,
		container.NewHBox(executeBtn, execStatus),
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
