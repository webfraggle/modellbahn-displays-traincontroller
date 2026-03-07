package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/webfraggle/mbd-tc/internal/api"
	"github.com/webfraggle/mbd-tc/internal/config"
)

func Run() {
	a := app.New()
	w := a.NewWindow("mbd-tc Konfiguration")
	w.Resize(fyne.NewSize(560, 320))

	configs, _ := config.List()

	endpointEntry := widget.NewEntry()
	endpointEntry.SetPlaceHolder("http://192.168.178.xxx")

	statusLabel := widget.NewLabel("")

	var selectedConfig string

	loadConfig := func(name string) {
		selectedConfig = name
		cfg, err := config.Load(name)
		if err == nil {
			endpointEntry.SetText(cfg.Endpoint)
		}
		statusLabel.SetText("")
	}

	configList := widget.NewList(
		func() int { return len(configs) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(configs[i])
		},
	)

	configList.OnSelected = func(id widget.ListItemID) {
		if id < len(configs) {
			loadConfig(configs[id])
		}
	}

	refreshList := func() {
		configs, _ = config.List()
		configList.Refresh()
	}

	if len(configs) > 0 {
		loadConfig(configs[0])
		configList.Select(0)
	}

	saveBtn := widget.NewButton("Speichern", func() {
		if err := config.Save(selectedConfig, &config.Config{Endpoint: endpointEntry.Text}); err != nil {
			statusLabel.SetText("Fehler: " + err.Error())
		} else {
			statusLabel.SetText("Gespeichert.")
		}
	})

	testBtn := widget.NewButton("Verbindung testen", func() {
		statusLabel.SetText("Teste...")
		go func() {
			client := api.NewClient(endpointEntry.Text, 5000)
			if err := client.Ping(); err != nil {
				statusLabel.SetText("Fehler: " + err.Error())
			} else {
				statusLabel.SetText("Verbindung OK")
			}
		}()
	})

	newBtn := widget.NewButton("Neue Config", func() {
		nameEntry := widget.NewEntry()
		nameEntry.SetPlaceHolder("z.B. gleis1")
		dialog.ShowCustomConfirm("Neue Konfiguration", "Erstellen", "Abbrechen",
			container.NewVBox(widget.NewLabel("Name:"), nameEntry),
			func(ok bool) {
				if !ok || nameEntry.Text == "" {
					return
				}
				cfg := &config.Config{Endpoint: "http://"}
				if err := config.Save(nameEntry.Text, cfg); err != nil {
					statusLabel.SetText("Fehler: " + err.Error())
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
		dialog.ShowConfirm("Löschen",
			"Konfiguration '"+selectedConfig+"' wirklich löschen?",
			func(ok bool) {
				if !ok {
					return
				}
				if err := config.Delete(selectedConfig); err != nil {
					statusLabel.SetText("Fehler: " + err.Error())
					return
				}
				refreshList()
				if len(configs) > 0 {
					configList.Select(0)
					loadConfig(configs[0])
				}
			}, w)
	})

	form := container.NewVBox(
		widget.NewForm(widget.NewFormItem("Endpoint URL", endpointEntry)),
		container.NewHBox(saveBtn, testBtn),
		container.NewHBox(newBtn, deleteBtn),
		statusLabel,
	)

	split := container.NewHSplit(configList, form)
	split.SetOffset(0.28)

	w.SetContent(split)
	w.ShowAndRun()
}
