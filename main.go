// Real-time massively multiplayer online space strategy arcade browser game!
package main

import (
	"os"
	"os/signal"
	"syscall"

	"warcluster/config"
	"warcluster/entities/db"
	"warcluster/leaderboard"
	"warcluster/server"
)

var cfg config.Config

func main() {
	go final()

	cfg.Load()
	db.InitPool(cfg.Database.Host, cfg.Database.Port, 8)
	server.ExportConfig(cfg)
	server.InitLeaderboard(leaderboard.New())
	server.SpawnDbMissions()
	server.Start()
}

func final() {
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT)
	signal.Notify(exitChan, syscall.SIGKILL)
	signal.Notify(exitChan, syscall.SIGTERM)
	<-exitChan

	server.Stop()
	os.Exit(0)
}
