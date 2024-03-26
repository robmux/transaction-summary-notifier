package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/robmux/transaction-summary-notifier/pkg/repositories"
)

func GetMailConfig() repositories.Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("MAIL_USER")
	if len(user) == 0 {
		panic("mail user missing")
	}

	password := os.Getenv("MAIL_PASSWORD")
	if len(password) == 0 {
		panic("password is empty")
	}
	host := os.Getenv("MAIL_HOST")
	if len(host) == 0 {
		panic("host is empty")
	}
	serverAddr := os.Getenv("MAIL_SERVER_ADDR")
	if len(serverAddr) == 0 {
		panic("server address is empty")
	}

	serverPortStr := os.Getenv("MAIL_SERVER_PORT")
	serverPort, err := strconv.ParseInt(serverPortStr, 10, 64)
	if err != nil {
		panic(err.Error())
	}

	return repositories.Config{
		UserMail: user,
		Password: password,

		Host:       host,
		ServerAddr: serverAddr,
		Port:       int(serverPort),
	}
}
