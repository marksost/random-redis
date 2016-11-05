// Tests the random-redis.go file
package main

import (
	// Standard lib
	// "fmt"

	// Third-party
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/redis.v5"
)

var _ = Describe("random-redis.go", func() {
	var (
		// Test Redis server
		s *RedisServer
	)

	// Spec for the RedisServer struct and it's methods
	Describe("RedisServer", func() {
		// Spec for the NewServer method
		Describe("`NewServer` method", func() {})

		// Spec for the RedisServer's command methods
		Describe("Redis server command methods", func() {
			Describe("`Flush` method", func() {
				BeforeEach(func() {
					// Set server to a struct with predictable properties
					s = &RedisServer{
						host:   ServerHost,
						port:   1234,
						status: STATUS_RUNNING,
					}
				})

				Context("Cannot connect to the Redis server via a client", func() {
					BeforeEach(func() {
						// Set status
						// NOTE: Forces an error from `connectToRedis`
						s.setStatus(STATUS_STARTING)
					})

					It("Returns an error", func() {
						// Call method
						err := s.Flush()

						// Verify return value
						Expect(err).To(HaveOccurred())
					})
				})

				Context("Can connect to the Redis server via a client", func() {
					BeforeEach(func() {
						// Set valid client
						s.client = redis.NewClient(&redis.Options{Addr: s.Addr()})
					})

					It("Returns an error if one occurred from the `FlushAll` method of the Redis client", func() {
						// Call method
						err := s.Flush()

						// Verify return value
						Expect(err).To(HaveOccurred())
					})
				})
			})

			Describe("`Ping` method", func() {})

			Describe("`Stop` method", func() {})
		})

		// Spec for the RedisServer's info methods
		Describe("Redis server info methods", func() {
			BeforeEach(func() {
				// Set server to a struct with predictable properties
				s = &RedisServer{
					host:   "mock-host",
					id:     "mock-id",
					port:   1234,
					status: STATUS_STARTING,
				}
			})

			It("Returns a Redis server's address", func() {
				Expect(s.Addr()).To(Equal("mock-host:1234"))
			})

			It("Returns a Redis server's host", func() {
				Expect(s.Host()).To(Equal("mock-host"))
			})

			It("Returns a Redis server's ID", func() {
				Expect(s.Id()).To(Equal("mock-id"))
			})

			It("Returns a Redis server's port", func() {
				Expect(s.Port()).To(Equal(1234))
			})

			It("Gets and sets a Redis server's status", func() {
				// Reset status
				s.setStatus(STATUS_KILLED)

				// Verify status was updated
				Expect(s.GetStatus()).To(Equal(STATUS_KILLED))
			})
		})

		// Spec for the RedisServer's utility methods
		Describe("Redis server utility methods", func() {
			Describe("`connectToRedis` method", func() {
				BeforeEach(func() {
					// Set server to a struct with predictable properties
					s = &RedisServer{
						host:   "mock-host",
						port:   1234,
						status: STATUS_RUNNING,
					}
				})

				Context("The server is not running", func() {
					BeforeEach(func() {
						// Reset status
						s.setStatus(STATUS_STARTING)
					})

					It("Returns an error", func() {
						// Call method
						err := s.connectToRedis()

						// Verify return value
						Expect(err).To(HaveOccurred())
					})
				})

				Context("The server's client is already set", func() {
					BeforeEach(func() {
						// Set valid client
						s.client = redis.NewClient(&redis.Options{Addr: s.Addr()})
					})

					It("Returns nil", func() {
						// Call method
						err := s.connectToRedis()

						// Verify return value
						Expect(err).To(Not(HaveOccurred()))
					})
				})

				Context("The server's client is not set", func() {
					It("Set's the server's client and returns nil", func() {
						// Call method
						err := s.connectToRedis()

						// Verify return value
						Expect(err).To(Not(HaveOccurred()))

						// Verify client was set
						Expect(s.client).To(Not(BeNil()))
					})
				})
			})

			Describe("`start` method", func() {
				BeforeEach(func() {
					// Set server to a struct with predictable properties
					s = &RedisServer{
						host:   ServerHost,
						status: STATUS_STARTING,
					}
				})

				Context("The server is already running", func() {
					BeforeEach(func() {
						// Reset status
						s.setStatus(STATUS_RUNNING)
					})

					It("Returns an error", func() {
						// Call method
						err := s.start()

						// Verify return value
						Expect(err).To(HaveOccurred())
					})
				})

				Context("The server could not be started", func() {
					BeforeEach(func() {
						// Set invalid port
						s.port = 22

						// Set server command
						s.cmd = getNewCommand(s.port, "mock-id")
					})

					It("Returns an error", func() {
						// Call method
						err := s.start()

						// Verify return value
						Expect(err).To(HaveOccurred())
					})
				})

				Context("The server is started", func() {
					BeforeEach(func() {
						// Get random port
						port, err := getEmptyPort()
						if err != nil {
							panic("Error getting an empty port. Testing cannot continue. Error was: " + err.Error())
						}

						// Set valid port
						s.port = port

						// Set server command
						s.cmd = getNewCommand(s.port, "mock-id")
					})

					It("Returns the response from the `Start` method of the command (nil)", func() {
						// Call method
						err := s.start()

						// Verify return value
						Expect(err).To(Not(HaveOccurred()))
					})
				})
			})
		})
	})

	// Spec for helper methods for this package
	Describe("Helper methods", func() {
		var (
			// Initial server host to reset to after each test
			initialServerHost string
		)

		BeforeEach(func() {
			// Store initial value for the server host
			initialServerHost = ServerHost
		})

		AfterEach(func() {
			// Restore server host to initial value
			ServerHost = initialServerHost
		})

		It("Returns a new start command based on input", func() {
			// Call method
			cmd := getNewCommand(1234, "mock-id")

			// Verify return value
			Expect(cmd.Path).To(ContainSubstring(RedisCommand))
			Expect(cmd.Args[1]).To(Equal("--dbfilename"))
			Expect(cmd.Args[2]).To(Equal("dump.1234.mock-id.rdb"))
			Expect(cmd.Args[3]).To(Equal("--dir"))
			Expect(cmd.Args[4]).To(Equal(RedisFileLocation))
			Expect(cmd.Args[5]).To(Equal("--pidfile"))
			Expect(cmd.Args[6]).To(Equal(RedisFileLocation + "/random-redis.1234.mock-id.pid"))
			Expect(cmd.Args[7]).To(Equal("--port"))
			Expect(cmd.Args[8]).To(Equal("1234"))
		})

		Describe("getEmptyPort", func() {
			Context("When one or more ports are free on a network device", func() {
				It("Returns the port", func() {
					// Call method
					port, err := getEmptyPort()

					// Verify return values
					Expect(port).To(Not(Equal(0)))
					Expect(err).To(Not(HaveOccurred()))
				})
			})

			Context("When no ports are free on a network device", func() {
				It("Returns an error", func() {
					// Set invalid server host
					ServerHost = "invalid-address"

					// Call method
					port, err := getEmptyPort()

					// Verify return values
					Expect(port).To(Equal(0))
					Expect(err).To(HaveOccurred())
				})
			})
		})

		It("Converts a string to an int", func() {
			// Set test data
			data := map[string]int{
				"foo":  0,
				"1234": 1234,
			}

			// Loop through test data
			for input, expected := range data {
				// Call method
				actual := string2Int(input)

				// Verify result
				Expect(actual).To(Equal(expected))
			}
		})

		It("Converts a string to an int64", func() {
			// Set test data
			data := map[string]int64{
				"foo":  0,
				"1234": 1234,
			}

			// Loop through test data
			for input, expected := range data {
				// Call method
				actual := string2Int64(input)

				// Verify result
				Expect(actual).To(Equal(expected))
			}
		})
	})
})
