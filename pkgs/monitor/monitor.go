package monitor

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var (
	filesToWatch = make([]string, 0)
	watcher      *fsnotify.Watcher
)

func Monitor(files, command []string) error {
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
					stdout, err := execCmd(command...)
					if err == nil {
						fmt.Print(stdout)
					} else {
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

func execCmd(command ...string) (string, error) {
	var stdout bytes.Buffer
	var cmdToExecute string
	var args []string
	if len(command) == 2 {
		//no args
		cmdToExecute = command[1]
		args = make([]string, 0)
	} else if len(command) > 2 {
		// with args
		cmdToExecute = command[1]
		args = command[2:len(command)]
	} else {
		// no command
		return "", fmt.Errorf("no command specified. %+v", command)
	}

	cmd := exec.Command(cmdToExecute, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = ioutil.Discard

	err := cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
			if waitStatus.ExitStatus() == 1 {
				return "", fmt.Errorf("Error executing command: %s", command)
			}
		}
		return "", err
	}
	return stdout.String(), nil
}
