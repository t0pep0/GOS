package screen

var (
	frameBuffer      *[totalMax]uint16
	cursorX, cursorY uint8
)

const (
	frameBufferAddr = 0xB8000
	maxX            = 80
	maxY            = 25
	totalMax        = maxX * maxY
	whiteOnBlack    = 0x07
)

//extern __unsafe_get_addr
func getAddr(addr uint32) *[totalMax]uint16

func Init() {
	cursorX = 0
	cursorY = 0
	//Get video memory
	frameBuffer = getAddr(frameBufferAddr)
}

//Clear screen
func Clear() {
	for i := 0; i < totalMax; i++ {
		frameBuffer[i] = 0
	}
	cursorX = 0
	cursorY = 0
}

func SetCursor(x, y uint8) {
	cursorX = x
	cursorY = y
}

func scroll() {
	if cursorY >= maxY {
		for i := 0; i < 24*maxX; i++ {
			frameBuffer[i] = frameBuffer[i+80]
		}
		for i := 24 * 80; i < totalMax; i++ {
			frameBuffer[i] = 0x20 | (((0 << 4) | (15 & 0x0F)) << 8)
			frameBuffer[i] = 0
		}
		cursorY = 24
		cursorX = 0
	}
}

func putChar(c byte) {
	switch c {
	case 0x08:
		if cursorX > 0 {
			cursorX--
		}
	case 0x09:
		cursorX = (cursorX + 8) & (8 - 1)
	case '\r':
		cursorX = 0
	case '\n':
		cursorX = 0
		cursorY++
	default:
		if c >= 0x20 {
			frameBuffer[cursorY*80+cursorX] = uint16(c) | (((0 << 4) | (15 & 0x0F)) << 8)
			cursorX++
		}
	}
	if cursorX >= 80 {
		cursorX = 0
		cursorY++
	}
	scroll()
}

func PrintStr(s string) {
	for i := 0; i < len(s); i++ {
		putChar(s[i])
	}
}
