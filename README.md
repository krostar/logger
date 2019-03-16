# logger

[![godoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge)](https://godoc.org/github.com/krostar/logger)
[![Licence](https://img.shields.io/github/license/krostar/logger.svg?style=for-the-badge)](https://tldrlegal.com/license/mit-license)
![Latest version](https://img.shields.io/github/tag/krostar/logger.svg?style=for-the-badge)

[![Build Status](https://img.shields.io/travis/krostar/logger/master.svg?style=for-the-badge)](https://travis-ci.org/krostar/logger)
[![Code quality](https://img.shields.io/codacy/grade/219a45ca1028442f816c745fcedbb111/master.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/logger/dashboard)
[![Code coverage](https://img.shields.io/codacy/coverage/219a45ca1028442f816c745fcedbb111.svg?style=for-the-badge)](https://app.codacy.com/project/krostar/logger/dashboard)

One logger to rule them all.

## Motivation

I was using **logrus** for some times when I discovered **zap**. I wanted to try out zap in a project where I already used logrus. I was first thinking it should only be a matter of replacing all logrus occurence with zap's one modulo some renaming. But it was not that easy: interfaces were different, exposed API contained special types (for instance `logrus.Field`) or tooks differents parameters, logrus was handling standard output redirection differently than zap, etc. The cost of trying **zap** over **logrus** was high and the outcome was not that much a blast. All loggers are differents and no-one can expect one to work like another.
Well, why the last sentence should be true ?

From that question is born this project which:

-   expose a unique abstract way of logging through a unique interface
-   expose a unique configuration
-   expose a unique way of redirect standard output
-   expose a unique way of writing to an io.Writer
-   is modulable (can use already built zap or logrus instances, or any other logger)
-   is easily mockable
-   help to test logs

## Example

```go
// var log = <some way of building one>

// same interface we are used to see
log.WithField(key, value).Info(args...)

// easy way to get a io.Writer to inject in every components that require a logger
logger.WriterLevel(log, logger.LevelError)

// easy way of redirecting golang log standard library
logger.RedirectStdLog(log, logger.LevelWarn)
```

Changing the underlying logger takes literraly **one** line.

```go
package main

import (
    "github.com/krostar/logger"
    "github.com/krostar/logger/logrus"
    "github.com/krostar/logger/zap"
)

func main() {
    var config = logger.Config{
        Formatter: "json",
    }

    // create a logrus-based logger with configuration
    log := logrus.New(logrus.WithConfig(config))
    log.Info("i'm a logrus-based logger")

    // as any logger keep the same interface
    // switching to a zap-based logger with configuration is easy
    log = zap.New(zap.WithConfig(config))
    log.Info("i'm a zap-based logger")
}
```

You have a lot of code that:

-   is using io.Writer or 
-   is logging directly with standard library `log`
-   takes a `*log.Logger` parameter

There is actually few methods to help with that:

```go
func setupHTTP(addr string, logger *log.Logger) *http.Server {
    return &http.Server{
		Addr:     addr,
		ErrorLog: logger,
	}
}

func doSomethingElse(doIt func() error, errWriter io.Writer) {
    if err := doIt(); err != nil {
        io.WriteString(errWriter, err.String())
    }
}

func main() {
    var (
        log         logger.Logger // init it the way you want (zap, logrus, ...)
        stdlog    = logger.StdLog(log, logger.LevelError)
        writerlog = logger.WriterLevel(log, logger.LevelError)
    )

    setupHTTP(":80", stdlog)

    doSomethingElse(func()error{
        return errors.New("bim bam boum")
    }, stdlog)
}
```

## License

This project is under the MIT licence, please see the LICENCE file.
