package entities

import (
	"encoding/json"
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"math"
	"time"
)

type Planet struct {
	Color               Color
	Coords              []int
	IsHome              bool
	Texture             int
	Size                int
	LastShipCountUpdate int64
	ShipCount           int
	MaxShipCount        int
	Owner               string
}

type marshalHook Planet

func (p *Planet) String() string {
	return fmt.Sprintf("Planet[%s, %s]", p.Coords[0], p.Coords[1])
}

func (p *Planet) GetKey() string {
	return fmt.Sprintf("planet.%d_%d", p.Coords[0], p.Coords[1])
}

func (p *Planet) HasOwner() bool {
	return len(p.Owner) > 0
}

func (p *Planet) MarshalJSON() ([]byte, error) {
	p.UpdateShipCount()
	return json.Marshal((*planetMarshalHook)(p))
}

func (p *Planet) GetShipCount() int {
	p.UpdateShipCount()
	return p.ShipCount
}

func (p *Planet) SetShipCount(count int) {
	p.UpdateShipCount()
	p.ShipCount = count
	p.LastShipCountUpdate = time.Now().Unix()
}

func (p *Planet) UpdateShipCount() {
	if p.HasOwner() {
		passedTime := time.Now().Unix() - p.LastShipCountUpdate
		timeModifier := int64(p.Size/3) + 1
		//TODO: To be completed for all planet size types
		//if getobject(Owner.getkey).gethomeplanet == p.getkey
		p.ShipCount += int(passedTime / (timeModifier * 10))
		p.LastShipCountUpdate = time.Now().Unix()
	}
}

/*
TODO: We need to add ship count on new planet creation
TODO: Put all funny numbers in a constans in our config file
NOTE: 5 in ring_offset is the distance between planets
*/
func GeneratePlanets(hash string, sun_position *vec2d.Vector) ([]*Planet, *Planet) {
	hashElement := func(index int) float64 {
		return float64(hash[index]) - 48 // The offset of simbol "1" in the ascii table
	}

	result := []*Planet{}
	ring_offset := float64(PLANETS_RING_OFFSET)
	planet_radius := float64(PLANETS_PLANET_RADIUS)

	for ix := 0; ix < PLANETS_PLANET_COUNT; ix++ {
		planet_in_creation := Planet{Color{200, 180, 140}, []int{0, 0}, false, 0, 0, time.Now().Unix(), 10, 0, ""}
		ring_offset += planet_radius + hashElement(4*ix)*5

		planet_in_creation.Coords[0] = int(float64(sun_position.X) + ring_offset*math.Cos(
			hashElement(4*ix+1)*40))
		planet_in_creation.Coords[1] = int(float64(sun_position.Y) + ring_offset*math.Sin(
			hashElement(4*ix+1)*40))

		planet_in_creation.Texture = int(hashElement(4*ix + 2))
		planet_in_creation.Size = 1 + int(hashElement(4*ix+3))
		planet_in_creation.LastShipCountUpdate = time.Now().Unix()
		result = append(result, &planet_in_creation)
	}
	// + 1 bellow stands for: after all the planet info is read the next element is the user's home planet idx
	homePlanetIdx := int(hashElement(PLANETS_PLANET_COUNT*PLANETS_PLANET_HASH_ARGS + 1))
	result[homePlanetIdx].IsHome = true
	return result, result[homePlanetIdx]
}
