// Tests the random-redis.go file
package main

import (
	// Standard lib
	"fmt"

	// Third-party
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("random-redis.go", func() {
	var (
		// Test Redis server
		s *RedisServer
	)

	BeforeEach(func() {
		// Error to be used when creating a new test server
		var err error

		// Create new Redis server
		s, err = NewServer()
		if err != nil {
			panic("Error creating new server. Testing cannot continue. Error was: " + err.Error())
		}
	})

	Describe("RedisServer stuct methods", func() {
		It("Should return a new server, operating on a random port", func() {})
		It("Should flush all key/value pairs from the Redis server", func() {})
		It("Should stop the Redis server", func() {})

		It("Should return the address of the Redis server", func() {
			// Verify method return value
			Expect(s.Addr()).To(Equal(fmt.Sprintf("%s:%d", s.host, s.port)))
		})

		It("Should return the host of the Redis server", func() {
			// Verify method return value
			Expect(s.Host()).To(Equal(s.host))
		})

		It("Should return the ID of the Redis server", func() {
			// Verify method return value
			Expect(s.Id()).To(Equal(s.id))
		})

		It("Should return the result of an `Info` call to the Redis server", func() {})

		It("Should return the port of the Redis server", func() {
			// Verify method return value
			Expect(s.Port()).To(Equal(s.port))
		})

		It("Should return the `status` property of the server", func() {
			// Set predictable status
			s.setStatus(STATUS_STOPPED)

			// Verify method return value
			Expect(s.GetStatus()).To(Not(Equal(STATUS_STARTING)))
			Expect(s.GetStatus()).To(Equal(STATUS_STOPPED))
		})

		It("Should set the `status` property of the server", func() {
			// Verify initial value
			Expect(s.GetStatus()).To(Equal(STATUS_STARTING))

			// Call method
			s.setStatus(STATUS_STOPPED)

			// Verify method return value
			Expect(s.GetStatus()).To(Not(Equal(STATUS_STARTING)))
			Expect(s.GetStatus()).To(Equal(STATUS_STOPPED))
		})

		It("Should start a new Redis server", func() {})
	})

	Describe("Utility methods", func() {
		It("Should return a new shell command based on input", func() {})
		It("Should return a psudo-random empty port for use by a new Redis server", func() {})

		It("Should convert a string to an int", func() {
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

		It("Should convert a string to an int64", func() {
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
