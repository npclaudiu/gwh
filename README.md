# gwh

Local data warehouse built using data from multiple Git repositories.

## Introduction

A warehouse is a local directory that stores one DuckDB database for each Git
repository that is linked to it. All these databases are orchestrated by an
SQLite control database.

## CLI

```sh
cd path/to/project_x

# Create an emtpy data warehouse. All files will be stored into `$PWD/.gwh`.
gwh init

# Link a couple of Git repositories. This will not import any data yet.
gwh link repo1 git_repository_1
gwh link repo2 subdir/git_repository_2 # --recursive

# Import changes from the specified repository incrementally.
gwh sync repo1
```

## Go Library

```sh
go get https://github.com/npclaudiu/gwh/v1
```

```go
import (
    "fmt"

    "github.com/npclaudiu/gwh/v1"
)

func demo() error {
    // Open warehouse directory named `.gwh`, creating it if it does not exist.
    //
    warehouse, err := gwh.Open("path/to/warehouse/store")

    if err != nil {
        return fmt.Errorf("failed to open warehouse: %w", err)
    }

    // Close handles at exit.
    //
    defer warehouse.Close()

    // Link Git repository at the specified path. This will not cause any
    // data synchronization at this point.
    //
    err := warehouse.LinkRepository("repo1", "path/to/git_repository")

    if err != nil {
        return fmt.Errorf("failed to link repository: %w", err)
    }

    // Pull data from the repository incrementally.
    //
    err := warehouse.SyncRepository("repo1")

    if err != nil {
        return fmt.Errorf("failed to pull data from repositories: %w", err)
    }

    return nil
}
```
