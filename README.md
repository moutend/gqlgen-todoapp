gqlgen-todoapp
==============

Demo app using `gqlgen`.

## Requirements

- golang
- docker
- sqlboiler
- migrate
- mysql

```console
# Install sqlboiler CLI
go get -u -t github.com/volatiletech/sqlboiler
go get github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql

# Install migrate CLI
go get -u -d github.com/golang-migrate/migrate/cmd/migrate
cd $GOPATH/src/github.com/golang-migrate/migrate/cmd/migrate
git checkout $TAG  # e.g. v4.1.0
go build -tags 'postgres' -ldflags="-X main.Version=$(git describe --tags)" -o $GOPATH/bin/migrate github.com/golang-migrate/migrate/cmd/migrate
```

## Build

```console
git clone https://github.com/moutend/gqlgen-todoapp
cd ./gqlgen-todoapp

./bin/setup.bash
./bin/up.bash
./bin/sqlboiler.bash

cd ./cmd/todoapp
go build
go install
```

## Test

Open a terminal and then run the command:

```console
todoapp -c config/common.toml
```

Now the GraphQL playground is available. Open the browser, go to localhost:8080.

Or, open another terminal and then run the command:

```console
./bin/create_user.bash
./bin/create_task.bash
./bin/tasks.bash
```

## LICENSE

MIT
