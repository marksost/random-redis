// Package randomredis is a utility for starting and stopping Redis servers on random ports
// Useful for testing applications that utilize Redis within code to provide predictable i/o
package randomredis

import (
	// Standard lib
	"fmt"
	"os/exec"
	"time"

	// Third-party
	log "github.com/Sirupsen/logrus"
	goutils "github.com/marksost/go-utils"
	"github.com/satori/go.uuid"
	"gopkg.in/redis.v5"
)

const (
	// Redis server statuses
	StatusStarting = 1
	StatusRunning  = 2
	StatusKilled   = 3
)

type (
	// RedisServer is a struct representing a single Redis server listening on a random port
	RedisServer struct {
		client *redis.Client // Redis client for interacting with the Redis server
		cmd    *exec.Cmd     // The command instance that runs the Redis server
		host   string        // The host the Redis server is running on
		id     string        // Unique ID for the Redis server
		port   int           // The port the Redis server is running on
		status int           // The current status of the Redis server
	}
)

var (
	// RedisFileLocation represents the location that all files relating
	// to a Redis server should be located in
	// NOTE: Should start with a leading slash but have no trailing slash
	// NOTE: Public variable to allow package authors the ability
	// to change this before starting the Redis server
	RedisFileLocation = "/tmp"
	// RedisCommand represents the command to be run when starting a Redis server
	// NOTE: Public variable to allow package authors the ability
	// to change this before starting the Redis server
	RedisCommand = "redis-server"
	// ServerHost is the host to run the Redis server on
	// NOTE: Public variable to allow package authors the ability
	// to change this before starting the Redis server
	ServerHost = "localhost"
	// startTimeout is the time to wait for a new Redis server to start before checking for errors
	startTimeout = 200 * time.Millisecond
)

// NewServer attempts to create, start, and return a new Redis server
// operating on a random port
func NewServer() (*RedisServer, error) {
	// Set goutils server host
	goutils.ServerHost = ServerHost

	// Get random port
	port, err := goutils.GetEmptyPort()
	if err != nil {
		return nil, err
	}

	// Generate new ID
	id := uuid.NewV4().String()

	// Form new server
	s := &RedisServer{
		cmd:    getNewCommand(port, id),
		host:   ServerHost,
		id:     id,
		port:   port,
		status: StatusStarting,
	}

	// Log server start
	log.WithField("server", s).Info("Attempting to start Redis server")

	// Attempt to start the server
	if err = s.start(); err != nil {
		return nil, err
	}

	// Set server status
	s.setStatus(StatusRunning)

	// Log running status
	log.WithField("server", s).Info("Redis server running")

	return s, nil
}

/* Begin Redis server command methods */

// Flush is used to flush all key/value pairs from a Redis server
// by running a `FlushAll` command
// For more information on the `FlushAll`  command, see http://redis.io/commands/flushall
func (s *RedisServer) Flush() error {
	// Connect to Redis client if needed
	err := s.connectToRedis()
	if err != nil {
		return err
	}

	// Return result of a `FlushAll` command to the Redis server
	return s.client.FlushAll().Err()
}

// Ping returns the output of running a `Ping` command on the Redis server
// For more information on the `Ping`  command, see http://redis.io/commands/ping
func (s *RedisServer) Ping() error {
	// Connect to Redis client if needed
	err := s.connectToRedis()
	if err != nil {
		return err
	}

	// Return result of a `Ping` command to the Redis server
	return s.client.Ping().Err()
}

// Stop attempts to stop a currently-running Redis server
func (s *RedisServer) Stop() error {
	// Check that Redis server is running
	if s.GetStatus() != StatusRunning {
		return fmt.Errorf("Attempted to stop a non-running Redis server")
	}

	// Log server stop
	log.WithField("server", s).Info("Attempting to stop Redis server")

	// Set server status
	s.setStatus(StatusKilled)

	// Attempt to kill the process
	s.cmd.Process.Kill()

	// Log killed status
	log.WithField("server", s).Info("Redis server killed")

	return nil
}

/* End Redis server command methods */

/* Begin Redis server info methods */

// Addr returns the address of the Redis server with the pattern of: {host}:{port}
func (s *RedisServer) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host(), s.Port())
}

// Host returns the host of the Redis server
func (s *RedisServer) Host() string {
	return s.host
}

// ID returns the unique ID of the Redis server
func (s *RedisServer) ID() string {
	return s.id
}

// Port returns the port of the Redis server
func (s *RedisServer) Port() int {
	return s.port
}

// GetStatus returns a server's internal status property
// NOTE: It is advisable to check the value returned from this method
// against one of the status contstants defined in this package
func (s *RedisServer) GetStatus() int {
	return s.status
}

// setStatus sets a server's internal status property
func (s *RedisServer) setStatus(status int) {
	s.status = status
}

/* End Redis server info methods */

/* Begin internal utility methods */

// connectToRedis attempts to connect to a currently-running Redis server
// and sets a `client` property on it for future use
func (s *RedisServer) connectToRedis() error {
	// Check that the Redis server is running
	if s.GetStatus() != StatusRunning {
		return fmt.Errorf("Attempted to connect to a non-running Redis server")
	}

	// Check if client has already been set up
	if s.client != nil {
		return nil
	}

	// Create Redis client based on Redis server
	s.client = redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	return nil
}

// start is an internal method for starting a Redis server on a random port
func (s *RedisServer) start() error {
	// Check if server is already running
	if s.GetStatus() == StatusRunning {
		return fmt.Errorf("Attempted to start a Redis server that is already running")
	}

	// Make channel that listens for wait errors
	ch := make(chan error)

	// Start command
	resp := s.cmd.Start()

	// Wait for command to complete in a goroutine
	go func() {
		ch <- s.cmd.Wait()
	}()

	select {
	// Handle start/wait errors
	case err := <-ch:
		return fmt.Errorf("Error starting Redis server: %s", err.Error())
	// Return the result of the start call after a specified timeout
	case <-time.After(startTimeout):
		return resp
	}
}

// getNewCommand is used to form a new Redis server start command and return it
// for use by a new Redis server
func getNewCommand(port int, id string) *exec.Cmd {
	return exec.Command(RedisCommand,
		"--dbfilename", fmt.Sprintf("dump.%d.%s.rdb", port, id),
		"--dir", RedisFileLocation,
		"--pidfile", fmt.Sprintf("%s/random-redis.%d.%s.pid", RedisFileLocation, port, id),
		"--port", fmt.Sprintf("%d", port),
	)
}

/* End internal utility methods */

// NOTE: Provided for package compliance
func main() {}
