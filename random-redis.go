// random-redis is a utility for starting and stopping Redis servers on random ports
// Useful for testing applications that utilize Redis within code to provide predictable i/o
package main

import (
	// Standard lib
	"fmt"
	"net"
	"regexp"
	"strconv"

	// Third-party
	log "github.com/Sirupsen/logrus"
)

const (
	// Redis server statuses
	STATUS_STARTING = 1
	STATUS_RUNNING  = 2
	STATUS_STOPPED  = 3
)

type (
	// Struct representing a single Redis server listening on a random port
	RedisServer struct {
		host   string // The host the Redis server is running on
		port   int    // The port the Redis server is running on
		status int    // The current status of the Redis server
	}
)

var (
	// Command to be run when starting a Redis server
	// NOTE: Public variable to allow package authors the ability
	// to change this before starting the Redis server
	RedisCommand string = "redis-server"
	// The host to run the Redis server on
	// NOTE: Public variable to allow package authors the ability
	// to change this before starting the Redis server
	ServerHost string = "localhost"
)

// NewServer attempts to create, start, and return a new Redis server
// operating on a random port
func NewServer() (*RedisServer, error) {
	// Get random port
	port, err := getEmptyPort()
	if err != nil {
		return nil, err
	}

	// Form new server
	s := &RedisServer{
		host:   ServerHost,
		port:   port,
		status: STATUS_STARTING,
	}

	// Log server
	log.WithField("server", s).Info("Attempting to start Redis server")

	// Attempt to start the server
	err = s.start()
	if err != nil {
		return nil, err
	}

	return s, nil
}

/* Begin Redis server command methods */

// Flush is used to flush all key/value pairs from a Redis server
func (r *RedisServer) Flush() error {
	// TO-DO: Fill this method in
	return nil
}

// Stop attempts to stop a currently-running Redis server
func (r *RedisServer) Stop() error {
	// TO-DO: Fill this method in
	return nil
}

/* End Redis server command methods */

/* Begin Redis server info methods */

// Address returns the address of the Redis server with the pattern of: {host}:{port}
func (r *RedisServer) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host(), r.Port())
}

// Host returns the host of the Redis server
func (r *RedisServer) Host() string {
	return r.host
}

// Info returns the output of running an "Info" command on the Redis server
// NOTE: the output will be returned as a map of strings
// For more information on the "Info" command, see http://redis.io/commands/info
func (r *RedisServer) Info() (map[string]string, error) {
	// TO-DO: Fill this method in
	return nil, nil
}

// Port returns the port of the Redis server
func (r *RedisServer) Port() int {
	return r.port
}

/* End Redis server info methods */

/* Begin internal utility methods */

// GetStatus returns a server's internal status property
// NOTE: It is advisable to check the value returned from this method
// against one of the status contstants defined in this package
func (r *RedisServer) GetStatus() int {
	return r.status
}

// setStatus sets a server's internal status property
func (r *RedisServer) setStatus(status int) {
	r.status = status
}

// start is an internal method for starting a Redis server on a random port
func (r *RedisServer) start() error {
	// TO-DO: Fill this method in
	return nil
}

// getEmptyPort returns a number to be used as a new server's port
// NOTE: Uses tcp to allow the kernel to give an open port
func getEmptyPort() (int, error) {
	// Create regex for extracting port
	r, _ := regexp.Compile("\\d+$")

	// NOTE: Uses "port" 0 to allow the kernal to chose a port for itself
	if l, err := net.Listen("tcp", fmt.Sprintf("%s:0", ServerHost)); err == nil {
		// Close listener
		defer l.Close()

		// Use regex to extract port
		port := r.FindString(l.Addr().String())

		if len(port) != 0 {
			return string2Int(port), nil
		}
	}

	return 0, fmt.Errorf("No random ports were found")
}

// string2Int converts a string to an int
func string2Int(v string) int {
	return int(string2Int64(v))
}

// string2Int64 converts a string to an int64
func string2Int64(v string) int64 {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0
	}

	return i
}

/* End internal utility methods */

// NOTE: Provided for package compliance
func main() {}
