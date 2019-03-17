# zap

Using zap instance to build a `logger.Logger`

```go
// there are few ways to build a zap instance

// using the logger configuration
var log, flush, err = zap.New(
    zap.WithConfig(cfg logger.Config),
)

// building it directly
var log, flush, err = zap.New(
    logrus.WithLevel(level logger.Level),
    logrus.WithConsoleFormatter(colored bool),
    logrus.WithJSONFormatter(),
    logrus.WithOutputPaths(output []string),
)

// or by giving an original zap.Config
var log, flush, err = zap.New(
    logrus.WithZapConfig(cfg zap.Config),
)
```

Once the logger has been built, it can be used like any other logger.Logger.
