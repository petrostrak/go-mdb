## go_mDB
A JSON API for retrieving and managing information about movies writen in Go.

To run the application on port 3030 in production mode
```bash
go run ./cmd/api -port 3030 -env=production
```

### Decoupling the DSN
Create a new GOMDB_DSN environment variable by adding the following line to either your `$HOME/.profile` or `$HOME/.bashrc` files:
```bash
export GOMDB_DS='postgres://postgres:password@localhost/gomdb?sslmode=disable'
```