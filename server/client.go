package server

import (
	"../db_manager"
	"../entities"
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type Client struct {
	conn     net.Conn
	nickname string
	channel  chan string
	player   *entities.Player
}

func authenticate(c net.Conn, bufc *bufio.Reader) (string, *entities.Player) {
	var player entities.Player

	io.WriteString(c, "Twitter Authenticating:\n")
	io.WriteString(c, "Username: ")
	nick, _, _ := bufc.ReadLine()
	nickname := string(nick)

	io.WriteString(c, "TwitterID: ")
	twitter, _, _ := bufc.ReadLine()
	twitter_id := string(twitter)

	entity, _ := db_manager.GetEntity(fmt.Sprintf("player.%s", nick))
	if entity == nil {
		all_suns_entities := db_manager.GetEntities("sun.*")
		all_suns := []entities.Sun{}
		for _, entity := range all_suns_entities {
			all_suns = append(all_suns, entity.(entities.Sun))
		}
		sun := entities.GenerateSun(all_suns, []entities.Sun{})
		hash := entities.GenerateHash(nickname)
		planets, home_planet := entities.GeneratePlanets(hash, sun.GetPosition())
		player = entities.CreatePlayer(nickname, twitter_id, home_planet)
		db_manager.SetEntity(player)
		db_manager.SetEntity(sun)
		for i := 0; i < len(planets); i++ {
			db_manager.SetEntity(planets[i])
		}
	} else {
		player = entity.(entities.Player)
	}
	return nickname, &player
}

func (self *Client) ReadLinesInto(ch chan<- string) {
	bufc := bufio.NewReader(self.conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			break
		}

		if strings.HasPrefix(line, "sm;") {
			params := strings.Split(line, ";")
			if len(params) != 4 {
				continue
			}
			fleet, _ := strconv.Atoi(params[3])
			if err := actionParser(self.nickname, params[1], params[2], fleet); err == nil {
				ch <- fmt.Sprintf("%s: %s", self.nickname, line)
			}
		} else if strings.HasPrefix(line, "scope:") {
			entity_list := db_manager.GetEntities("*")
			line := "{"
			for _, entity := range entity_list {
				switch t := entity.(type) {
				case entities.Mission, entities.Planet, entities.Player, entities.Sun:
					if key, json, err := t.Serialize(); err == nil {
						line += fmt.Sprintf("%v: %s, ", key, json)
					}
				}
			}
			line += "}"
			io.WriteString(self.conn, fmt.Sprintf("%v", line))
		}
	}
}

func (self *Client) WriteLinesFrom(ch <-chan string) {
	for msg := range ch {
		if _, err := io.WriteString(self.conn, msg); err != nil {
			return
		}
	}
}
