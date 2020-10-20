package logger

import (
	"errors"
	stdlog "log"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NoopImplementLogger(t *testing.T) {
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
	var log Noop

	output, err := CaptureOutput(func() {
		oldFlags := stdlog.Flags()
		stdlog.SetFlags(0)
		defer func() {
			stdlog.SetFlags(oldFlags)
		}()

		restore := RedirectStdLog(&log, LevelDebug)

		stdlog.Println("first")
		restore()
		stdlog.Println("second")
	})
	require.NoError(t, err)
	require.Equal(t, "second\n", output)
}

func Test_NoopUselessButUntestable(*testing.T) {
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
	_ = log.SetLevel(LevelError)
}
