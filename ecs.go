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

// If Gustav had invented the Entity-Component-System (ECS) architectural
// pattern, he would have named it Id-Data-Function (IDF), as in ECS
// entities are ids, components are data and systems are functions.

// Entities are ids.
type Id uint64

// Components are data containers for entities.  They contain no game logic.
type Component interface {
	// TODO
}

/* Systems contain game logic.  They operate on and manipulate component data.
   Each system handles a subset of components, which may overlap with other
   component subsets handled by other systems.

   Examples of systems include physics, combat and player trading exchanges.

   While components are the primary data containers, systems track some data
   such as which entities "are in the system"
   (have data in all the system's components)
*/
type System interface {
	Init() error
	Update(elapsed float64, rf *RefFrame, entMap map[Id]bool) error
}
