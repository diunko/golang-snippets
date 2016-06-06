package main

import (
	//"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
)

func doFork(command string, argRest ...string) error {

	cmdParts := append([]string{}, argRest...)
	log.Infof("Execute: %s %s", command, strings.Join(cmdParts, " "))

	cmd := exec.Command(command, cmdParts...)
	cmd.Env = os.Environ()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}

	sigs := make(chan os.Signal)
	stop := make(chan bool)

	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		for {
			select {
			case sig := <-sigs:
				log.Infof("received %s, sending it to terraform process", sig)
				cmd.Process.Signal(sig)
			case <-stop:
				return
			}
		}
	}()
	defer func() {
		stop <- true
	}()

	return cmd.Wait()
}

func main() {
	go func() {
		time.Sleep(2 * time.Second)
		log.Infof("panicking")
		//err := log.Errorf("ohmy!")
		panic("ohmy")
	}()

	if err := doFork("./run.sh"); err != nil {
		log.Errorf("error in doFork: %# v, %s", err, err)
	}

}
