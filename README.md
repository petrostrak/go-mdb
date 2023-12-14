## go_mDB
A JSON API for retrieving and managing information about movies writen in Go.

To run the application on port 3030 in production mode
```bash
go run ./cmd/api -port 3030 -env=production
```

### Decoupling the DSN
Create a new GOMDB_DSN environment variable by adding the following line to either your `$HOME/.profile` or `$HOME/.bashrc` files:
```bash
export GOMDB_DSN='postgres://postgres:password@localhost/gomdb?sslmode=disable'
```
### Working with SQL Migrations
```bash
brew install golang-migrate
```
To create the first migration:
```bash
migrate create -seq -ext=.sql -dir=./migrations create_movies_table
```
In this command:
* The `-seq` flag indicates that we want to use sequential numbering instead of Unix timestamps.
* The `-ext` flag indicates that we want to give the migration files the extension `.sql`.
* The `-dir` flag indicates that we want to store the migration files in the `./migrations`
directory (which will be created automatically if it doesn’t already exist).
* The name `create_movies_table` is a descriptive label that we give the migration files to signify their contents.

To add the constraints:
```bash
migrate create -seq -ext=.sql -dir=./migrations add_movies_check_constraints
```

To execute the migrations:
```bash
migrate -path=./migrations -database=$GOMDB_DSN up
```

Migrating to a specific version

As an alternative to looking at the schema_migrations table, if you want to see which migration version your database is currently on you can run the migrate tool’s version command, like so:
```bash
migrate -path=./migrations -database=$GOMDB_DSN version
```
You can also migrate up or down to a specific version by using the goto command:
```bash
migrate -path=./migrations -database=$GOMDB_DSN goto 1
```

### Filtering, Sorting and Pagination
The page, page_size and sort query string parameters in action:
```bash
curl "localhost:8080/v1/movies?title=godfather&genres=crime,drama&page=1&page_size=5&sort=year"
```