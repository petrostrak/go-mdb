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

To add the indexes:
```bash
migrate create -seq -ext .sql -dir ./migrations add_movies_indexes
```

To add users table:
```bash
migrate create -seq -ext=.sql -dir=./migrations create_users_table
```

To add tokens table
```bash
migrate create -seq -ext .sql -dir ./migrations create_tokens_table
```

To add permissions table:
```bash
migrate create -seq -ext .sql -dir ./migrations add_permissions
```

To execute the migrations:
```bash
migrate -path=./migrations -database=$GOMDB_DSN up
```

In case of an error `error: Dirty database version x. Fix and force version.`
```bash
migrate -path=./migrations -database=$GOMDB_DSN force x
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
A `reductive filter` which allows clients to search based on a case- insensitive exact match for movie title and/or one or more movie genres. For example:
```go
// List all movies.
/v1/movies
// List movies where the title is a case-insensitive exact match for 'black panther'.
/v1/movies?title=black+panther
// List movies where the genres includes 'adventure'.
/v1/movies?genres=adventure
// List movies where the title is a case-insensitive exact match for 'moana' AND the // genres include both 'animation' AND 'adventure'. /v1/movies?title=moana&genres=animation,adventure
```

* The `page` value is between 1 and 10,000,000.
* The `page_size` value is between 1 and 100.
* The `sort` parameter contains a known and supported value for our movies table. Specifically, we’ll allow `"id"`, `"title"`, `"year"`, `"runtime"`, `"-id"`, `"-title"`, `"-year"` or `"-runtime"`.
<sub><sup>The `-` character to denotes descending sort order.<sub><sup>

### CORS
To pass an arbitrary list (space separated) of URIs as trusted origins:
```bash
go run ./cmd/api -cors-trusted-origins="http://localhost:9000 http://localhost:9001"
```