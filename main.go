package main

import (
	"fmt"
	"os"

	"github.com/theandyeh/gator/internal/app"
	"github.com/theandyeh/gator/internal/cmd"
	"github.com/theandyeh/gator/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No command provided")
		return
	}

	state := &app.State{}
	var err error

	state.Cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd_list := cmd.CreateCommandsList()
	cmd_list.Register("login", cmd.HandlerLogin)

	c_name := os.Args[1]
	c_args := os.Args[2:]
	usr_command := cmd.Command{
		Name: c_name,
		Args: c_args,
	}

	err = cmd_list.Run(state, usr_command)
	if err != nil {
		fmt.Println(err)
		return
	}
}
