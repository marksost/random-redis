// Tests the random-redis.go file
package randomredis

import (
	// Third-party
	goutils "github.com/marksost/go-utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/redis.v5"
)

var _ = Describe("random-redis.go", func() {
	var (
		// Initial Redis file location to reset to after each test
		initialRedisFileLocation string
		// Initial server host to reset to after each test
		initialServerHost string
		// Error to use throughout tests
		err error
		// Test Redis server
		s *RedisServer
	)

	BeforeEach(func() {
		// Store initial value for the Redis file location
		initialRedisFileLocation = RedisFileLocation

		// Store initial value for the server host
		initialServerHost = ServerHost
	})

	AfterEach(func() {
		// Restore Redis file location to initial value
		RedisFileLocation = initialRedisFileLocation

		// Restore server host to initial value
		ServerHost = initialServerHost
	})

	// Spec for the RedisServer struct and it's methods
	Describe("RedisServer", func() {
		// Spec for the NewServer method
		Describe("`NewServer` method", func() {
			Context("Cannot get an empty port", func() {
				BeforeEach(func() {
					// Set invalid server host
					ServerHost = "invalid-address"
				})

				It("Returns an error", func() {
					// Call method
					s, err := NewServer()

					// Verify return values
					Expect(s).To(BeNil())
					Expect(err).To(HaveOccurred())
				})
			})

			Context("Cannot start the redis server", func() {
				BeforeEach(func() {
					// Set invalid Redis file location
					RedisFileLocation = "/foo/bar/baz"
				})

				It("Returns an error", func() {
					// Call method
					s, err := NewServer()

					// Verify return values
					Expect(s).To(BeNil())
					Expect(err).To(HaveOccurred())
				})
			})

			Context("The server is started", func() {
				AfterEach(func() {
					// Stop server
					s.Stop()
				})

				It("Returns the new Redis server", func() {
					// Call method
					s, err = NewServer()

					// Verify return values
					Expect(s).To(Not(BeNil()))
					Expect(err).To(Not(HaveOccurred()))

					// Verify status was set
					Expect(s.GetStatus()).To(Equal(StatusRunning))
				})
			})
		})

		// Spec for the RedisServer's command methods
		Describe("Redis server command methods", func() {
			BeforeEach(func() {
				// Set server to a struct with predictable properties
				s = &RedisServer{
					host:   ServerHost,
					port:   1234,
					status: StatusRunning,
				}
			})

			Describe("`Flush` method", func() {
				Context("Cannot connect to the Redis server via a client", func() {
					BeforeEach(func() {
						// Set status
						// NOTE: Forces an error from `connectToRedis`
						s.setStatus(StatusStarting)
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

			Describe("`Ping` method", func() {
				Context("Cannot connect to the Redis server via a client", func() {
					BeforeEach(func() {
						// Set status
						// NOTE: Forces an error from `connectToRedis`
						s.setStatus(StatusStarting)
					})

					It("Returns an error", func() {
						// Call method
						err := s.Ping()

						// Verify return value
						Expect(err).To(HaveOccurred())
					})
				})

				Context("Can connect to the Redis server via a client", func() {
					BeforeEach(func() {
						// Set valid client
						s.client = redis.NewClient(&redis.Options{Addr: s.Addr()})
					})

					It("Returns an error if one occurred from the `Ping` method of the Redis client", func() {
						// Call method
						err := s.Ping()

						// Verify return value
						Expect(err).To(HaveOccurred())
					})
				})
			})

			Describe("`Stop` method", func() {
				BeforeEach(func() {
					// Set status
					s.setStatus(StatusStarting)
				})

				Context("The server is not already running", func() {
					It("Returns an error", func() {
						// Call method
						err := s.Stop()

						// Verify return value
						Expect(err).To(HaveOccurred())
					})
				})

				Context("The server is stopped", func() {
					BeforeEach(func() {
						// Use `NewServer` to create a server
						s, err = NewServer()
					})

					It("Returns nil", func() {
						// Call method
						err := s.Stop()

						// Verify return value
						Expect(err).To(Not(HaveOccurred()))

						// Verify status was set
						Expect(s.GetStatus()).To(Equal(StatusKilled))
					})
				})
			})
		})

		// Spec for the RedisServer's info methods
		Describe("Redis server info methods", func() {
			BeforeEach(func() {
				// Set server to a struct with predictable properties
				s = &RedisServer{
					host:   "mock-host",
					id:     "mock-id",
					port:   1234,
					status: StatusStarting,
				}
			})

			It("Returns a Redis server's address", func() {
				Expect(s.Addr()).To(Equal("mock-host:1234"))
			})

			It("Returns a Redis server's host", func() {
				Expect(s.Host()).To(Equal("mock-host"))
			})

			It("Returns a Redis server's ID", func() {
				Expect(s.ID()).To(Equal("mock-id"))
			})

			It("Returns a Redis server's port", func() {
				Expect(s.Port()).To(Equal(1234))
			})

			It("Gets and sets a Redis server's status", func() {
				// Reset status
				s.setStatus(StatusKilled)

				// Verify status was updated
				Expect(s.GetStatus()).To(Equal(StatusKilled))
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
						status: StatusRunning,
					}
				})

				Context("The server is not running", func() {
					BeforeEach(func() {
						// Reset status
						s.setStatus(StatusStarting)
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
						status: StatusStarting,
					}
				})

				Context("The server is already running", func() {
					BeforeEach(func() {
						// Reset status
						s.setStatus(StatusRunning)
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
						port, err := goutils.GetEmptyPort()
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
	})
})
