package tinygfx

import (
	"fmt"
	"image/color"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinyfont"
)

var (
	WHITE   = color.RGBA{255, 255, 255, 255}
	BLACK   = color.RGBA{0, 0, 0, 255}
	RED     = color.RGBA{255, 0, 0, 255}
	GREEN   = color.RGBA{0, 255, 0, 255}
	BLUE    = color.RGBA{0, 0, 255, 255}
	YELLOW  = color.RGBA{255, 255, 0, 255}
	CYAN    = color.RGBA{0, 255, 255, 255}
	MAGENTA = color.RGBA{255, 0, 255, 255}
	ORANGE  = color.RGBA{255, 165, 0, 255}
	PURPLE  = color.RGBA{128, 0, 128, 255}
	BROWN   = color.RGBA{165, 42, 42, 255}
	GRAY    = color.RGBA{128, 128, 128, 255}
)

type rotation = int

const (
	Rotation0   = 0
	Rotation90  = 1
	Rotation180 = 2
	Rotation270 = 3
)

var (
	SPACER_CONSTANT_X uint = 1
	SPACER_CONSTANT_Y uint = 1
)

/*
Provides basic functions for drawing shapes and printing text
for drivers which follow drivers.Displayer interface signature
*/
type TinyGFX struct {
	DisplayDevice    drivers.Displayer
	font             tinyfont.Fonter
	cursorX, cursorY int16
	fontSize         uint
	rotation         rotation
	inversed         bool
}

func check_out_of_bounds(display *TinyGFX, x, y int16) bool {
	width, height := display.Size()
	return (x > width || y > height) || (x < 0 || y < 0)
}

func check_rotation_supported(r rotation) bool {
	switch r {
	case Rotation0, Rotation90, Rotation180, Rotation270:
		return true
	}
	return false
}

/*
Creates a new TinyGFX Object with the display driver provided
Automatically creates a default configuration and returns a pointer to the object
*/
func NewTinyGFX(display drivers.Displayer) *TinyGFX {
	result := TinyGFX{DisplayDevice: display}
	result.SetFont(&tinyfont.Org01)
	result.SetCursor(0, 10)
	result.SetTextSize(1)
	result.SetRotation(Rotation0)
	result.SetInverse(false)

	return &result
}

/*
Changes renderer font, Uses Font format defined by the tinyfont library
*/
func (display *TinyGFX) SetFont(font tinyfont.Fonter) {
	display.font = font
}

/*
Changes renderer draw orientation, Only supports Rotation provided by the library
*/
func (display *TinyGFX) SetRotation(rot rotation) error {
	if !check_rotation_supported(rot) {
		return fmt.Errorf("rotation not supported")
	}

	display.rotation = rot

	return nil
}

/*
Changes rendered text size
*/
func (display *TinyGFX) SetTextSize(size uint) {
	display.fontSize = size
}

/*
Sets the inverse flag of the renderer
*/
func (display *TinyGFX) SetInverse(inv bool) {
	display.inversed = inv
}

/*
Puts the renderer cursor on a specific coordinate on the drawable region
*/
func (display *TinyGFX) SetCursor(x, y int16) error {
	if check_out_of_bounds(display, x, y) {
		return fmt.Errorf("coordinates given outside of drawable region")
	}

	display.cursorX = x
	display.cursorY = y

	return nil
}

/*
Puts the renderer cursor on specific line and column of the drawable region
Also accounts for variable text size

works best for monospaced fonts
*/
func (display *TinyGFX) SetCursorLine(column, line int16) error {
	var (
		CHARACTER_WIDTH  = int16(display.font.GetGlyph('A').Info().Width) * int16(display.fontSize)
		CHARACTER_HEIGHT = int16(display.font.GetGlyph('A').Info().Height) * int16(display.fontSize)
	)

	if check_out_of_bounds(display, column*CHARACTER_WIDTH, line*CHARACTER_HEIGHT) {
		return fmt.Errorf("coordinates given outside of drawable region")
	}

	display.cursorX = column * CHARACTER_WIDTH
	display.cursorY = (line * CHARACTER_HEIGHT) + ((line % (display.VirtualHeight() / CHARACTER_HEIGHT)) * int16(display.fontSize))

	return nil
}

// Exact same thing as `Update` but allows the object to be used by other libraries by following the `Displayer` function signatures
func (display *TinyGFX) Display() error {
	return display.DisplayDevice.Display()
}

/*
Draw a color on the pixel with (x, y) coordinate with a color
Also accounts for orientation
*/
func (display *TinyGFX) SetPixel(x, y int16, color color.RGBA) {
	switch display.rotation {
	case Rotation90:
		x, y = display.Width()-1-y, x
	case Rotation180:
		x, y = display.Width()-1-x, display.Height()-1-y
	case Rotation270:
		x, y = y, display.Height()-1-x
	}

	if display.inversed {
		color.R = 255 - color.R
		color.G = 255 - color.G
		color.B = 255 - color.B
	}

	display.DisplayDevice.SetPixel(x, y, color)
}

/*
Returns the physical (width, height) of the screen
*/
func (display *TinyGFX) Size() (int16, int16) {
	return display.DisplayDevice.Size()
}

/*
Returns the physical width of the screen
*/
func (display *TinyGFX) Width() int16 {
	width, _ := display.Size()
	return width
}

/*
Returns the physical height of the screen
*/
func (display *TinyGFX) Height() int16 {
	_, height := display.Size()
	return height
}

/*
Returns the Drawable width of the screen
*/
func (display *TinyGFX) VirtualWidth() int16 {
	switch display.rotation {
	case Rotation0, Rotation180:
		return display.Width()
	case Rotation90, Rotation270:
		return display.Height()
	}

	return display.Width()
}

/*
Returns the Drawable height of the screen
*/
func (display *TinyGFX) VirtualHeight() int16 {
	switch display.rotation {
	case Rotation0, Rotation180:
		return display.Height()
	case Rotation90, Rotation270:
		return display.Width()
	}

	return display.Height()
}

/*
Fills the entire screen with the color specified
*/
func (display *TinyGFX) Fill(color color.RGBA) {
	for x := int16(0); x < display.VirtualWidth(); x++ {
		for y := int16(0); y < display.VirtualHeight(); y++ {
			display.SetPixel(x, y, color)
		}
	}
}

/*
An alias for Fill(BLACK)
*/
func (display *TinyGFX) Clear() {
	display.Fill(BLACK)
}

/*
An alias of Display(), doesn't return any error
*/
func (display *TinyGFX) Update() {
	display.DisplayDevice.Display()
}

// Exact same thing as `SetPixel` but follows this library's naming conventions
func (display *TinyGFX) DrawPixel(x, y int16, color color.RGBA) {
	display.SetPixel(x, y, color)
}

// Identical function to tinyfont's DrawChar() but supports various text size
// Function inherited from tinyfont's original github repo to support various text sizes
func (display *TinyGFX) drawChar(char rune, x, y int16, color color.RGBA) {
	glyph := display.font.GetGlyph(char).(tinyfont.Glyph)

	bitmapOffset := 0
	bitmap := byte(0)

	if len(glyph.Bitmaps) > 0 {
		bitmap = glyph.Bitmaps[bitmapOffset]
	}

	bit := uint8(0)
	for j := int16(0); j < int16(glyph.Height); j++ {
		for i := int16(0); i < int16(glyph.Width); i++ {
			if (bitmap & 0x80) != 0x00 {
				display.DrawSquare(
					x+int16(glyph.XOffset)+(int16(display.fontSize)*i),
					y+int16(glyph.YOffset)+(int16(display.fontSize)*j),
					int16(display.fontSize),
					true,
					color,
				)
			}

			bitmap <<= 1
			bit++

			if bit > 7 {
				bitmapOffset++
				if bitmapOffset < len(glyph.Bitmaps) {
					bitmap = glyph.Bitmaps[bitmapOffset]
				}
				bit = 0
			}
		}
	}
}

/*
Draws text onto the screen at (cursorX, cursorY) coordinates stored inside the TinyGFX Object
The font used is also stored inside TinyGFX object which can be changed with SetFont() method
*/
func (display *TinyGFX) DrawText(text string, wrap bool, color color.RGBA) {
	drawX := display.cursorX
	drawY := display.cursorY

	if check_out_of_bounds(display, drawX, drawY) {
		return
	}

	for _, char := range text {
		var (
			CHARACTER_WIDTH  = int16(display.font.GetGlyph(char).Info().Width) * int16(display.fontSize)
			CHARACTER_HEIGHT = int16(display.font.GetGlyph(char).Info().Height) * int16(display.fontSize)
			SPACER_X         = int16(uint(SPACER_CONSTANT_X) * display.fontSize)
			SPACER_Y         = int16(uint(SPACER_CONSTANT_Y) * display.fontSize)
		)

		display.drawChar(char, drawX, drawY, color)

		if drawX+CHARACTER_WIDTH+SPACER_X > display.VirtualWidth() {
			if !wrap {
				return
			}

			drawY += CHARACTER_HEIGHT + SPACER_Y
			drawX = 0

			if drawY > display.VirtualHeight() {
				return
			}
		} else {
			drawX += CHARACTER_WIDTH + SPACER_X
		}
	}

	return
}

/*
Draws a bitmap image stored as []byte onto screen
The bitmap image should be encoded as Horizontal - 1 bit per pixel
*/
func (display *TinyGFX) DrawBitmap(bitmap []byte, width, height, x, y int16, inverse bool) {
	DISPLAY_WIDTH, DISPLAY_HEIGHT := display.DisplayDevice.Size()

	drawX := x
	drawY := y

	var white color.RGBA = WHITE
	var black color.RGBA = BLACK

	if inverse {
		white = BLACK
		black = WHITE
	}

	for _, item := range bitmap {
		if drawY > DISPLAY_HEIGHT {
			return
		}

		for i := uint(0); i < 8; i++ {
			if drawX < DISPLAY_WIDTH {
				if item&(1<<(7-i)) != 0 {
					display.SetPixel(drawX, drawY, white)
				} else {
					display.SetPixel(drawX, drawY, black)
				}
			}

			drawX++
			if drawX >= width+x {
				drawX = x
				drawY++
			}
		}
	}
}

/*
Draws a filled/unfilled rectangle onto the screen at (x, y) coordinate with the given color
*/
func (display *TinyGFX) DrawRectangle(x, y, width, height int16, fill bool, color color.RGBA) {
	if check_out_of_bounds(display, x, y) {
		return
	}

	if width < 0 || height < 0 {
		return
	}

	if width+x > display.VirtualWidth() {
		width = display.VirtualWidth() - x
	}

	if height+y > display.VirtualHeight() {
		height = display.VirtualHeight() - y
	}

	for drawX := x; drawX < width+x; drawX++ {
		for drawY := y; drawY < height+y; drawY++ {
			if fill {
				display.SetPixel(drawX, drawY, color)
			} else if drawX == (width+x)-1 || drawX == x || drawY == (height+y)-1 || drawY == y {
				display.SetPixel(drawX, drawY, color)
			}
		}
	}
}

/*
Draws a Horizontal line starting from (x, y) coordinate to (x + length, y) with the given color
*/
func (display *TinyGFX) DrawHLine(x, y, length int16, color color.RGBA) {
	if length < 0 {
		return
	}

	if length+x > display.Width() {
		length = display.Width() - x
	}

	for drawX := x; drawX < length+x; drawX++ {
		display.SetPixel(drawX, y, color)
	}
}

/*
Draws a Vertical line starting from (x, y) coordinate to (x, y + length) with the given color
*/
func (display *TinyGFX) DrawVLine(x, y, length int16, color color.RGBA) {
	if length < 0 {
		return
	}

	if length+y > display.Height() {
		length = display.Height() - y
	}

	for drawY := y; drawY < length+y; drawY++ {
		display.SetPixel(x, drawY, color)
	}
}

/*
Draws a filled/unfilled Circle on the screen using Bresenham's Circle Algo. at (x, y) coordinate with the given color
*/
func (display *TinyGFX) DrawCircle(x, y, radius int16, fill bool, color color.RGBA) {
	if radius < 0 {
		return
	}

	drawX := int16(0)
	drawY := radius
	d := int16(3) - 2*radius

	plotCirclePoints := func(cx, cy, px, py int16) {
		if fill {
			for drawX <= drawY {
				display.DrawHLine(x-drawX, y+drawY, 2*drawX+1, color)
				display.DrawHLine(x-drawX, y-drawY, 2*drawX+1, color)
				display.DrawHLine(x-drawY, y+drawX, 2*drawY+1, color)
				display.DrawHLine(x-drawY, y-drawX, 2*drawY+1, color)

				if d > 0 {
					drawY--
					d = d + 4*int16(drawX-drawY) + 10
				} else {
					d = d + 4*int16(drawX) + 6
				}
				drawX++
			}
		} else {
			display.SetPixel(cx+px, cy+py, color)
			display.SetPixel(cx-px, cy+py, color)
			display.SetPixel(cx+px, cy-py, color)
			display.SetPixel(cx-px, cy-py, color)
			display.SetPixel(cx+py, cy+px, color)
			display.SetPixel(cx+py, cy-px, color)
			display.SetPixel(cx-py, cy+px, color)
			display.SetPixel(cx-py, cy-px, color)
		}
	}

	plotCirclePoints(x, y, drawX, drawY)

	for drawX <= drawY {
		drawX++

		if d > 0 {
			drawY--
			d = d + 4*(drawX-drawY) + 10
		} else {
			d = d + 4*drawX + 6
		}

		plotCirclePoints(x, y, drawX, drawY)
	}
}

/*
An alias for DrawRectangle(x, y, side, side, fill, color)
*/
func (display *TinyGFX) DrawSquare(x, y, side int16, fill bool, color color.RGBA) {
	display.DrawRectangle(x, y, side, side, fill, color)
}

/*
Draw a line on the screen using Bresenham's Line Algo. starting from (x1, y1) ending at (x2, y2) with the given color
*/
func (display *TinyGFX) DrawLine(x1, y1, x2, y2 int16, color color.RGBA) {
	abs := func(x int16) int16 {
		if x < 0 {
			return -x
		}
		return x
	}

	dx := abs(x2 - x1)
	dy := -abs(y2 - y1)

	var sx, sy int16 = 1, 1
	if x1 > x2 {
		sx = -1
	}
	if y1 > y2 {
		sy = -1
	}

	err := dx + dy

	for {
		display.SetPixel(x1, y1, color)
		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x1 += sx
		}
		if e2 <= dx {
			err += dx
			y1 += sy
		}
	}
}

/*
Draws a multi-point line on the screen using DrawLine() starting from the first pair of (x, y) coordinates ending at the last pair of (x, y) with the given color
*/
func (display *TinyGFX) DrawMultiPointLine(color color.RGBA, points ...int16) {
	length := len(points)

	if length < 4 {
		return
	}

	for i := 0; i < length-2; i += 2 {
		x1 := points[i]
		y1 := points[i+1]
		x2 := points[i+2]
		y2 := points[i+3]

		display.DrawLine(x1, y1, x2, y2, color)
	}
}
