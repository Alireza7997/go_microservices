package color

const (
	// Reset resets all text attributes including color.
	Reset = "\033[0m"

	// Text formatting
	Bold      = "\033[1m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
	Reverse   = "\033[7m"
	Concealed = "\033[8m"

	// Regular text colors
	Black   = "\033[30m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	// Background colors
	BlackBg   = "\033[40m"
	RedBg     = "\033[41m"
	GreenBg   = "\033[42m"
	YellowBg  = "\033[43m"
	BlueBg    = "\033[44m"
	MagentaBg = "\033[45m"
	CyanBg    = "\033[46m"
	WhiteBg   = "\033[47m"
)
