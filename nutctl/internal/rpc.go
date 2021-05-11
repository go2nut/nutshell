package internal

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// Nutctl is RPC server
type Nutctl struct {
	rpcChan chan<- *RpcMessage
}

type RpcMessage struct {
	Msg  string
	Args []string
	// sending error (if any) when the task completes
	ErrCh chan error
}

// Start do start
func (r *Nutctl) Start(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	for _, arg := range args {
		if err = startProc(arg, nil, nil); err != nil {
			break
		}
	}
	return err
}

// Stop do stop
func (r *Nutctl) Stop(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	errChan := make(chan error, 1)
	r.rpcChan <- &RpcMessage{
		Msg:   "stop",
		Args:  args,
		ErrCh: errChan,
	}
	err = <-errChan
	return
}

// StopAll do stop all
func (r *Nutctl) StopAll(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	for _, proc := range Procs {
		if err = stopProc(proc.Name, nil); err != nil {
			break
		}
	}
	return err
}

// Restart do restart
func (r *Nutctl) Restart(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	for _, arg := range args {
		if err = restartProc(arg); err != nil {
			break
		}
	}
	return err
}

// RestartAll do restart all
func (r *Nutctl) RestartAll(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	for _, proc := range Procs {
		if err = restartProc(proc.Name); err != nil {
			break
		}
	}
	return err
}

// List do list
func (r *Nutctl) List(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	*ret = ""
	for _, proc := range Procs {
		*ret += proc.Name + "\n"
	}
	return err
}

// Status do status
func (r *Nutctl) Status(args []string, ret *string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	*ret = ""
	for _, proc := range Procs {
		if proc.Cmd != nil {
			*ret += "*" + proc.Name + "\n"
		} else {
			*ret += " " + proc.Name + "\n"
		}
	}
	return err
}

// command: run.
func Run(cmd string, args []string, serverPort uint) error {
	client, err := rpc.Dial("tcp", DefaultServer(serverPort))
	if err != nil {
		return err
	}
	defer client.Close()
	var ret string
	switch cmd {
	case "start":
		return client.Call("Nutctl.Start", args, &ret)
	case "stop":
		return client.Call("Nutctl.Stop", args, &ret)
	case "stop-all":
		return client.Call("Nutctl.StopAll", args, &ret)
	case "restart":
		return client.Call("Nutctl.Restart", args, &ret)
	case "restart-all":
		return client.Call("Nutctl.RestartAll", args, &ret)
	case "list":
		err := client.Call("Nutctl.List", args, &ret)
		fmt.Print(ret)
		return err
	case "status":
		err := client.Call("Nutctl.Status", args, &ret)
		fmt.Print(ret)
		return err
	}
	return errors.New("unknown command")
}

// start rpc server.
func StartServer(ctx context.Context, rpcChan chan<- *RpcMessage, listenPort uint) error {
	gm := &Nutctl{
		rpcChan: rpcChan,
	}
	rpc.Register(gm)
	server, err := net.Listen("tcp", fmt.Sprintf("%s:%d", DefaultAddr(), listenPort))
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	var acceptingConns = true
	for acceptingConns {
		conns := make(chan net.Conn, 1)
		go func() {
			conn, err := server.Accept()
			if err != nil {
				return
			}
			conns <- conn
		}()
		select {
		case <-ctx.Done():
			acceptingConns = false
			break
		case client := <-conns: // server is not canceled.
			wg.Add(1)
			go func() {
				defer wg.Done()
				rpc.ServeConn(client)
			}()
		}
	}
	done := make(chan struct{}, 1)
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	select {
	case <-done:
		return nil
	case <-time.After(10 * time.Second):
		return errors.New("RPC server did not shut down in 10 seconds, quitting")
	}
}
