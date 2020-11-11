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
	"math"
)

//
// Collision Detection
//
type BoundingShape uint8

const (
	Sphere = iota
	Box
)

// The BoundingVolume interface enables any 3D volume that fully bounds
// one or more 3D volumes to be used with the Bounding Volume Hierarchy Tree
// and other constructs used in collision detection.
type BoundingVolume interface {
	// Returns the shape type of the bounding volume.
	Shape() BoundingShape

	// Returns whether this volume overlaps with the passed volume.
	Overlaps(BoundingVolume) bool

	// Returns a bounding volume that fully bounds this volume
	// and the passed volume.
	NewBoundingVolume(BoundingVolume) BoundingVolume

	// Returns how much this bounding volume would have to grow in order to
	// bound the passed bounding volume.
	// This value can be derived from the NewBoundingVolume method and is
	// provided a separate method for optimization purposes.
	// The returned value is dimensionless and not a volume unit.
	CalcGrowth(BoundingVolume) float64

	// Returns the volume in cubic meters (m^3)
	Volume() float64

	// Returns the surface area in square meters (m^2)
	SurfaceArea() float64
}

type BoundingSphere struct {
	P *V3
	R float64
}

// BoundingVolume interface
func (s *BoundingSphere) Shape() BoundingShape {
	return Sphere
}

func (s *BoundingSphere) Overlaps(bv BoundingVolume) bool {
	if bv.Shape() != Sphere {
		panic("unsupported bounding volume shape")
	}

	s2 := bv.(*BoundingSphere)
	return new(V3).Sub(s.P, s2.P).SquareMagnitude() < (s.R+s2.R)*(s.R+s2.R)
}

func (s *BoundingSphere) NewBoundingVolume(bv BoundingVolume) BoundingVolume {
	if bv.Shape() != Sphere {
		panic("unsupported bounding volume shape")
	}

	s2 := bv.(*BoundingSphere)
	radiusDelta := s2.R - s.R
	posDelta := new(V3).Sub(s2.P, s.P)
	distance := posDelta.Magnitude()

	// return copy of either sphere if it fully encloses the other
	if math.Abs(radiusDelta) >= distance {
		if s.R > s2.R {
			return &BoundingSphere{new(V3).Set(s.P), s.R}
		} else {
			return &BoundingSphere{new(V3).Set(s2.P), s2.R}
		}
	} else {
		// overlapping spheres; create new sphere enclosing both
		newR := (s.R + s2.R + distance) * 0.5
		newP := new(V3).Set(s.P)
		// adjust position of new sphere towards s2.P
		if distance > 0 {
			// newP += posDelta * ((newR - s.R) / distance)
			newP.AddScaledVector(posDelta, ((newR - s.R) / distance))
		}
		return &BoundingSphere{newP, newR}
	}
}

func (s *BoundingSphere) CalcGrowth(bv BoundingVolume) float64 {
	if bv.Shape() != Sphere {
		panic("unsupported bounding volume shape")
	}

	return bv.SurfaceArea() - s.SurfaceArea()
}

func (s *BoundingSphere) Volume() float64 {
	return (4.0 / 3.0) * math.Pi * s.R * s.R * s.R
}

func (s *BoundingSphere) SurfaceArea() float64 {
	return 4.0 * math.Pi * s.R * s.R
}

// Bounding Volume Hierarchy Tree.
// Each non-leaf holds a bounding volume encompassing all its child nodes.
// Each leaf     holds a bounding volume of a single entity.
type BVHNode struct {
	parent, left, right *BVHNode
	entity              Id // nil for non-leaf nodes
	volume              BoundingVolume
}

func (n *BVHNode) IsLeaf() bool {
	return n.entity != 0
}

func (n *BVHNode) Insert(e Id, v BoundingVolume) {
	if n.IsLeaf() {
		n.left = &BVHNode{n, nil, nil, n.entity, n.volume}
		n.right = &BVHNode{n, nil, nil, e, v}
		n.entity = 0
		n.UpdateBoundingVolume()
	} else {
		// TODO: handle when left or right is empty
		if n.left.volume.CalcGrowth(v) < n.right.volume.CalcGrowth(v) {
			n.left.Insert(e, v)
		} else {
			n.right.Insert(e, v)
		}
	}
}

func (n *BVHNode) Delete() {
	if n.parent != nil {
		var sibling *BVHNode
		if n.parent.left == n {
			sibling = n.parent.right
		} else {
			sibling = n.parent.left
		}
		n.parent.left = sibling.left
		n.parent.right = sibling.right
		n.parent.entity = sibling.entity
		n.parent.volume = sibling.volume
		n.parent.UpdateBoundingVolume()
	}
	n.left = nil
	n.right = nil
}

func (n *BVHNode) UpdateBoundingVolume() {
	n.volume = n.left.volume.NewBoundingVolume(n.right.volume)
	if n.parent != nil {
		n.parent.UpdateBoundingVolume()
	}
}

func (n *BVHNode) PotentialContacts() [][2]Id {
	if n.IsLeaf() {
		return nil
	}

	contacts := make([][2]Id, 0)
	// recursively descend into child nodes, appending contacts
	potentialContactsWith(n.left, n.right, &contacts)
	return contacts
}

func potentialContactsWith(n1, n2 *BVHNode, contacts *[][2]Id) {
	if !n1.volume.Overlaps(n2.volume) {
		return
	}

	if n1.IsLeaf() && n2.IsLeaf() {
		*contacts = append(*contacts, [2]Id{n1.entity, n2.entity})
		return
	}

	// If either node is a leaf, then we descend the other.
	// If both nodes are branches, then we descent the node with larger volume.
	if n2.IsLeaf() || (!n1.IsLeaf() && n1.volume.Volume() > n2.volume.Volume()) {
		potentialContactsWith(n1.left, n2, contacts)
		potentialContactsWith(n1.right, n2, contacts)
	} else {
		potentialContactsWith(n2.left, n1, contacts)
		potentialContactsWith(n2.right, n1, contacts)
	}
}
