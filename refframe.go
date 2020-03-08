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

/*  See https://en.wikipedia.org/wiki/Frame_of_reference
    and https://en.wikipedia.org/wiki/Celestial_coordinate_system

    The game world is a hiearchical tree of reference frames.
    Each reference frame except the root has one parent reference frame.
    Child frames are "dragged along" their parent frames.

    Most gameplay logic only knows about the local frame, but some gameplay
    involves multiple frames.  For example, warp drive involves leaving one
    ref frame, spending some time in a noninteractive / locked frame and
    then arriving in a destination frame.

    The top-level or root reference frame is the Milky Way galaxy.
    It has no parent or surrounding context and can be thought of as a
    static/stationary 3D grid.

    2nd level reference frames are generally star systems, centered on the
    system's approximate barycenter.  In a single-star system the reference
    frame center is the center of the star.

    3rd level frames can be planets orbiting stars, 4th level moons of planets.

    The location of a reference frame relative its parent is encoded
    as either a 3D X,Y,Z position or as orbital elements, with the other
    set to nil.

    If the frame is stationary relative its parent - for example the inside
    of a building on a planet surface - then it has a 3D position (X,Y,Z)
    but no orbital elements.
*/
type RefFrame struct {
	// The top-level reference frame (the Milky Way galaxy) has Parent,
	// Position, Orbit and Orientation all set to nil.
	Parent *RefFrame

	Pos   *V3
	Orbit *OE

	// Except for the top-level frame, Orientation is always non-nil;
	// it's required to translate local coordinates to outer frame(s).
	// A zero orientation equals inheriting the parent's frame orientation
	Orientation *Q

	Radius float64

	// Rotation is unsupported for now.
	// DragCoef1, DragCoef2 float64
}

func (rf *RefFrame) IsRoot() bool {
	return rf.Parent == nil
}
