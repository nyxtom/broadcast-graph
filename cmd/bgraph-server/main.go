package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/nyxtom/bgraph"
	"github.com/nyxtom/broadcast/backends/bdefault"
	"github.com/nyxtom/broadcast/protocols/line"
	"github.com/nyxtom/broadcast/server"
)

type Configuration struct {
	port int    // port of the server
	host string // host of the server
}

var LogoHeader = `%s %s %s, Port: %d, PID: %d`

func main() {
	// Leverage all cores available
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse out flag parameters
	var host = flag.String("h", "127.0.0.1", "bgraph host to bind to")
	var port = flag.Int("p", 7331, "bgraph port to bind to")
	var configFile = flag.String("config", "", "bgraph configuration file (/etc/bgraph.conf)")
	var cpuProfile = flag.String("cpuprofile", "", "write cpu profile to file")

	flag.Parse()

	cfg := &Configuration{*port, *host}
	if len(*configFile) == 0 {
		fmt.Printf("[%d] %s # WARNING: no config file specified, using the default config\n", os.Getpid(), time.Now().Format(time.RFC822))
	} else {
		data, err := ioutil.ReadFile(*configFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = toml.Decode(string(data), cfg)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			fmt.Println(err)
			return
		}
		pprof.StartCPUProfile(f)
	}

	// create a new broadcast server
	app, err := server.ListenProtocol(cfg.port, cfg.host, lineProtocol.NewLineProtocol())
	app.Header = ""
	app.Name = "BGraph"
	app.Version = "0.1.0"
	app.Header = LogoHeader
	if err != nil {
		fmt.Println(err)
		return
	}

	// setup default backend
	backend, err := bdefault.RegisterBackend(app)
	if err != nil {
		fmt.Println(err)
		return
	}
	app.LoadBackend(backend)

	backend, err = bgraph.RegisterBackend(app)
	if err != nil {
		fmt.Println(err)
		return
	}
	app.LoadBackend(backend)

	// wait for all events to fire so we can log them
	pid := os.Getpid()
	go func() {
		for !app.Closed {
			event := <-app.Events
			t := time.Now()
			delim := "#"
			if event.Level == "error" {
				delim = "ERROR:"
			}
			msg := fmt.Sprintf("[%d] %s %s %s", pid, t.Format(time.RFC822), delim, event.Message)
			if event.Err != nil {
				msg += fmt.Sprintf(" %v", event.Err)
			}

			fmt.Println(msg)
		}
	}()

	go func() {
		<-app.Quit
		pprof.StopCPUProfile()
		os.Exit(0)
	}()

	// attach to any signals that would cause our app to close
	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt)

	go func() {
		<-sc
		app.Close()
	}()

	// accept incomming connections!
	app.AcceptConnections()
}
