package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type guiState struct {
	scriptPath widget.Editor
	startBtn   widget.Clickable
	stopBtn    widget.Clickable
	status     string
	matcher    *Matcher
	running    bool
}

func runGUI() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("Match-Go"),
			app.Size(unit.Dp(600), unit.Dp(400)),
		)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()
	var ops op.Ops
	state := &guiState{
		status: "Ready",
	}
	state.scriptPath.SetText("script.txt")

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if state.startBtn.Clicked(gtx) && !state.running {
				path := state.scriptPath.Text()
				go func() {
					state.running = true
					state.status = "Loading..."
					w.Invalidate()

					config, err := parseScript(path)
					if err != nil {
						state.status = fmt.Sprintf("Error: %v", err)
						state.running = false
						w.Invalidate()
						return
					}

					dataSets, err := readDataSets(config)
					if err != nil {
						state.status = fmt.Sprintf("Error: %v", err)
						state.running = false
						w.Invalidate()
						return
					}

					preprocessData(config, dataSets)
					state.matcher = NewMatcher(config, dataSets)
					state.matcher.OnUpdate = func(dist float64) {
						state.status = fmt.Sprintf("Best Distance: %.4f", dist)
						w.Invalidate()
					}

					state.status = "Running..."
					w.Invalidate()
					state.matcher.Run()
					
					state.running = false
					state.status = fmt.Sprintf("Finished. Final Dist: %.4f", state.matcher.BestDist)
					w.Invalidate()
				}()
			}

			if state.stopBtn.Clicked(gtx) && state.running && state.matcher != nil {
				state.matcher.Stop()
				state.status = "Stopping..."
			}

			layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return material.H4(th, "Match-Go").Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						e := material.Editor(th, &state.scriptPath, "Script Path")
						state.scriptPath.SingleLine = true
						return e.Layout(gtx)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.Flex{}.Layout(gtx,
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if state.running {
									gtx = gtx.Disabled()
								}
								return material.Button(th, &state.startBtn, "Start").Layout(gtx)
							}),
							layout.Rigid(layout.Spacer{Width: unit.Dp(16)}.Layout),
							layout.Rigid(func(gtx layout.Context) layout.Dimensions {
								if !state.running {
									gtx = gtx.Disabled()
								}
								btn := material.Button(th, &state.stopBtn, "Stop")
								btn.Background = color.NRGBA{R: 200, A: 255}
								return btn.Layout(gtx)
							}),
						)
					})
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.UniformInset(unit.Dp(16)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						lbl := material.Body1(th, state.status)
						lbl.Alignment = text.Middle
						return lbl.Layout(gtx)
					})
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
