package colorutil

import "github.com/gookit/color"

type ColorFunc func(...interface{}) string

var (
	ColorRed     = ColorFunc(color.FgRed.Render)
	ColorYellow  = ColorFunc(color.FgYellow.Render)
	ColorGreen   = ColorFunc(color.FgGreen.Render)
	ColorBlue    = ColorFunc(color.FgBlue.Render)
	ColorMagenta = ColorFunc(color.FgMagenta.Render)
	ColorCyan    = ColorFunc(color.FgCyan.Render)
	ColorWhite   = ColorFunc(color.FgLightWhite.Render)
)
