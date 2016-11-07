// Tests the random-redis.go file
package main

import (
	// Third-party
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RedisServer integration tests", func() {
	var (
		// Error to use throughout tests
		err error
		// Test Redis server
		s *RedisServer
	)

	BeforeEach(func() {
		// Start new Redis server
		s, err = NewServer()
		if err != nil {
			panic("Error starting test Redis server. Testing cannot continue. Error was: " + err.Error())
		}

		// Set up Redis client for tests
		s.connectToRedis()
	})

	AfterEach(func() {
		// Stop redis server
		err = s.Stop()
		if err != nil {
			panic("Error stopping test Redis server. Testing cannot continue. Error was: " + err.Error())
		}
	})

	It("Should start and stop a functional Redis server without errors", func() {
		// Start new Redis server
		s, err := NewServer()
		Expect(s).To(Not(BeNil()))
		Expect(err).To(Not(HaveOccurred()))

		// Stop Redis server
		err = s.Stop()
		Expect(err).To(Not(HaveOccurred()))
	})

	It("Should be a ping-able Redis server", func() {
		// Attempt to call `Ping` on the Redis server
		err = s.Ping()
		Expect(err).To(Not(HaveOccurred()))
	})

	It("Should be a flush-able Redis server", func() {
		// Set keys in Redis
		s.client.Set("foo", "bar", 0)
		s.client.Set("test", 1234, 0)

		// Verify keys were set
		keys, err := s.client.Keys("*").Result()
		Expect(len(keys)).To(Equal(2))
		Expect(err).To(Not(HaveOccurred()))

		// Attempt to call `FlushAll` on the Redis server
		s.Flush()

		// Verify keys were flushed
		keys, err = s.client.Keys("*").Result()
		Expect(len(keys)).To(Equal(0))
		Expect(err).To(Not(HaveOccurred()))
	})
})
