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

// Ship classes are analogous to classical navy ship classes
// (https://en.wikipedia.org/wiki/Ship_class) and are uniquely identified
// by a set of constants.
//
// Note that ship classes are not OOP classes and there is no inheritance or
// instancing of ship classes.  Rather, each ship class is encoded as
// a set of constants and an empty struct implements the ShipClass interface.
//
// This enables other packages to read immutable ship params while
// game state components store dynamic parameters like ship mass and
// hull/armor/shield hit points.
//
// As per the naming convention, "Base" and "Cap" denotes a value intrinsic
// to the ship class - before effects are applied from any attached modules
// or player character skills.  "Base" and "Cap" are always suffixes.
//
// Note that a parameter's real-time value can be lower than its base value.
// An example of this is a ship without any modules or cargo flown by a player
// with trained aerodynamic skills - the ship's effective drag would be lower
// than the base hull value.
type ShipClass interface {
	// The mass in kilograms (kg) of the ship when no modules are attached
	// and nothing is in the cargo bay.
	MassBase() float64

	// Volume in cubic meters (m^3) of ship excluding any external modules.
	VolumeBase() float64
	// The volume of the ship when packed inside a cargo bay or other storage.
	PackedVolumeBase() float64

	// Control Moment Gyroscope Max Torque in newtons (N).
	//
	// This is a built-in, non-modular engine situated at the ship's
	// center of mass.  The CMG has a single function: generate torque around
	// the ship's center of mass.  This rotates the ship in any direction
	// without changing velocity.
	//
	// This is a simple mechanism to enable the physics engine to simulate
	// "turning inertia" realistically - players will notice their ships
	// turn slower when storing heavy cargo or fielding heavy modules
	// like armor plates.
	//
	// By using a 3D vector we can configure ships that turn faster "up"/"down"
	// (pitch) than "left"/"right" (yaw), making banked turns useful
	// even in zero-g / vacuum (https://en.wikipedia.org/wiki/Banked_turn).
	CMGTorqueCap() V3

	// Hull/Armor/Shield Capacity in hit points.
	HullHPCap() float64
	//ArmorHPCap float64
	//ShieldHPCap float64

	// Cargo bay capacity in cubic meters (m^3).
	CargoBayCap() float64

	// Hull Aerodynamic Lift and Drag Coefficients (dimensionless quantity).
	// This is the base lift/drag of the ship hull prior to taking into account
	// modules that affect lift/drag and aerodynamic player skills.
	// TODO: split into subsonic, supersonic, hypersonic
	AeroLiftBase() float64
	AeroDragBase() float64

	// External hard points and internal module slots are the main mechanism
	// for ship modularity; the exact counts of each type of module slot
	// is a crucial part of the balance for each ship class.

	// Hard Point integer count.
	// Each Hard Point can attach one external module such as a weapon module.
	HardPoints() uint8

	// High Power Slot integer count.
	// Each internal high power slot can fit one high-powered module.
	HighPowerSlots() uint8

	// Low Power Slot integer count.
	// Each internal low power slot can fit one low-powered module.
	LowPowerSlots() uint8

	// Each ship class also has one engine module slot, one reactor slot and
	// one mainframe slot.  For now this is hard-coded in the gameplay/ship
	// system of the game engine as those module slot counts are singular.

	// TODO: ship bonuses
}

type Engine interface {
	MaxThrust() float64
	SetThrust(float64) error
}

type WarmJet struct{}

const (
	massBaseWarmjet = 42000 // kg

	volumeBaseWarmjet       = 60 // m3
	packedVolumeBaseWarmjet = 60 // m3

	cmgTorqueCapWarmjetX = 1000000.0 // Newton (N)
	cmgTorqueCapWarmjetY = 1000000.0
	cmgTorqueCapWarmjetZ = 1000000.0

	hullHPCapWarmjet = 100 // hit points

	cargoBayCapWarmjet = 10 // m3

	aeroLiftBaseWarmjet = 0.2 // dimensionless coefficient
	aeroDragBaseWarmjet = 0.2

	hardPointsWarmjet     = 2 // integer count
	highPowerSlotsWarmjet = 1
	lowPowerSlotsWarmjet  = 1

	// TODO: add shape / size
)

func (s *WarmJet) MassBase() float64         { return massBaseWarmjet }
func (s *WarmJet) VolumeBase() float64       { return volumeBaseWarmjet }
func (s *WarmJet) PackedVolumeBase() float64 { return packedVolumeBaseWarmjet }
func (s *WarmJet) CMGTorqueCap() V3 {
	return V3{cmgTorqueCapWarmjetX, cmgTorqueCapWarmjetY, cmgTorqueCapWarmjetZ}
}
func (s *WarmJet) HullHPCap() float64    { return hullHPCapWarmjet }
func (s *WarmJet) CargoBayCap() float64  { return cargoBayCapWarmjet }
func (s *WarmJet) AeroLiftBase() float64 { return aeroLiftBaseWarmjet }
func (s *WarmJet) AeroDragBase() float64 { return aeroDragBaseWarmjet }
func (s *WarmJet) HardPoints() uint8     { return hardPointsWarmjet }
func (s *WarmJet) HighPowerSlots() uint8 { return highPowerSlotsWarmjet }
func (s *WarmJet) LowPowerSlots() uint8  { return lowPowerSlotsWarmjet }
