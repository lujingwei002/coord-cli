package main

import (
	"bytes"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/urfave/cli"
)

func getCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	return path, nil
}

func getProjectDir() (string, error) {
	scriptPath, err := getCurrentPath()
	if err != nil {
		return "", err
	}
	return path.Dir(path.Dir(scriptPath)), nil
}

func main() {
	app := cli.NewApp()
	app.Name = "coord-cli"
	app.Author = "lujingwei"
	app.Email = "lujingwei002@qq.com"
	app.Description = "Coord CLI"
	// flags
	app.Flags = []cli.Flag{}
	app.Action = runAction

	//log.SetFlags(log.LstdFlags | log.Lshortfile)

	app.Commands = []cli.Command{
		{
			Name: "start",
			//Aliases: []string{"start"},
			Usage:  "start <server name>",
			Action: startAction,
		},
		{
			Name:   "status",
			Usage:  "status <server name>",
			Action: statusAction,
		},
		{
			Name:   "stop",
			Usage:  "stop <server name>",
			Action: stopAction,
		},
		{
			Name:   "restart",
			Usage:  "restart <server name>",
			Action: restartAction,
		},
		{
			Name:   "reload",
			Usage:  "reload <server name>",
			Action: reloadAction,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("[main] startup server error %+v", err)
	}
}

func execCommand(name string, arg []string) (string, error) {
	cmd := exec.Command(name, arg...)
	//cmd.Stdin = strings.NewReader("")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return "", err
	}
	fmt.Println(out.String())
	return out.String(), nil
}

func startAction(args *cli.Context) error {
	if args.NArg() <= 0 {
		return nil
	}
	name := args.Args().Get(0)
	projectDir, err := getProjectDir()
	if err != nil {
		return err
	}
	configPath := fmt.Sprintf("%s/supervisord/supervisord.conf", projectDir)
	execCommand("supervisorctl", []string{"-c", configPath, "start", name})
	return nil
}

func stopAction(args *cli.Context) error {
	if args.NArg() <= 0 {
		return nil
	}
	name := args.Args().Get(0)
	projectDir, err := getProjectDir()
	if err != nil {
		return err
	}
	configPath := fmt.Sprintf("%s/supervisord/supervisord.conf", projectDir)
	execCommand("supervisorctl", []string{"-c", configPath, "stop", name})
	return nil
}

func restartAction(args *cli.Context) error {
	if args.NArg() <= 0 {
		return nil
	}
	name := args.Args().Get(0)
	projectDir, err := getProjectDir()
	if err != nil {
		return err
	}
	configPath := fmt.Sprintf("%s/supervisord/supervisord.conf", projectDir)
	execCommand("supervisorctl", []string{"-c", configPath, "restart", name})
	return nil
}

func statusAction(args *cli.Context) error {
	projectDir, err := getProjectDir()
	if err != nil {
		return err
	}
	configPath := fmt.Sprintf("%s/supervisord/supervisord.conf", projectDir)
	execCommand("supervisorctl", []string{"-c", configPath, "status"})
	return nil
}

func reloadAction(args *cli.Context) error {
	projectDir, err := getProjectDir()
	if err != nil {
		return err
	}
	configPath := fmt.Sprintf("%s/supervisord/supervisord.conf", projectDir)
	execCommand("supervisorctl", []string{"-c", configPath, "reload"})
	return nil
}

func runAction(args *cli.Context) error {
	return nil
}
