package main

import "messagechannel/internal/pkg/command"

func main() {
	root := command.InitCmd()
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
