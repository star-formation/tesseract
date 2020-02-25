/*  Copyright 2019 The tesseract Authors

    This file is part of tesseract.

    tesseract is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    tesseract is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package tesseract

import (
	"fmt"
	"unicode"
)

// TODO: check if we can use package unicode for more of this
type Script struct {
	Name         string
	UnicodeStart int
	CharCount    int
	Weight       int
}

var (
	// https://en.wikipedia.org/wiki/List_of_writing_systems#List_of_writing_scripts_by_adoption
	scripts = []*Script{
		&Script{"Latin", 65, 26, wLatin},
		// https://stackoverflow.com/questions/1366068/whats-the-complete-range-for-chinese-characters-in-unicode
		&Script{"Han", 0x4E00, 0x9FFF - 0x4E00, wChinese},
		&Script{"Devanagari", 0x0900, 0x097F - 0x0900, wDevanagari},
		&Script{"Arabic", 0x0600, 0x06FF - 0x0600, wArabic},
		// TODO: add Tirhuta?
		&Script{"Bengali", 0x0980, 0x09FF - 0x0980, wBengaliAssamese},
		&Script{"Cyrillic", 0x0400, 0x04FF - 0x0400, wCyrillic},
		&Script{"Katakana", 0x30A0, 0x30FF - 0x30A0, wKana},
		&Script{"Javanese", 0xA980, 0xA9DF - 0xA980, wJavanese},
		&Script{"Hangul", 0x1100, 0x11FF - 0x1100, wHangul},
	}

	numerals = map[string]map[int]rune{
		"Latin": map[int]rune{
			0: '0',
			1: '1',
			2: '2',
			3: '3',
			4: '4',
			5: '5',
			6: '6',
			7: '7',
			8: '8',
			9: '9',
		},
		// https://en.wikipedia.org/wiki/Chinese_numerals#Characters_with_military_usage
		"Han": map[int]rune{
			0: '洞',
			1: '幺',
			2: '两',
			3: '三',
			4: '四',
			5: '五',
			6: '六',
			7: '拐',
			8: '八',
			9: '勾',
		},
		"Devanagari": map[int]rune{
			0: '०',
			1: '१',
			2: '२',
			3: '३',
			4: '४',
			5: '५',
			6: '६',
			7: '७',
			8: '८',
			9: '९',
		},
		"EasternArabic": map[int]rune{
			0: '٠',
			1: '١',
			2: '٢',
			3: '٣',
			4: '٤',
			5: '٥',
			6: '٦',
			7: '٧',
			8: '٨',
			9: '٩',
		},
		"Bengali": map[int]rune{
			0: '০',
			1: '১',
			2: '২',
			3: '৩',
			4: '৪',
			5: '৫',
			6: '৬',
			7: '৭',
			8: '৮',
			9: '৯',
		},
		"Japanese": map[int]rune{
			0: '〇',
			1: '一',
			2: '二',
			3: '三',
			4: '四',
			5: '五',
			6: '六',
			7: '七',
			8: '八',
			9: '九',
		},
		"Javanese": map[int]rune{
			0: '꧐',
			1: '꧑',
			2: '꧒',
			3: '꧓',
			4: '꧔',
			5: '꧕',
			6: '꧖',
			7: '꧗',
			8: '꧘',
			9: '꧙',
		},
	}
)

func starName() string {
	s := randScript()
	//fmt.Printf("%s: ", s.Name)
	letters := randLetters(s)
	//fmt.Printf("letters: %v\n", letters)
	var d0, d1 rune
	var nMap map[int]rune
	switch s.Name {
	case "Han", "Devanagari", "Bengali", "Javanese":
		nMap = numerals[s.Name]
	case "Katakana":
		nMap = numerals["Japanese"]
	case "Arabic":
		if Rand.Float32() < 0.33 {
			nMap = numerals["EasternArabic"]
		} else {
			nMap = numerals["Latin"]
		}
	default:
		nMap = numerals["Latin"]
	}
	d0, d1 = nMap[Rand.Intn(10)], nMap[Rand.Intn(10)]
	runes := append(letters, []rune{'-', d0, d1}...)
	return string(runes)
}

func randScript() *Script {
	wSum := 0
	for _, s := range scripts {
		wSum += s.Weight
	}

	x := Rand.Intn(wSum)
	for _, s := range scripts {
		if x < s.Weight {
			return s
		}
		x -= s.Weight
	}
	return nil
}

func randLetters(s *Script) []rune {
	return []rune{
		randLetter(s),
		randLetter(s),
		randLetter(s),
	}
}

func randLetter(s *Script) rune {
	r := rune(s.UnicodeStart + Rand.Intn(s.CharCount))
	rt, ok := unicode.Scripts[s.Name]
	if !ok {
		panic(fmt.Errorf("unicode script name not found: %s", s.Name))
	}
	if unicode.Is(rt, r) {
		return r
	} else {
		return randLetter(s)
	}
}
