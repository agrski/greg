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
	reset AnsiCode = codePrefix + "0" + codeSuffix

	bold          AnsiCode = codePrefix + "1" + codeSuffix
	faint         AnsiCode = codePrefix + "2" + codeSuffix
	italic        AnsiCode = codePrefix + "3" + codeSuffix
	underline     AnsiCode = codePrefix + "4" + codeSuffix
	invert        AnsiCode = codePrefix + "7" + codeSuffix
	conceal       AnsiCode = codePrefix + "8" + codeSuffix
	strikethrough AnsiCode = codePrefix + "9" + codeSuffix

	fgBlack   AnsiCode = codePrefix + "30" + codeSuffix
	fgRed     AnsiCode = codePrefix + "31" + codeSuffix
	fgGreen   AnsiCode = codePrefix + "32" + codeSuffix
	fgYellow  AnsiCode = codePrefix + "33" + codeSuffix
	fgBlue    AnsiCode = codePrefix + "34" + codeSuffix
	fgMagenta AnsiCode = codePrefix + "35" + codeSuffix
	fgCyan    AnsiCode = codePrefix + "36" + codeSuffix
	fgWhite   AnsiCode = codePrefix + "37" + codeSuffix
	fgDefault AnsiCode = codePrefix + "39" + codeSuffix

	bgBlack   AnsiCode = codePrefix + "40" + codeSuffix
	bgRed     AnsiCode = codePrefix + "41" + codeSuffix
	bgGreen   AnsiCode = codePrefix + "42" + codeSuffix
	bgYellow  AnsiCode = codePrefix + "43" + codeSuffix
	bgBlue    AnsiCode = codePrefix + "44" + codeSuffix
	bgMagenta AnsiCode = codePrefix + "45" + codeSuffix
	bgCyan    AnsiCode = codePrefix + "46" + codeSuffix
	bgWhite   AnsiCode = codePrefix + "47" + codeSuffix

	fgIntenseBlack   AnsiCode = codePrefix + "90" + codeSuffix
	fgIntenseRed     AnsiCode = codePrefix + "91" + codeSuffix
	fgIntenseGreen   AnsiCode = codePrefix + "92" + codeSuffix
	fgIntenseYellow  AnsiCode = codePrefix + "93" + codeSuffix
	fgIntenseBlue    AnsiCode = codePrefix + "94" + codeSuffix
	fgIntenseMagenta AnsiCode = codePrefix + "95" + codeSuffix
	fgIntenseCyan    AnsiCode = codePrefix + "96" + codeSuffix
	fgIntenseWhite   AnsiCode = codePrefix + "97" + codeSuffix

	bgIntenseBlack   AnsiCode = codePrefix + "100" + codeSuffix
	bgIntenseRed     AnsiCode = codePrefix + "101" + codeSuffix
	bgIntenseGreen   AnsiCode = codePrefix + "102" + codeSuffix
	bgIntenseYellow  AnsiCode = codePrefix + "103" + codeSuffix
	bgIntenseBlue    AnsiCode = codePrefix + "104" + codeSuffix
	bgIntenseMagenta AnsiCode = codePrefix + "105" + codeSuffix
	bgIntenseCyan    AnsiCode = codePrefix + "106" + codeSuffix
	bgIntenseWhite   AnsiCode = codePrefix + "107" + codeSuffix
)
