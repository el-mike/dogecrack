package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/el-mike/dogecrack/shepherd/persist"
	"github.com/el-mike/dogecrack/shepherd/server"
	"github.com/el-mike/dogecrack/shepherd/vast"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	sshUser := os.Getenv("SSH_USER")
	sshPassword := os.Getenv("SSH_PASSWORD")
	sshDirPath := os.Getenv("SSH_DIR")

	vastApiSecret := os.Getenv("VAST_API_SECRET")

	// walletString := os.Getenv("WALLET_STRING")

	mongoUser := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")
	mongoPort := os.Getenv("MONGO_PORT")

	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	sshIp := vast.GetFakeVastIp(path)

	client, err := vast.NewVastClient(sshUser, sshPassword, sshDirPath, sshIp)
	if err != nil {
		panic(err)
	}

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	defer client.Close()

	manager := vast.NewVastManager(vastApiSecret)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := persist.InitMongo(ctx, mongoUser, mongoPassword, mongoHost, mongoPort)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	s := server.NewServer(manager, client)
	s.Run()

	// if err := client.GetUser(); err != nil {
	// 	log.Fatal(err)
	// }
}
