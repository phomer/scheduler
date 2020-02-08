package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
)

func main() {
	fmt.Println("register")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("%s", err)
	}
	_ = watcher
}
