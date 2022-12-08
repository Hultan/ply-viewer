package ply_viewer

import (
	"image/color"
	"os"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"

	"github.com/hultan/ply-viewer/internal/ply"
	"github.com/hultan/softteam/framework"
)

const applicationTitle = "go3d"
const applicationVersion = "v 0.01"
const applicationCopyRight = "©SoftTeam AB, 2022"

var (
	WHITE  = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	RED    = color.RGBA{R: 185, G: 0, B: 0, A: 255}
	BLUE   = color.RGBA{R: 0, G: 69, B: 173, A: 255}
	GREEN  = color.RGBA{R: 0, G: 155, B: 72, A: 255}
	ORANGE = color.RGBA{R: 255, G: 89, B: 0, A: 255}
	YELLOW = color.RGBA{R: 255, G: 213, B: 0, A: 255}
)

var thetaX float64 = 0.1
var thetaY float64 = 0.6
var thetaZ float64 = 0.6
var dX = 0.01
var dY = 0.02
var dZ = 0.03
var object ply.PLY

type Axis int

const (
	XAxis Axis = iota
	YAxis
	ZAxis
)

type RotationMatrix [3][3]float32

type MainForm struct {
	Window      *gtk.ApplicationWindow
	builder     *framework.GtkBuilder
	AboutDialog *gtk.AboutDialog
}

// NewMainForm : Creates a new MainForm object
func NewMainForm() *MainForm {
	mainForm := new(MainForm)
	return mainForm
}

// OpenMainForm : Opens the MainForm window
func (m *MainForm) OpenMainForm(app *gtk.Application) {
	// Initialize gtk
	gtk.Init(&os.Args)

	// Create a new softBuilder
	fw := framework.NewFramework()
	builder, err := fw.Gtk.CreateBuilder("main.glade")
	if err != nil {
		panic(err)
	}
	m.builder = builder

	// Get the main window from the glade file
	m.Window = m.builder.GetObject("main_window").(*gtk.ApplicationWindow)

	// Set up main window
	m.Window.SetApplication(app)
	m.Window.SetTitle("go3d main window")
	m.Window.SetSizeRequest(1024, 768)
	// Hook up the destroy event
	m.Window.Connect("destroy", m.Window.Close)

	// Quit button
	button := m.builder.GetObject("main_window_quit_button").(*gtk.ToolButton)
	button.Connect("clicked", m.Window.Close)

	// Status bar
	statusBar := m.builder.GetObject("main_window_status_bar").(*gtk.Statusbar)
	statusBar.Push(statusBar.GetContextId("go3d"), "go3d : version 0.1.0")

	// Drawing area
	da := builder.GetObject("drawingArea").(*gtk.DrawingArea)
	da.Connect("draw", m.draw)
	// Menu
	m.setupMenu(fw)

	// Show the main window
	m.Window.ShowAll()

	err = object.Load("/home/per/code/ply-viewer/data/dragon.ply")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			// time.Sleep(50 * time.Millisecond)
			da.QueueDraw()
		}
	}()
}

func (m *MainForm) setupMenu(fw *framework.Framework) {
	menuQuit := m.builder.GetObject("menu_file_quit").(*gtk.MenuItem)
	menuQuit.Connect("activate", m.Window.Close)
}

func (m *MainForm) draw(da *gtk.DrawingArea, cx *cairo.Context) {
	m.drawBackground(da, cx)
	m.rotateObject()
	m.drawFaces(cx)
}

func (m *MainForm) drawBackground(da *gtk.DrawingArea, cx *cairo.Context) {
	cx.SetSourceRGBA(0, 0, 0, 1)
	w, h := da.GetAllocatedWidth(), da.GetAllocatedHeight()
	cx.Rectangle(0, 0, float64(w), float64(h))
	cx.Fill()
}

func (m *MainForm) drawFaces(cx *cairo.Context) {
	for _, face := range object.Faces {
		var v []*ply.Vertex
		for i := 0; i < len(face.Indexes); i++ {
			v = append(v, object.Vertexes[face.Indexes[i]])
		}
		m.drawFace(cx, v...)
	}
}

func (m *MainForm) drawFace(cx *cairo.Context, args ...*ply.Vertex) {
	zMean := args[0].Z
	cx.MoveTo(coord(args[0].X, args[0].Z), coord(args[0].Y, args[0].Z))
	for i := 0; i < len(args); i++ {
		zMean += args[i].Z
		cx.LineTo(coord(args[i].X, args[i].Z), coord(args[i].Y, args[i].Z))
	}
	cx.SetSourceRGBA(float64((zMean+2500)/5000), float64((zMean+2500)/5000), float64((zMean+2500)/5000), 1)
	cx.Fill()
}

func coord(c, z float32) float64 {
	return float64(c + 400)
}
