# Codecov

[![codecov](https://codecov.io/gh/taiidani/go-codecov/branch/master/graph/badge.svg)](https://codecov.io/gh/taiidani/go-codecov)

This library is made to access the [Codecov](https://codecov.io/) REST API, allowing programmatic management of the resources exposed through it.

## Using The Library

First, generate a Codecov API token for your user at `https://codecov.io/account/gh/<user>/access`. You will need it in order to authenticate against the API.

Bring the library into your project with:

```sh
go get github.com/taiidani/go-codecov
```

And that's it!

## Examples

Listing all repositories for your user:

```go
c := codecov.NewClient("<token>")
repos, err := c.ListRepositories(context.Background(), "<username/org>")
if err != nil {
    panic(err)
}

for _, repo := range repos {
    fmt.Println("%#v", repo)
}
```

Extracting data for a single repository:

```go
c := codecov.NewClient("<token>")
repo, err := c.GetRepository(context.Background(), "<username/org>", "<repo>")
if err != nil {
    panic(err)
}

fmt.Println("%#v", repo)
```
