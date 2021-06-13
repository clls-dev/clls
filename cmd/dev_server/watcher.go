package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
)

func changeWatcher(p string) error {
	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	w.SetMaxEvents(1)

	// Only notify rename and move events.
	//w.FilterOps(watcher.Rename, watcher.Move)

	// Only files that match the regular expression during file listings
	// will be watched.
	w.AddFilterHook(watcher.RegexFilterHook(regexp.MustCompile(`(\.go)$`), false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
				cmd := exec.Command("go", "mod", "tidy")
				fmt.Println("1) ", cmd)
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				if err := cmd.Run(); err != nil {
					fmt.Println(err)
				}
				cmd = exec.Command("make", "main.wasm")
				cmd.Dir = filepath.Join(os.Getenv("PWD"), "cmd", "web")
				fmt.Println("2) ", cmd)
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				if err := cmd.Run(); err != nil {
					fmt.Println(err)
				}
				cmd = exec.Command("go", "build", "-o", filepath.Join(os.Getenv("HOME"), "go", "bin", "clls"), "./cmd/clls")
				fmt.Println("3) ", cmd)
				cmd.Stderr = os.Stderr
				cmd.Stdout = os.Stdout
				if err := cmd.Run(); err != nil {
					fmt.Println(err)
				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch test_folder recursively for changes.
	if err := w.AddRecursive(p); err != nil {
		return err
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	/*for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}*/

	//fmt.Println()

	// Trigger 2 events after watcher started.
	/*go func() {
		w.Wait()
		w.TriggerEvent(watcher.Create, nil)
		w.TriggerEvent(watcher.Remove, nil)
	}()*/

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		return err
	}

	return nil
}
