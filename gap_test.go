package gap_test

import (
	"io"
	"log"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
	m.Run()
}
