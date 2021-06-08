# git-stale -- lists stale branches and remove

## Getting Started

```
go install github.com/oakcask/git-stale/cmd/git-stale@latest
```

## Options

```
git stale
git stale -d
git stale -d -f
```

* `-d`, `--delete`: removes gone branches.
* `-f`, `--force`: combined with `-d`, remove branches even if it wasn't merged.
