package entities

import (
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"github.com/Vladimiroff/vec2d"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Color struct {
	R int
	G int
	B int
}

func (c *Color) String() string {
	return fmt.Sprintf("Color[%s, %s, %s]", c.R, c.G, c.B)
}

func Construct(key string, data []byte) Entity {
	entityType := strings.Split(key, ".")[0]

	switch entityType {
	case "player":
		player := new(Player)
		json.Unmarshal(data, player)
		player.username = strings.Split(key, "player.")[1]
		return player
	case "planet":
		planet := new(Planet)
		json.Unmarshal(data, planet)
		return planet
	case "mission":
		mission := new(Mission)
		json.Unmarshal(data, mission)
		return mission
	case "sun":
		sun := new(Sun)
		json.Unmarshal(data, sun)
		sun.position = ExtractSunKey(key)
		return sun
	}
	return nil
}

func generateHash(username string) string {
	return simplifyHash(usernameHash(username))
}

func ExtractSunKey(key string) *vec2d.Vector {
	paramsRaw := strings.Split(key, ".")[1]
	params := strings.Split(paramsRaw, "_")
	sunCoords_0, _ := strconv.ParseFloat(params[0], 64)
	sunCoords_1, _ := strconv.ParseFloat(params[1], 64)
	coords := vec2d.New(sunCoords_0, sunCoords_1)
	return coords
}

func usernameHash(username string) []byte {
	hash := sha512.New()
	io.WriteString(hash, username)
	return hash.Sum(nil)
}

func simplifyHash(hash []byte) string {
	result := ""
	for ix := 0; ix < len(hash); ix++ {
		lastDigit := hash[ix] % 10
		result += strconv.Itoa(int(lastDigit))
	}
	return result
}

func getRandomStartPosition(scope int) *vec2d.Vector {
	xSeed := time.Now().UTC().UnixNano()
	ySeed := time.Now().UTC().UnixNano()
	xGenerator := rand.New(rand.NewSource(xSeed))
	yGenerator := rand.New(rand.NewSource(ySeed))
	return vec2d.New(float64(xGenerator.Intn(scope)-scope/2), float64(yGenerator.Intn(scope)-scope/2))
}
