package logger

import (
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNoopImplementLogger(t *testing.T) {
	var i interface{} = new(Noop)
	if _, ok := i.(Logger); !ok {
		t.Fatalf("expected %t to implement Logger", i)
	}

	i = Noop{}
	if _, ok := i.(Logger); !ok {
		t.Fatalf("expected %t to implement Logger", i)
	}
}

func TestNoop_RedirectStdLog(t *testing.T) {
	var nope Noop

	output, err := CaptureOutput(func() {
		oldFlags := log.Flags()
		log.SetFlags(0)
		defer func() {
			log.SetFlags(oldFlags)
		}()

		restore := RedirectStdLog(&nope, LevelDebug)

		log.Println("first")
		restore()
		log.Println("second")
	})
	require.NoError(t, err)
	require.Equal(t, "second\n", output)
}

func TestNoopUselessButUntestable(*testing.T) {
	var log Noop

	log.Debug("debug")
	log.Debugf("debug")
	log.Info("info")
	log.Infof("info")
	log.Warn("warn")
	log.Warnf("warn")
	log.Error("error")
	log.Errorf("error")
	log.WithError(errors.New("eww")).Info("info")
	log.WithField("a", "b").Info("info")
	log.WithFields(map[string]interface{}{"a": "b"}).Info("info")
	log.SetLevel(LevelError) // nolint: errcheck, gosec
}
