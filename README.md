# Casbin Redis Watcher

Casbin Redis Watcher is a Go library meant to be used as a [Casbin Watcher](https://casbin.org/docs/en/watchers).

The main difference between this watcher and the currently recommended [Redis Watcher](https://github.com/billcobbler/casbin-redis-watcher) is that this one supports using a Redis Cluster.


## Installation

```sh
go get github.com/escaletech/casbin-redis-watcher/watcher
```

## Usage

The simplest usage would be:

```go
package main

import (
    "github.com/casbin/casbin/v2"
    "github.com/escaletech/casbin-redis-watcher/watcher"
)

func main() {
    // Initialize the watcher
    w, err := watcher.New(watcher.Options{
        RedisURL: "redis://127.0.0.1:6379",
    })
    if err != nil {
        panic(err)
    }

    // Create your enforcer
    e := casbin.NewEnforcer("rbac_model.conf", "rbac_policy.csv")

    // Tell the enforcer to use the watcher
    e.SetWatcher(w)

    // Now whenever e.SavePolicy() is called, the watcher is going to be notified
}
```

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
