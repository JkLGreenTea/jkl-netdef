package translit

import (
	"bytes"
	"golang.org/x/exp/utf8string"
	"strings"
	"unicode"
)

var RuTransiltMap = map[rune]string{
	'а': "a",
	'б': "b",
	'в': "v",
	'г': "g",
	'д': "d",
	'е': "e",
	'ё': "yo",
	'ж': "zh",
	'з': "z",
	'и': "i",
	'й': "j",
	'к': "k",
	'л': "l",
	'м': "m",
	'н': "n",
	'о': "o",
	'п': "p",
	'р': "r",
	'с': "s",
	'т': "t",
	'у': "u",
	'ф': "f",
	'х': "h",
	'ц': "c",
	'ч': "ch",
	'ш': "sh",
	'щ': "sch",
	'ъ': "'",
	'ы': "y",
	'ь': "",
	'э': "e",
	'ю': "ju",
	'я': "ya",
}

func Transliterate(text string, translitMap map[rune]string) string {
	var result bytes.Buffer
	utf8text := utf8string.NewString(text)
	length := utf8text.RuneCount()
	for index := 0; index < length; index++ {
		runeValue := utf8text.At(index)
		switch str, ok := translitMap[unicode.ToLower(runeValue)]; {
		case !ok:
			result.WriteRune(runeValue)
		case str == "":
			continue
		case unicode.IsUpper(runeValue):

			if (length > index+1 && unicode.IsUpper(utf8text.At(index+1))) ||
				(index > 0 && unicode.IsUpper(utf8text.At(index-1))) {
				str = strings.ToUpper(str)
			} else {
				str = strings.Title(str)
			}
			fallthrough
		default:
			result.WriteString(str)
		}
	}
	return result.String()
}

// Ru - выполняет транслитерацию строки с учетом словаря для русской транслитерации.
func Ru(text string) string {
	return Transliterate(text, RuTransiltMap)
}
