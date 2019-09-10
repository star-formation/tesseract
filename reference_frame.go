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
)

/*  See https://en.wikipedia.org/wiki/Frame_of_reference
    and https://en.wikipedia.org/wiki/Celestial_coordinate_system

    The game world can be seen as a hiearchical tree of reference frames,
    each frame except the top-level one having a parent reference frame.

    Child frames are "dragged along" their parent frames.

    Example:
    - a space station orbits a star
    - the station's orbit is defined within the star's reference frame
    - the station in turn "hosts" a reference frame, centered on the station
    - a player starship is at standstill next to the station
    - the ship is at 0 velocity within the local frame
    - the ship _also_ "inherits" the station's orbit from the outer frame

    Most gameplay takes place within the context of the local reference frame,
    but some aspects can cross frames.  In the example above; if the player ship
    is at standstill and bathing in starlight, it may end up in the shadow of
    the station after some time, if the station (and its frame) is not rotating
    with respect to the star's reference frame.
    - then, the ship may get much less power from any deployed solar panels,
      even though nothing changed in the local frame.
    - the game logic of solar panels and starlight operates across
      two (or more!) reference frames.

    The top-level or root reference frame is the Milky Way galaxy.
    It has no parent or surrounding context, and can be thought of as a
    static/stationary 3D grid.

    2nd level reference frames are generally star systems, centered on the
    system's approximate barycenter.  In a single-star system the reference
    frame center is the center of the star.

    3rd level frames can be planets orbiting stars, 4th level moons of planets.

    There is no limit on the depth of reference frames, but in practice the max
    depth/level is probably 6 or 7.
    (inside structure on asteroid orbiting moon orbiting planet orbiting star)
*/
type RefFrame struct {
    // The top-level reference frame (the Milky Way galaxy) has Parent,
    // Position, Orbit and Orientation all set to nil.
    Parent *RefFrame

    /* The location of a reference frame relative its parent is encoded
       as either a 3D X,Y,Z position or as orbital elements, with the other
       set as nil.

       If the frame is stationary relative its parent - for example the inside
       of a building on a planet surface - then it has a 3D position (X,Y,Z)
       but no orbital elements.

       For now, the only allowed movement of a reference frame relative
       its parent is orbits.  An example of this is the grid surrounding
       a space station orbiting a planet.  Such a frame has 3D orbital
       elements (E,S,I,L,A,T) but no position (X,Y,Z) set.

       X,Y,Z coordinates can always be derived from E,S,I,L,A,T
       (see orbit.go) but to translate carteesian coordinates to orbital
       elements we need both position and velocity.
    */
    Position *V3
    Orbit *OE
    
    // Except for the top-level frame, Orientation is always non-nil;
    // as it's required to translate local coordinates to outer frame(s).
    // A zero (0,0,0) orientation equals inheriting parent frame orientation
    Orientation *V3

    // TODO: size/shape

    // Rotation is unsupported for now.
}
