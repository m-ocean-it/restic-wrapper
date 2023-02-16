package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"restic-wrapper/config"
)

func main() {
	var mode string

	switch len(os.Args) {
	case 0:
		panic("wtf") // first argument is always a name of the executable
	case 1:
		mode = "backup"
	case 2:
		arg := os.Args[1]
		if arg == "init" {
			mode = "init"
		} else {
			log.Fatalf("unknown arg: %s\n", arg)
		}
	default:
		log.Fatalln("too many args")
	}

	conf, err := config.Build()
	if err != nil {
		log.Fatalln(err)
	}

	secr := conf.Secrets()
	{
		// Restic requires those environment variables for authenticating
		// with an S3 storage-provider.
		err = os.Setenv("AWS_ACCESS_KEY_ID", secr.Aws.KeyId)
		if err != nil {
			log.Fatalln(err)
		}
		err = os.Setenv("AWS_SECRET_ACCESS_KEY", secr.Aws.Key)
		if err != nil {
			log.Fatalln(err)
		}
	}
	err = os.Setenv("RESTIC_PASSWORD", secr.ResticPassword)
	if err != nil {
		log.Fatalln(err)
	}

	switch mode {

	case "init":
		init_repo(conf.Url())

	case "backup":
		err = backup_paths(conf.Url(), conf.BACKUP_PATHS...)
		if err != nil {
			log.Fatalln(err)
		}

		err = list_snapshots(conf.Url())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func init_repo(url string) {
	log.Println("initializing a repository")

	cmd := []string{
		"restic",
		"-r",
		"s3:" + url,
		"init",
	}

	_ = exec_command(cmd)

	// if err != nil {
	// 	if strings.Contains(err.Error(), "repository master key and config already initialized") {
	// 		return nil
	// 	}
	// }
	// return err
}

func backup_paths(url string, paths ...string) error {
	cmd := []string{
		"restic",
		"-r",
		"s3:" + url,
		"--verbose",
		"backup",
	}
	cmd = append(cmd, paths...)

	err := exec_command(cmd)

	return err
}

func list_snapshots(url string) error {
	fmt.Println("Listing snapshots:")
	cmd := []string{
		"restic",
		"-r",
		"s3:" + url,
		"snapshots",
	}
	err := exec_command(cmd)
	return err
}

func exec_command(cmd []string) error {
	if len(cmd) == 0 {
		return errors.New("empty string as input")
	}

	log.Printf("executing shell command: %s\n", cmd)

	executable, args := cmd[0], cmd[1:]

	out, err := exec.Command(executable, args...).CombinedOutput()
	log.Println("\n", string(out))

	return err
}
