package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func WriteConfig(Config string, File string, ReloadCmd string) {
	oldConfig := ""
	oldConfBytes, err := os.ReadFile(File)
	if err != nil && err.Error() != fmt.Sprintf("open %s: no such file or directory", File) {
		log.Fatal(err)
	} else if err == nil {
		oldConfig = string(oldConfBytes)
	}

	if oldConfig == Config {
		return
	}

	log.Print("Updating config")
	os.WriteFile(File, []byte(Config), 0644)

	if ReloadCmd == "" {
		return
	}

	log.Print("Executing ReloadCmd")
	out, cmdErr := exec.Command("/bin/sh", "-c", ReloadCmd).CombinedOutput()
	log.Print(out)
	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
}
