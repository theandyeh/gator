package main

import (
	"fmt"

	"github.com/theandyeh/gator/internal/config"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	cfg.SetUser("andy")

	cfg, err = config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}

	fmt.Println("Current DB User:", cfg.Current_db_user)
	fmt.Println("DB URL:", cfg.Db_url)

	return

}
