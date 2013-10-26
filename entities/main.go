package entities

import (
	"encoding/json"
	"warcluster/entities/db"
)

const (
	PLANETS_RING_OFFSET           = 300
	PLANETS_PLANET_RADIUS         = 300
	PLANETS_PLANET_COUNT          = 10
	PLANETS_PLANET_HASH_ARGS      = 4
	SUNS_RANDOM_SPAWN_ZONE_RADIUS = 50000
	SUNS_SOLAR_SYSTEM_RADIUS      = 9000
)

type Entity interface {
	GetKey() string
	String() string
}

type missionMarshalHook Mission
type planetMarshalHook Planet

// Finds records in the database, by given key
// All Redis wildcards are allowed.
func Find(query string) []Entity {
}

// Fetches a single record in the database, by given concrete key.
// If there is no entity with such key, returns error.
func Get(key string) (Entity, error) {
	if record, err := db.Get(key); err != nil {
		return nil, err
	}
	return Construct(key, record), nil
}

// Saves an entity to the database. Records' key is entity.GetKey()
// If there is a record with such key in the database, simply updates
// the record. Otherwise creates a new one.
//
// Failed marshaling of the given entity is pretty much the only
// point of failure in this function... I supose.
func Save(entity Entity) error {
	key := entity.GetKey()
	value, err := json.Marshal(entity); if err != nil {
		return err
	}
	return db.Save(key, value)
}

// Deletes a record by the given key
func Delete(key string) error {
	return db.Delete(key)
}