package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"

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
	sshIp := getVastIp()

	client, err := NewVastClient(sshUser, sshPassword, sshDirPath, sshIp)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer client.conn.Close()
}

func getVastIp() string {
	cmd := exec.Command("./scripts/get_fake_vast_ip.sh")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	// Sice echo command retirns a newline at the end, we want to
	// make sure ip is correctly trimmed.
	ip := strings.Trim(out.String(), "\n")

	return ip
}
