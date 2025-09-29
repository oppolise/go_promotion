package serverconfig

import "os"

type config struct {
	PORT string
}

func ServerConfig() *config {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &config{
		PORT: port,
	}
}
