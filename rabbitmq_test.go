package valize

import (
	"os"
	"strings"
	"testing"
)

var RMQ_URI string = "amqp://guest:guest@localhost:5672/"
var RMQ_QUEUE string = "go.valize.test"

func skipIfNotRMQIntegrationTest(t *testing.T) {
	env := strings.ToLower(os.Getenv("RMQ_INTEGRATION_TEST"))
	if env == "" || env == "false" {
		t.Skip("skipping test; $RMQ_INTEGRATION_TEST not set")
	}
}

func TestRabbitMQQueueInit(t *testing.T) {
	skipIfNotRMQIntegrationTest(t)
	actual := &Queue{Backend: NewRabbitMQStrategy(RMQ_URI, RMQ_QUEUE)}
	if actual == nil {
		t.Error("RabbitMQStraegy failed to initialize")
	}
}

func TestRMQPop(t *testing.T) {
	skipIfNotRMQIntegrationTest(t)
	queue := &Queue{Backend: NewRabbitMQStrategy(RMQ_URI, RMQ_QUEUE)}
	queue.Backend.Clear()

	if err := queue.Push([]byte("a")); err != nil {
		t.Error("Error Pushing to queue")
	}

	if err := queue.Push([]byte("b")); err != nil {
		t.Error("Error Pushing to queue")
	}

	actual1, err1 := queue.Pop()
	actual2, err2 := queue.Pop()

	expected1 := []byte("a")
	if actual1[0] != expected1[0] {
		t.Error("Assert Fail", actual1, expected1)
	}

	expected2 := []byte("b")
	if actual2[0] != expected2[0] {
		t.Error("Assert Fail", actual2, expected2)
	}

	if err1 != nil || err2 != nil {
		t.Error("Error Popping from queue")
	}

	// Should return nothing
	actual3, err3 := queue.Pop()

	if actual3 != nil {
		t.Error("Assert Fail", actual3, nil)
	}

	// TODO: capture
	_ = err3

}
