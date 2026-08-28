package main

import (
	gfx "PICO/TINY_GFX"
	"fmt"
	"machine"
	"math/rand"
	"time"

	"tinygo.org/x/drivers/sh1106"
)

const seed = 42

//go:noinline
func TestRectangles(display *gfx.TinyGFX) time.Duration {
	prng := rand.New(rand.NewSource(seed))

	start_time := time.Now()

	display.Clear()
	for i := 0; i < 64; i++ {
		x := int16(prng.Intn(128))
		y := int16(prng.Intn(64))
		width := int16(prng.Intn(128))
		height := int16(prng.Intn(64))

		display.DrawRectangle(x, y, width, height, false, gfx.WHITE)
		display.Update()
	}

	return time.Since(start_time)
}

//go:noinline
func TestFilledRectangles(display *gfx.TinyGFX) time.Duration {
	prng := rand.New(rand.NewSource(seed))

	start_time := time.Now()

	display.Clear()
	for i := 0; i < 64; i++ {
		x := int16(prng.Intn(128))
		y := int16(prng.Intn(64))
		width := int16(prng.Intn(128))
		height := int16(prng.Intn(64))

		display.DrawRectangle(x, y, width, height, true, gfx.WHITE)
		display.Update()
	}

	return time.Since(start_time)
}

//go:noinline
func TestCircles(display *gfx.TinyGFX) time.Duration {
	start_time := time.Now()

	display.Clear()

	for i := int16(0); i < display.VirtualWidth()/2-5; i++ {
		display.DrawCircle(display.VirtualWidth()/2, display.VirtualHeight()/2, i, false, gfx.WHITE)
		display.Update()
	}

	return time.Since(start_time)
}

//go:noinline
func TestFilledCircles(display *gfx.TinyGFX) time.Duration {
	start_time := time.Now()

	display.Clear()

	for i := int16(0); i < display.VirtualWidth()/2-5; i++ {
		display.DrawCircle(display.VirtualWidth()/2, display.VirtualHeight()/2, i, true, gfx.WHITE)
		display.Update()
	}

	return time.Since(start_time)
}

//go:noinline
func TestInverse(display *gfx.TinyGFX) time.Duration {
	start_time := time.Now()

	invert := false

	for i := 0; i < 3; i++ {
		display.SetInverse(invert)

		display.Clear()
		display.DrawCircle(display.VirtualWidth()/2, display.VirtualHeight()/2, 50, true, gfx.WHITE)
		display.Update()

		invert = !invert

		time.Sleep(500 * time.Millisecond)
	}

	return time.Since(start_time)
}

//go:noinline
func TestLines(display *gfx.TinyGFX) time.Duration {
	start_time := time.Now()

	display.Clear()

	for i := int16(0); i < display.VirtualHeight(); i += 4 {
		display.DrawLine(0, 0, display.VirtualWidth(), i, gfx.WHITE)
		display.Update()
	}

	for i := int16(0); i < display.VirtualWidth(); i += 4 {
		display.DrawLine(0, 0, i, display.VirtualHeight(), gfx.WHITE)
		display.Update()
	}

	return time.Since(start_time)
}

//go:noinline
func TestText(display *gfx.TinyGFX) time.Duration {
	start_time := time.Now()

	display.Clear()

	display.SetTextSize(1)
	display.SetCursorLine(0, 10)
	display.DrawText("Hello, World!", false, gfx.WHITE)
	display.Update()

	display.SetTextSize(2)
	display.SetCursor(0, 20)
	display.DrawText("3.14159", false, gfx.WHITE)
	display.Update()

	display.SetTextSize(3)
	display.SetCursor(0, 40)
	display.DrawText("TinyGFX", false, gfx.WHITE)
	display.Update()

	return time.Since(start_time)
}

//go:noinline
func TestOrientation(display *gfx.TinyGFX) time.Duration {
	start_time := time.Now()

	display.SetTextSize(1)
	display.SetRotation(gfx.Rotation90)
	display.Clear()
	display.SetCursor(0, 10)
	display.DrawText("Rotation90", true, gfx.WHITE)
	display.Update()

	time.Sleep(500 * time.Millisecond)

	display.SetRotation(gfx.Rotation180)
	display.Clear()
	display.SetCursor(0, 10)
	display.DrawText("Rotation180", true, gfx.WHITE)
	display.Update()

	time.Sleep(500 * time.Millisecond)

	display.SetRotation(gfx.Rotation270)
	display.Clear()
	display.SetCursor(0, 10)
	display.DrawText("Rotation270", true, gfx.WHITE)
	display.Update()

	time.Sleep(500 * time.Millisecond)

	display.SetRotation(gfx.Rotation0)
	display.Clear()
	display.SetCursor(0, 10)
	display.DrawText("Rotation0", true, gfx.WHITE)
	display.Update()

	return time.Since(start_time)
}

func main() {
	bus := machine.I2C0
	bus.Configure(machine.I2CConfig{
		Frequency: 800000,
		SDA:       machine.Pin(4),
		SCL:       machine.Pin(5),
	})

	display_driver := sh1106.NewI2C(bus)
	display_driver.Configure(sh1106.Config{
		Width:  128,
		Height: 64,
	})

	display := gfx.NewTinyGFX(&display_driver)

	rect := TestRectangles(display)
	time.Sleep(1000 * time.Millisecond)
	frect := TestFilledRectangles(display)
	time.Sleep(1000 * time.Millisecond)
	circ := TestCircles(display)
	time.Sleep(1000 * time.Millisecond)
	fcirc := TestFilledCircles(display)
	time.Sleep(1000 * time.Millisecond)
	inve := TestInverse(display)
	time.Sleep(1000 * time.Millisecond)
	lines := TestLines(display)
	time.Sleep(1000 * time.Millisecond)
	text := TestText(display)
	time.Sleep(1000 * time.Millisecond)
	orie := TestOrientation(display)

	display.Clear()
	display.SetTextSize(1)
	display.SetCursorLine(0, 1)
	display.DrawText(fmt.Sprintf("Rect: %dms", rect.Milliseconds()), false, gfx.WHITE)
	display.SetCursorLine(0, 2)
	display.DrawText(fmt.Sprintf("FRect: %dms", frect.Milliseconds()), false, gfx.WHITE)
	display.SetCursorLine(0, 3)
	display.DrawText(fmt.Sprintf("Circ: %dms", circ.Milliseconds()), false, gfx.WHITE)
	display.SetCursorLine(0, 4)
	display.DrawText(fmt.Sprintf("FCirc: %dms", fcirc.Milliseconds()), false, gfx.WHITE)
	display.SetCursorLine(0, 5)
	display.DrawText(fmt.Sprintf("Inve: %dms", inve.Milliseconds()), false, gfx.WHITE)
	display.SetCursorLine(0, 6)
	display.DrawText(fmt.Sprintf("Lines: %dms", lines.Milliseconds()), false, gfx.WHITE)
	display.SetCursorLine(0, 7)
	display.DrawText(fmt.Sprintf("Text: %dms", text.Milliseconds()), false, gfx.WHITE)
	display.SetCursorLine(0, 8)
	display.DrawText(fmt.Sprintf("Orie: %dms", orie.Milliseconds()), false, gfx.WHITE)
	display.Update()

	for {
	}
}
