# logmid

`logmid` is a **http middleware** that log using a `logger.Logger`.

By default no fields (like _latencies_, _status_, ...) are added to the logger
as it's up to the developer to know which fields he / she may want. Instead, a
callback can be given to the construtor to add any fields to the logger (or
do any other kind of action).

Inside a `http.Handler` any fields can be added using `AddFieldInContext` and / or
`AddErrorInContext` which respectively call `logger.WithField` and `logger.WithError`.

Custom options can be applied to the middleware (for example the verbosity of the log
based on whatever please you, the message wrote, ...)

## Example

```go
func setupRouter(router) {
    router.Use(
        httpinfo.Record(), // middleware to record latencies, status, ...

        logmid.New(
            logmid.WithDefaultFields(), // adds default fields like latencies, status, ...
            logmid.WithLogLevelFunc(customLevelFunc), // default is Info, but custom level can be applied
            logmid.WithCallback(customCallback),
        ),
    )
}

func customLevelFunc(r *http.Request) logger.Level {
    var lvl logger.Level

    switch status := httpinfo.Status(); status {
    case status >= 400 && status < 500:
        lvl = logger.LevelWarn
    case status >= 500:
        lvl = logger.LevelError
    default:
        lvl = logger.LevelDebug
    }
    
    return lvl
}

func customCallback(r *http.Request) {
    logmid.AddFieldInContext(r.Context(), "hello", "world")
}
```
