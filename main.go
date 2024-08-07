/*
Copyright © 2024 Subham Bharadwaz subhamsbharadwaz@gmail.com
*/
package main

import (
	"log"

	"github.com/subhamBharadwaz/go-todo-cli-app/cmd"
	"github.com/subhamBharadwaz/go-todo-cli-app/utils"
)

func main() {
	err := utils.OpenDB()
	if err != nil {
		log.Panic(err)
	}
	defer utils.CloseDB()

	err = utils.SetupDB()
	if err != nil {
		log.Panic(err)
	}

	cmd.Execute()

}
