package internal

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
)


type Config struct {
	Procfile string `yaml:"procfile"`
	// Port for RPC server
	Port     uint   `yaml:"port"`
	BaseDir  string `yaml:"basedir"`
	BasePort uint   `yaml:"baseport"`
	Args     []string
	// If true, exit the supervisor process if a subprocess exits with an error.
	ExitOnError bool `yaml:"exit_on_error"`
}


// -- process information structure.
type ProcInfo struct {
	Name       string
	Cmdline    string
	Cmd        *exec.Cmd
	Port       uint
	SetPort    bool
	ColorIndex int

	// True if we called stopProc to kill the process, in which case an
	// *os.ExitError is not the fault of the subprocess
	StoppedBySupervisor bool

	Mu      sync.Mutex
	Cond    *sync.Cond
	waitErr error
}


func DefaultServer(serverPort uint) string {
	if s, ok := os.LookupEnv("GOREMAN_RPC_SERVER"); ok {
		return s
	}
	return fmt.Sprintf("127.0.0.1:%d", DefaultPort())
}

func DefaultAddr() string {
	if s, ok := os.LookupEnv("GOREMAN_RPC_ADDR"); ok {
		return s
	}
	return "0.0.0.0"
}

// default Port
func DefaultPort() uint {
	s := os.Getenv("GOREMAN_RPC_PORT")
	if s != "" {
		i, err := strconv.Atoi(s)
		if err == nil {
			return uint(i)
		}
	}
	return 8555
}


// show timestamp in log
var logTime = flag.Bool("logtime", true, "show timestamp in log")
var setPorts = flag.Bool("set-ports", true, "False to avoid setting PORT env var for each subprocess")

var Mu sync.Mutex
// process informations named with proc.
var Procs []*ProcInfo
var MaxProcNameLength = 0
var re = regexp.MustCompile(`\$([a-zA-Z]+[a-zA-Z0-9_]+)`)

// read Procfile and parse it.
func ReadProcfile(cfg *Config) error {
	content, err := ioutil.ReadFile(cfg.Procfile)
	if err != nil {
		return err
	}
	Mu.Lock()
	defer Mu.Unlock()

	Procs = []*ProcInfo{}
	index := 0
	for _, line := range strings.Split(string(content), "\n") {
		tokens := strings.SplitN(line, ":", 2)
		if len(tokens) != 2 || tokens[0][0] == '#' {
			continue
		}
		k, v := strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])
		if runtime.GOOS == "windows" {
			v = re.ReplaceAllStringFunc(v, func(s string) string {
				return "%" + s[1:] + "%"
			})
		}
		proc := &ProcInfo{Name: k, Cmdline: v, ColorIndex: index}
		if *setPorts == true {
			proc.SetPort = true
			proc.Port = cfg.BasePort
			cfg.BasePort += 100
		}
		proc.Cond = sync.NewCond(&proc.Mu)
		Procs = append(Procs, proc)
		if len(k) > MaxProcNameLength {
			MaxProcNameLength = len(k)
		}
		index++
		if index >= len(Colors) {
			index = 0
		}
	}
	if len(Procs) == 0 {
		return errors.New("no valid entry")
	}
	return nil
}


func FindProc(name string) *ProcInfo {
	Mu.Lock()
	defer Mu.Unlock()

	for _, proc := range Procs {
		if proc.Name == name {
			return proc
		}
	}
	return nil
}
