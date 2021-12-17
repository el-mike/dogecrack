package main

import (
	"log"
	"os"
	"path/filepath"

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

	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	sshIp := vast.GetFakeVastIp(path)

	client, err := vast.NewVastClient(sshUser, sshPassword, sshDirPath, sshIp)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	manager := vast.NewVastManager()

	s := server.NewServer(manager, client)

	s.Run()

	// if err := client.GetUser(); err != nil {
	// 	log.Fatal(err)
	// }
}
