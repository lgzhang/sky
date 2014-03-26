package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/davecheney/profile"
	"github.com/skydb/sky/cmd"
	"github.com/skydb/sky/server"
)

const (
	// DefaultPort is the port that Sky listens on by default.
	DefaultPort = 8585

	// DefaultNoSync is the default setting for syncing the data.
	DefaultNoSync = false

	// DefaultMaxDBs is the default LMDB setting for MaxDBs.
	DefaultMaxDBs = 4096

	// DefaultMaxReaders is the default LMDB setting for MaxReaders.
	DefaultMaxReaders = 126 // lmdb's default

	// DefaultStreamFlushPeriod is the default length of time between flushing
	// events from the bulk endpoint into LMDB.
	DefaultStreamFlushPeriod = 60 // seconds

	// DefaultStreamFlushThreshold is the default max number of events to stream
	// in before flushing from the bulk endpoint into LMDB.
	DefaultStreamFlushThreshold = 1000
)

var branch, commit string

var config = NewConfig()
var configPath string

var profileFlag = flag.Bool("profile", false, "enable profiling")

func init() {
	log.SetFlags(0)
	flag.UintVar(&config.Port, "port", config.Port, "the port to listen on")
	flag.UintVar(&config.Port, "p", config.Port, "the port to listen on")
	flag.StringVar(&config.DataDir, "data-dir", config.DataDir, "the data directory (defaults to ~/.sky)")
	flag.BoolVar(&config.NoSync, "no-sync", config.NoSync, "use mdb.NOSYNC option, or not")
	flag.UintVar(&config.MaxDBs, "max-dbs", config.MaxDBs, "max number of named btrees in the database (mdb.MaxDBs)")
	flag.UintVar(&config.MaxReaders, "max-readers", config.MaxReaders, "max number of concurrenly executing queries (mdb.MaxReaders)")
	flag.UintVar(&config.StreamFlushPeriod, "stream-flush-period", config.StreamFlushPeriod, "time period on which to flush streamed events")
	flag.UintVar(&config.StreamFlushThreshold, "stream-flush-threshold", config.StreamFlushThreshold, "the maximum number of events (per table) in event stream before flush")
	flag.StringVar(&configPath, "config", "", "the path to the config file")
}

func main() {
	// Parse the command line arguments and load the config file (if specified).
	flag.Parse()
	if configPath != "" {
		file, err := os.Open(configPath)
		if err != nil {
			fmt.Printf("Unable to open config: %v\n", err)
			return
		}
		defer file.Close()
		if err = config.Decode(file); err != nil {
			fmt.Printf("Unable to parse config: %v\n", err)
			os.Exit(1)
		}
	}

	// Initialize profiling.
	if *profileFlag {
		defer profile.Start(&profile.Config{CPUProfile: true, MemProfile: true, BlockProfile: true}).Stop()
	}

	// Default the data directory to ~/.sky
	if config.DataDir == "" {
		u, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		config.DataDir = filepath.Join(u.HomeDir, ".sky")
	}

	// Hardcore parallelism right here.
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Initialize
	s := server.NewServer(config.Port, config.DataDir)
	s.Version = cmd.Version()
	s.NoSync = config.NoSync
	s.MaxDBs = config.MaxDBs
	s.MaxReaders = config.MaxReaders

	// Print configuration.
	log.Printf("Sky %s (%s %s)", cmd.Version(), branch, commit)
	log.Printf("Listening on http://localhost%s", s.Addr)
	log.Println("")
	log.Println("[config]")
	log.Printf("port        = %v", config.Port)
	log.Printf("data-dir    = %v", config.DataDir)
	log.Printf("no-sync     = %v", s.NoSync)
	log.Printf("max-dbs     = %v", s.MaxDBs)
	log.Printf("max-readers = %v", s.MaxReaders)
	log.Println("")

	// Start the server.
	log.SetFlags(log.LstdFlags)
	log.Fatal(s.ListenAndServe())
}

// Config represents the configuration settings used to start skyd.
type Config struct {
	Port                 uint   `toml:"port"`
	DataDir              string `toml:"data-dir"`
	NoSync               bool   `toml:"no-sync"`
	MaxDBs               uint   `toml:"max-dbs"`
	MaxReaders           uint   `toml:"max-readers"`
	StreamFlushPeriod    uint   `toml:"stream-flush-period"`
	StreamFlushThreshold uint   `toml:"stream-flush-threshold"`
}

// NewConfig creates a new Config object with the default settings.
func NewConfig() *Config {
	return &Config{
		Port:                 DefaultPort,
		NoSync:               DefaultNoSync,
		MaxDBs:               DefaultMaxDBs,
		MaxReaders:           DefaultMaxReaders,
		StreamFlushPeriod:    DefaultStreamFlushPeriod,
		StreamFlushThreshold: DefaultStreamFlushThreshold,
	}
}

// Decode reads the contents of configuration file and populates the config object.
// Any properties that are not set in the configuration file will default to
// the value of the property before the decode.
func (c *Config) Decode(r io.Reader) error {
	if _, err := toml.DecodeReader(r, &c); err != nil {
		return err
	}
	return nil
}
