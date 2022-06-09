package console

// ANSI colour codes

type ansiCode string

const (
	escape     ansiCode = "\u001b"
	csi        ansiCode = "["
	codePrefix          = escape + csi
	codeSuffix          = "m"
)

const (
	reset ansiCode = codePrefix + "0" + codeSuffix

	bold          ansiCode = codePrefix + "1" + codeSuffix
	faint         ansiCode = codePrefix + "2" + codeSuffix
	italic        ansiCode = codePrefix + "3" + codeSuffix
	underline     ansiCode = codePrefix + "4" + codeSuffix
	invert        ansiCode = codePrefix + "7" + codeSuffix
	conceal       ansiCode = codePrefix + "8" + codeSuffix
	strikethrough ansiCode = codePrefix + "9" + codeSuffix

	fgBlack   ansiCode = codePrefix + "30" + codeSuffix
	fgRed     ansiCode = codePrefix + "31" + codeSuffix
	fgGreen   ansiCode = codePrefix + "32" + codeSuffix
	fgYellow  ansiCode = codePrefix + "33" + codeSuffix
	fgBlue    ansiCode = codePrefix + "34" + codeSuffix
	fgMagenta ansiCode = codePrefix + "35" + codeSuffix
	fgCyan    ansiCode = codePrefix + "36" + codeSuffix
	fgWhite   ansiCode = codePrefix + "37" + codeSuffix
	fgDefault ansiCode = codePrefix + "39" + codeSuffix

	bgBlack   ansiCode = codePrefix + "40" + codeSuffix
	bgRed     ansiCode = codePrefix + "41" + codeSuffix
	bgGreen   ansiCode = codePrefix + "42" + codeSuffix
	bgYellow  ansiCode = codePrefix + "43" + codeSuffix
	bgBlue    ansiCode = codePrefix + "44" + codeSuffix
	bgMagenta ansiCode = codePrefix + "45" + codeSuffix
	bgCyan    ansiCode = codePrefix + "46" + codeSuffix
	bgWhite   ansiCode = codePrefix + "47" + codeSuffix

	fgIntenseBlack   ansiCode = codePrefix + "90" + codeSuffix
	fgIntenseRed     ansiCode = codePrefix + "91" + codeSuffix
	fgIntenseGreen   ansiCode = codePrefix + "92" + codeSuffix
	fgIntenseYellow  ansiCode = codePrefix + "93" + codeSuffix
	fgIntenseBlue    ansiCode = codePrefix + "94" + codeSuffix
	fgIntenseMagenta ansiCode = codePrefix + "95" + codeSuffix
	fgIntenseCyan    ansiCode = codePrefix + "96" + codeSuffix
	fgIntenseWhite   ansiCode = codePrefix + "97" + codeSuffix

	bgIntenseBlack   ansiCode = codePrefix + "100" + codeSuffix
	bgIntenseRed     ansiCode = codePrefix + "101" + codeSuffix
	bgIntenseGreen   ansiCode = codePrefix + "102" + codeSuffix
	bgIntenseYellow  ansiCode = codePrefix + "103" + codeSuffix
	bgIntenseBlue    ansiCode = codePrefix + "104" + codeSuffix
	bgIntenseMagenta ansiCode = codePrefix + "105" + codeSuffix
	bgIntenseCyan    ansiCode = codePrefix + "106" + codeSuffix
	bgIntenseWhite   ansiCode = codePrefix + "107" + codeSuffix
)
