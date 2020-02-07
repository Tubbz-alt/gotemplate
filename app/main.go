package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Semior001/gotemplate/app/cmd"

	"github.com/hashicorp/logutils"
	"github.com/jessevdk/go-flags"
)

// Options struct defines all cli commands and flags
type Options struct {
	ServeCmd         cmd.ServeCommand     `command:"serve"`
	MigrateDbCmd     cmd.MigrateDbCommand `command:"migrate"`
	RegisterAdminCmd cmd.RegisterAdmin    `command:"register_admin"`

	Dbg bool `long:"dbg" env:"DEBUG" description:"debug mode"`
}

const appName = "gotemplate"
const appAuthor = "semior"
const version = "unknown"

var logFlags int = log.Ldate | log.Ltime

func main() {
	fmt.Printf("%s version: %s\n", appName, version)
	var opts Options
	p := flags.NewParser(&opts, flags.Default)

	p.CommandHandler = func(command flags.Commander, args []string) error {
		setupLog(opts.Dbg)
		command.(cmd.CommonCommander).SetCommonOptions(cmd.CommonOptions{
			AppName:     appName,
			AppAuthor:   appAuthor,
			Version:     version,
			LoggerFlags: logFlags,
		})
		err := command.Execute(args)
		if err != nil {
			log.Printf("[ERROR] failed to execute command %+v", err)
		}
		return nil
	}

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

}

func setupLog(dbg bool) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: "INFO",
		Writer:   os.Stdout,
	}

	if dbg {
		logFlags = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
		filter.MinLevel = "DEBUG"
	}

	log.SetFlags(logFlags)

	log.SetOutput(filter)
}
