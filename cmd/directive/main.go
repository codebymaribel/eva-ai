package main

import(
	"fmt"
	"os"
	"github.com/codebymaribel/eva-ai/internal/app"
)

func main(){
	if err := app.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}