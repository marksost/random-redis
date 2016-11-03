// Tests the random-redis.go file
package main

import (
	// Third-party
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("random-redis.go", func() {
	BeforeEach(func() {})

	Describe("RedisServer stuct methods", func() {
		It("Should return a new server, operating on a random port", func() {})

		It("Should flush all key/value pairs from the Redis server", func() {})

		It("Should stop the Redis server", func() {})

		It("Should return the address of the Redis server", func() {})

		It("Should return the host of the Redis server", func() {})

		It("Should return the result of an `Info` call to the Redis server", func() {})

		It("Should return the port of the Redis server", func() {})

		It("Should return the `status` property of the server", func() {})

		It("Should set the `status` property of the server", func() {})

		It("Should start a new Redis server", func() {})
	})

	Describe("Utility methods", func() {
		It("Should return a psudo-random empty port for use by a new Redis server", func() {})
	})
})
