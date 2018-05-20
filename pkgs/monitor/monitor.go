package monitor

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	filesToWatch = make([]string, 0)
	watcher      *fsnotify.Watcher
)

func Monitor(files []string, commands [][]string) error {
	for _, file := range files {
		f, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("Error reading file %s", file)
		}
		switch mode := f.Mode(); {
		case mode.IsDir():
			filepath.Walk(file, addToMonitor)
		case mode.IsRegular():
			fmt.Printf("Monitoring file %s\n", f.Name())
			filesToWatch = append(filesToWatch, f.Name())
		}
	}

	watcher, _ = fsnotify.NewWatcher()
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch event.Op {
				case fsnotify.Write, fsnotify.Rename:
					err := execCmd(commands)
					if err != nil {
						fmt.Println(err)
					}
					if event.Op == fsnotify.Rename {
						watcher.Add(event.Name)
					}
				default:
					//Do nothing
				}
			case err := <-watcher.Errors:
				fmt.Printf("Watcher error: %v", err)
			}

			continue
		}
	}()

	for _, file := range filesToWatch {
		err := watcher.Add(file)
		if err != nil {
			return fmt.Errorf("failed to watch file %q: %v", file, err)
		}
	}
	<-done
	return nil
}

func addToMonitor(path string, f os.FileInfo, err error) error {
	switch mode := f.Mode(); {
	case mode.IsRegular():
		fmt.Printf("Monitoring file %s\n", path)
		filesToWatch = append(filesToWatch, path)
	}
	return nil
}

func execCmd(commands [][]string) (error) {
	var cmdToExecute string
	for _, command := range commands {
		var args []string
		if len(command) == 2 {
			//no args
			cmdToExecute = command[1]
			args = make([]string, 0)
			trimmedCommand := strings.Split(command[1], " ")
			if len(trimmedCommand) > 1 {
				cmdToExecute = trimmedCommand[0]
				args = trimmedCommand[1:len(trimmedCommand)]
			}
		} else if len(command) > 2 {
			// with args
			cmdToExecute = command[1]
			args = command[2:len(command)]

		} else {
			// no command
			return fmt.Errorf("no command specified. %+v", command)
		}

		fmt.Println(cmdToExecute, args)
		output, err := exec.Command(cmdToExecute, args...).CombinedOutput()
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Print(string(output))
	}
	return nil

}
