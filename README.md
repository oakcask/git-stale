# git-stale -- list stale branches and remove

## Getting Started

```
go install github.com/oakcask/git-stale/cmd/git-stale@latest
```

### Cheat Sheet

Cleaning up stale remote reference in local repo, then remove those stale branches:

```
git prune origin
git stale -d
```

## Options

```
git stale
git stale -d
git stale -d -f
```

- `-d`, `--delete`: remove gone branches.
- `-f`, `--force`: combined with `-d`, remove branches even if it wasn't merged.
