package main

import (
	"embed"
	"fmt"
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
	"github.com/devendrasankhla/short-my-url/oltp/internals/database"
	messagequeue "github.com/devendrasankhla/short-my-url/oltp/internals/messageQueue"
	"github.com/joho/godotenv"
)

type application struct {
	models        database.Models
	counter       uint64
	messaageQueue messagequeue.MessageQueue
	snowflakeNode *snowflake.Node
	cfg           config
}

type config struct {
	port   int
	nodeId int64
}

//go:embed web/*
var embedFrontend embed.FS

func loadConfig() (config, error) {
	godotenv.Load()

	nodeStr := os.Getenv("NODE_ID")
	if nodeStr == "" {
		fmt.Println("NODE_ID is not set, setting to 0")
		nodeStr = "0"
	}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		fmt.Println("Port number not valid, defaulting to 4000")
		portStr = "4000"
	}

	nodeID, err := strconv.ParseInt(nodeStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		return config{}, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println(err)
		return config{}, err
	}

	return config{
		nodeId: nodeID,
		port:   port,
	}, err
}

func main() {
	db, err := database.NewFileDB("D:\\Downloads\\urlShortner")
	if err != nil {
		fmt.Println("Error initializing database : ", err)
		return
	}

	rabbitMQ, err := messagequeue.NewRabbmitMQ()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rabbitMQ.Close()

	cfg, err := loadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("cfg", cfg)

	snowflakeNode, err := snowflake.NewNode(cfg.nodeId)
	if err != nil {
		fmt.Println(err)
		return
	}

	app := &application{
		models:        database.NewModels(db),
		messaageQueue: rabbitMQ,
		snowflakeNode: snowflakeNode,
		cfg:           cfg,
	}

	err = app.serve()
	if err != nil {
		fmt.Println(err)
	}
}
