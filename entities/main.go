package entities

import (
	"encoding/json"
	"fmt"
	"strings"
	"warcluster/entities/db"
)

const (
	ENTITIES_RANGE_SIZE           = 10000
	PLANETS_RING_OFFSET           = 300
	PLANETS_PLANET_RADIUS         = 300
	PLANETS_PLANET_COUNT          = 10
	PLANETS_PLANET_HASH_ARGS      = 4
	SUNS_RANDOM_SPAWN_ZONE_RADIUS = 50000
	SUNS_SOLAR_SYSTEM_RADIUS      = 9000
)

// Entity interface is implemented by all entity types here
type Entity interface {
	SortedSet(string) (string, float64)
	Key() string
}

// Simple RGB color struct
type Color struct {
	R uint8
	G uint8
	B uint8
}

// Finds records in the database, by given key
// All Redis wildcards are allowed.
func Find(query string) []Entity {
	var entityList []Entity

	if records, err := db.GetList(query); err == nil {
		results := fmt.Sprintf("%s", records)
		for _, key := range strings.Split(results[1:len(results)-1], " ") {
			if entity, err := Get(key); err == nil {
				entityList = append(entityList, entity)
			}
		}
	}

	return entityList
}

// Fetches a single record in the database, by given concrete key.
// If there is no entity with such key, returns error.
func Get(key string) (Entity, error) {
	record, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	return Construct(key, record), nil
}

// Saves an entity to the database. Records' key is entity.Key()
// If there is a record with such key in the database, simply updates
// the record. Otherwise creates a new one.
//
// Failed marshaling of the given entity is pretty much the only
// point of failure in this function... I supose.
func Save(entity Entity) error {
	key := entity.Key()
	value, err := json.Marshal(entity)
	if err != nil {
		return err
	}
	err = db.Save(key, value)
	xSet, xWeight := entity.SortedSet("X")
	ySet, yWeight := entity.SortedSet("Y")
	db.Zadd(xSet, xWeight, key)
	db.Zadd(ySet, yWeight, key)
	return err
}

// Deletes a record by the given key
func Delete(key string) error {
	return db.Delete(key)
}
