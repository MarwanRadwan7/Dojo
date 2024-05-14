package main

import (
	"fmt"
	"os"
	"os/user"

	"dojo/repl"
)

const DOJO_LOGO = `
   ____    ___       _   ___  
  |  _ \  / _ \     | | / _ \ 
  | | | || | | | _  | || | | |
  | |_| || |_| || |_| || |_| |
  |____/  \___/  \___/  \___/ 
`

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Println(DOJO_LOGO)
	fmt.Printf("Hello %s! This is the Dojo programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
