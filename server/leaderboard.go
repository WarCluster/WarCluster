package server

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"warcluster/entities"
	"warcluster/leaderboard"
)

type searchResult struct {
	Username string
	Page     int
}

func leaderboardPlayersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	pageQuery, ok := r.URL.Query()["page"]
	if !ok {
		http.Error(w, "Bad Request", 400)
		return
	}

	page, intErr := strconv.ParseInt(pageQuery[0], 10, 0)
	if intErr != nil {
		http.Error(w, "Page Not Found", 404)
		return
	}

	boardPage, err := leaderBoard.Page(page)
	if err != nil {
		http.Error(w, "Page Not Found", 404)
		return
	}

	result, _ := json.Marshal(boardPage)
	fmt.Fprintf(w, string(result))
}

func leaderboardRacesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	races, err := json.Marshal(leaderBoard.Races())
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	fmt.Fprintf(w, string(races))
}

func leaderboardRacesInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	races, err := json.Marshal(entities.Races)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	fmt.Fprintf(w, string(races))
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	username, ok := r.URL.Query()["player"]
	if !ok || len(username[0]) < 3 {
		http.Error(w, "Bad Request", 400)
		return
	}

	players, err := entities.GetList(fmt.Sprintf("player.%s*", username[0]))
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	result := make([]searchResult, 0)

	for _, player := range players {
		username := player[7:]
		page := math.Ceil(float64(leaderBoard.Place(username)+1) / 10)
		result = append(result, searchResult{username, int(page)})
	}

	marshalledResult, _ := json.Marshal(result)
	fmt.Fprintf(w, string(marshalledResult))
}

// Initialize the leaderboard
func InitLeaderboard(board *leaderboard.Leaderboard) {
	log.Println("Initializing the leaderboard...")
	allPlayers := make(map[string]*leaderboard.Player)
	playerEntities := entities.Find("player.*")
	planetEntities := entities.Find("planet.*")

	for key, value := range cfg.Race {
		board.AddRace(
			key,
			value.Id,
		)
	}

	for _, playerEntity := range playerEntities {
		player := playerEntity.(*entities.Player)

		leaderboardPlayer := &leaderboard.Player{
			Username: player.Username,
			RaceId:   player.RaceID,
			Planets:  0,
		}
		allPlayers[player.Username] = leaderboardPlayer
		board.Add(leaderboardPlayer)
	}

	for _, entity := range planetEntities {
		planet, ok := entity.(*entities.Planet)

		if !planet.HasOwner() || !ok {
			continue
		}

		player, _ := allPlayers[planet.Owner]

		if planet.IsHome {
			player.HomePlanet = planet.Name
		}

		player.Planets++
	}
	board.Sort()
	board.RecountRacesPlanets()
	leaderBoard = board
}
