Creating a route in Go
Using templates
Building database connections
Collecting data
Using web middleware
Using the Ace template engine
Integrating HTTP routers like gorilla/mux
Authenticating users
Optimizing a Go codebase


```

```

Example DB setup
```
ngans-mbp:Go-Library ngan$ sqlite3 dev.db
SQLite version 3.19.3 2017-06-27 16:48:08
Enter ".help" for usage hints.
sqlite> .schema
sqlite> create table books(
   ...> pk integer primary key autoincrement,
   ...> title text,
   ...> author text,
   ...> id text,
   ...> classification text
   ...> );
sqlite> .schema
CREATE TABLE books(
pk integer primary key autoincrement,
title text,
author text,
id text,
classification text
);
CREATE TABLE sqlite_sequence(name,seq);
sqlite> 

```