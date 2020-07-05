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

func removeEnt(ents []Id, i int) []Id {
	ents[i] = ents[len(ents)-1]
	return ents[:len(ents)-1]
}

//
// TODO: move/sort these
//
var (
	list1 = []string{
		"star",
		"planet",
		"moon",
		"ice",
		"rock",
		"ring",
		"steel",
		"free",
		"easy",
		"range",
		"way",
		"halo",
		"light",
		"dark",
		"night",
		"fall",
		"dawn",
		"dusk",
		"day",
		"first",
		"head",
		"arch",
		"fore",
		"dare",
		"war",
		"clash",
		"tri",
		"blast",
		"up",
		"down",
		"next",
		"wild",
		"dust",
	}
	list2 = []string{
		"finder",
		"bender",
		"seeker",
		"catcher",
		"rover",
		"drifter",
		"booter",
		"nomad",
		"grim",
		"grimer",
		"farer",
		"roamer",
		"cutter",
		"ray",
		"valier",
		"guard",
		"guarder",
		"ment",
		"able",
		"fall",
		"diver",
		"buckler",
		"basher",
		"taker",
		"ing",
		"ling",
		"ler",
		"blaze",
		"blink",
		"link",
		"beam",
		"flare",
		"burn",
		"warper",
		"siren",
		"edge",
		"mark",
		"draw",
		//"",
	}
)

func nameGen() []string {
	count := 36
	names := []string{}

	for i := 0; i < count; i++ {
		prefix := list1[Rand.Intn(len(list1))]
		suffix := list2[Rand.Intn(len(list2))]
		names = append(names, prefix+""+suffix)
	}

	return names
}
