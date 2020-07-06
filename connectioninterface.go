package ledstrip

// LedConnection is a generic interface for additional connection types (e.g. PWM)
type LedConnection interface {
	RenderLEDs(pixels []RGBPixel)
	getTranslatedColor(pixel [3]uint8) []uint32
	Clear(nr int)
	Close()
}
