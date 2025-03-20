# gwh

Local data warehouse built using data from multiple Git repositories.

## CLI

```sh
cd path/to/project_x

# Create an emtpy data warehouse. All files will be stored into `$PWD/.gwh`.
gwh init

# Link a couple of Git repositories. This will not import any data yet.
gwh link git_repository_1
gwh link subdir/git_repository_2

# Import changes from all repositories incrementally.
gwh sync
```

## Go Library

```sh
go get https://github.com/npclaudiu/gwh
```

```go
import (
	"fmt"

    "github.com/npclaudiu/gwh/v1"
)

func demo() error {
    warehouse, err := gwh.Open("path/to/warehouse/store")

    if err != nil {
        return fmt.Errorf("failed to open warehouse: %w", err)
    }

    defer warehouse.Close()

    err := warehouse.LinkRepository("path/to/git_repository")

    if err != nil {
        return fmt.Errorf("failed to link repository: %w", err)
    }

    return nil
}
```
