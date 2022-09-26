# How to run locally in Linux?

## Install Golang 1.15

Link to Golang 1.15 Linux [archive](https://go.dev/dl/go1.15.linux-amd64.tar.gz).

Follow instructions as mentioned in step 2 for Linux here https://go.dev/doc/install.

## Install Postgresql 9.6

Follow instructions here to install postgresql 9.6 and create a database with name `savethislife`.

## Run the site locally

From the top folder, run this command:

```bash
go run ./server/main.go -db-info="postgres://<POSTGRES_USER>:<POSTGRES_PASS>@localhost/savethislife?sslmode=disable"
```

You should see an output similar to this:
```text
server listening at :8080 on http...
```

Now, go the browser and open the URL `http://localhost:8080`to see the site locally.

Whenever you do any changes to the site, you need to stop the command you ran before and run it again.