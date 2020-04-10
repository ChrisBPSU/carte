package carte

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)
const (
	colorBrightBlack = iota + 90
	colorBrightRed
	colorBrightGreen
	colorBrightYellow
	colorBrightBlue
	colorBrightMagenta
	colorBrightCyan
	colorBrightWhite
)

// TODO: allow the user to assign colors to the different log types

//colorCode: []byte(fmt.Sprintf("\033[%dm", colorWhite))

// INFO white
// DEBG green
// WARN yellow

// ERR  red
// CRIT purple
