package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/theandyeh/gator/internal/app"
	"github.com/theandyeh/gator/internal/cmd"
	"github.com/theandyeh/gator/internal/config"
	"github.com/theandyeh/gator/internal/database"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No command provided")
		os.Exit(1)
	}

	state := &app.State{}
	var err error

	state.Cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", state.Cfg.Db_url)
	dbQueries := database.New(db)
	state.Db = dbQueries

	cmd_list := cmd.CreateCommandsList()
	cmd_list.Register("login", cmd.HandlerLogin)
	cmd_list.Register("register", cmd.HandlerRegister)
	cmd_list.Register("reset", cmd.HandlerReset)
	cmd_list.Register("users", cmd.HandlerUsers)
	cmd_list.Register("agg", cmd.HandlerAgg)

	c_name := os.Args[1]
	c_args := os.Args[2:]
	usr_command := cmd.Command{
		Name: c_name,
		Args: c_args,
	}

	err = cmd_list.Run(state, usr_command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
