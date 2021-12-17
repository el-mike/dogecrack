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

	// cmd := exec.Command("ssh", "root@172.20.0.2")
	// out, err := cmd.Output()
	// if err != nil {
	// 	fmt.Print(err.Error())
	// }

	// fmt.Print(string(out))

	sshUser := os.Getenv("SSH_USER")
	sshPassword := os.Getenv("SSH_PASSWORD")
	sshDirPath := os.Getenv("SSH_DIR")
	sshIp := getVastIp()

	client, err := NewVastClient(sshUser, sshPassword, sshDirPath, sshIp)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	if err := client.GetUser(); err != nil {
		log.Fatal(err)
	}
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
