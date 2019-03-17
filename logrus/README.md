# logrus

Using logrus instance to build a `logger.Logger`

```go
// there are few ways to build a logrus instance

// using the logger configuration
var log, err = logrus.New(
    logrus.WithConfig(cfg logger.Config),
)

// building it directly
var log, err = logrus.New(
    logrus.WithLevel(level logger.Level),
    logrus.WithConsoleFormatter(colored bool),
    logrus.WithJSONFormatter(),
    logrus.WithOutput(writer io.Writer),
)

// or by giving an already built logrus instance
var log, err = logrus.New(
    logrus.WithInstance(log *logrus.Logger),
)
```

Once the logger has been built, it can be used like any other logger.Logger.
