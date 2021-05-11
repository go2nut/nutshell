package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"nutshell/nutctl/internal"
	"os"
	"sort"
	"strings"
)

// version is the git tag at the time of build and is used to denote the
// binary's current version. This value is supplied as an ldflag at compile
// time by goreleaser (see .goreleaser.yml).
const (
	name     = "nutctl"
	version  = "0.3.7"
	revision = "HEAD"
)

func usage() {
	fmt.Fprint(os.Stderr, `Tasks:
  nutctl check                      # Show entries in Procfile
  nutctl help [TASK]                # Show this help
  nutctl export [FORMAT] [LOCATION] # Export the apps to another process
                                       (upstart)
  nutctl run COMMAND [PROCESS...]   # Run a command
                                       start
                                       stop
                                       stop-all
                                       restart
                                       restart-all
                                       list
                                       status
  nutctl start [PROCESS]            # Start the application
  nutctl version                    # Display Nutctl version

Options:
`)
	flag.PrintDefaults()
	os.Exit(0)
}


// filename of Procfile.
var procfile = flag.String("f", ".nutshell/apps.Procfile", "proc file")

// rpc port number.
var port = flag.Uint("p", internal.DefaultPort(), "port")

// base directory
var basedir = flag.String("basedir", "", "base directory")

// base of port numbers for app
var baseport = flag.Uint("b", 5000, "base number of port")


// true to exit the supervisor
var exitOnError = flag.Bool("exit-on-error", false, "Exit nutctl if a subprocess quits with a nonzero return code")

func readConfig() *internal.Config {
	var cfg internal.Config

	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}

	cfg.Procfile = *procfile
	cfg.Port = *port
	cfg.BaseDir = *basedir
	cfg.BasePort = *baseport
	cfg.ExitOnError = *exitOnError
	cfg.Args = flag.Args()

	b, err := ioutil.ReadFile(".nutctl")
	if err == nil {
		yaml.Unmarshal(b, &cfg)
	}
	return &cfg
}

// command: check. show Procfile entries.
func Check(cfg *internal.Config) error {
	err := internal.ReadProcfile(cfg)
	if err != nil {
		return err
	}

	internal.Mu.Lock()
	defer internal.Mu.Unlock()

	keys := make([]string, len(internal.Procs))
	i := 0
	for _, proc := range internal.Procs {
		keys[i] = proc.Name
		i++
	}
	sort.Strings(keys)
	fmt.Printf("valid procfile detected (%s)\n", strings.Join(keys, ", "))
	return nil
}


// command: start. spawn procs.
func start(ctx context.Context, sig <-chan os.Signal, cfg *internal.Config) error {
	err := internal.ReadProcfile(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	// Cancel the RPC server when procs have returned/errored, cancel the
	// context anyway in case of early return.
	defer cancel()
	if len(cfg.Args) > 1 {
		tmp := make([]*internal.ProcInfo, 0, len(cfg.Args[1:]))
		internal.MaxProcNameLength = 0
		for _, v := range cfg.Args[1:] {
			proc := internal.FindProc(v)
			if proc == nil {
				return errors.New("unknown proc: " + v)
			}
			tmp = append(tmp, proc)
			if len(v) > internal.MaxProcNameLength {
				internal.MaxProcNameLength = len(v)
			}
		}
		internal.Mu.Lock()
		internal.Procs = tmp
		internal.Mu.Unlock()
	}
	godotenv.Load()
	rpcChan := make(chan *internal.RpcMessage, 10)
	go internal.StartServer(ctx, rpcChan, cfg.Port)
	procsErr := internal.StartProcs(sig, rpcChan, cfg.ExitOnError)
	return procsErr
}

func showVersion() {
	fmt.Fprintf(os.Stdout, "%s\n", version)
	os.Exit(0)
}

func main() {
	var err error
	cfg := readConfig()

	if cfg.BaseDir != "" {
		err = os.Chdir(cfg.BaseDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "nutctl: %s\n", err.Error())
			os.Exit(1)
		}
	}

	cmd := cfg.Args[0]
	switch cmd {
	case "check":
		err = Check(cfg)
	case "help":
		usage()
	case "run":
		if len(cfg.Args) >= 2 {
			cmd, args := cfg.Args[1], cfg.Args[2:]
			err = internal.Run(cmd, args, cfg.Port)
		} else {
			usage()
		}
	case "export":
		if len(cfg.Args) == 3 {
			format, path := cfg.Args[1], cfg.Args[2]
			err = internal.Export(cfg, format, path)
		} else {
			usage()
		}
	case "start":
		c := internal.NotifyCh()
		err = start(context.Background(), c, cfg)
	case "version":
		showVersion()
	default:
		usage()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err.Error())
		os.Exit(1)
	}
}
