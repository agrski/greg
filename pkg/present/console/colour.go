package console

// ANSI colour codes

type AnsiCode string

const (
	escape     AnsiCode = "\u001b"
	csi        AnsiCode = "["
	codePrefix          = escape + csi
	codeSuffix          = "m"
)

const (
	Reset AnsiCode = codePrefix + "0" + codeSuffix

	Bold          AnsiCode = codePrefix + "1" + codeSuffix
	Faint         AnsiCode = codePrefix + "2" + codeSuffix
	Italic        AnsiCode = codePrefix + "3" + codeSuffix
	Underline     AnsiCode = codePrefix + "4" + codeSuffix
	Invert        AnsiCode = codePrefix + "7" + codeSuffix
	Conceal       AnsiCode = codePrefix + "8" + codeSuffix
	Strikethrough AnsiCode = codePrefix + "9" + codeSuffix

	FgBlack   AnsiCode = codePrefix + "30" + codeSuffix
	FgRed     AnsiCode = codePrefix + "31" + codeSuffix
	FgGreen   AnsiCode = codePrefix + "32" + codeSuffix
	FgYellow  AnsiCode = codePrefix + "33" + codeSuffix
	FgBlue    AnsiCode = codePrefix + "34" + codeSuffix
	FgMagenta AnsiCode = codePrefix + "35" + codeSuffix
	FgCyan    AnsiCode = codePrefix + "36" + codeSuffix
	FgWhite   AnsiCode = codePrefix + "37" + codeSuffix
	FgDefault AnsiCode = codePrefix + "39" + codeSuffix

	BgBlack   AnsiCode = codePrefix + "40" + codeSuffix
	BgRed     AnsiCode = codePrefix + "41" + codeSuffix
	BgGreen   AnsiCode = codePrefix + "42" + codeSuffix
	BgYellow  AnsiCode = codePrefix + "43" + codeSuffix
	BgBlue    AnsiCode = codePrefix + "44" + codeSuffix
	BgMagenta AnsiCode = codePrefix + "45" + codeSuffix
	BgCyan    AnsiCode = codePrefix + "46" + codeSuffix
	BgWhite   AnsiCode = codePrefix + "47" + codeSuffix

	FgIntenseBlack   AnsiCode = codePrefix + "90" + codeSuffix
	FgIntenseRed     AnsiCode = codePrefix + "91" + codeSuffix
	FgIntenseGreen   AnsiCode = codePrefix + "92" + codeSuffix
	FgIntenseYellow  AnsiCode = codePrefix + "93" + codeSuffix
	FgIntenseBlue    AnsiCode = codePrefix + "94" + codeSuffix
	FgIntenseMagenta AnsiCode = codePrefix + "95" + codeSuffix
	FgIntenseCyan    AnsiCode = codePrefix + "96" + codeSuffix
	FgIntenseWhite   AnsiCode = codePrefix + "97" + codeSuffix

	BgIntenseBlack   AnsiCode = codePrefix + "100" + codeSuffix
	BgIntenseRed     AnsiCode = codePrefix + "101" + codeSuffix
	BgIntenseGreen   AnsiCode = codePrefix + "102" + codeSuffix
	BgIntenseYellow  AnsiCode = codePrefix + "103" + codeSuffix
	BgIntenseBlue    AnsiCode = codePrefix + "104" + codeSuffix
	BgIntenseMagenta AnsiCode = codePrefix + "105" + codeSuffix
	BgIntenseCyan    AnsiCode = codePrefix + "106" + codeSuffix
	BgIntenseWhite   AnsiCode = codePrefix + "107" + codeSuffix
)
