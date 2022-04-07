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

Lists stale branches that have matching prefix:

```
git stale hotfix/ feature/
```

Show local abandand branches at least 3 months passed, then do `git push --delete` for them:

```
git stale --since 3mo
git stale --since 3mo --delete --push
```

## Options

```
git stale [prefix...]
git stale -d [prefix...]
git stale -d -f [prefix...]
git stale --push -d [-f] [prefix...]
```

- `-d`, `--delete`: remove selected branches.
- `-f`, `--force`: combined with `-d`, remove branches even if it wasn't merged.
- `--since <date>`: select branches which have older last commit date, instead selecting gone branches. Check out relative date format section.
- `--push`: combined with `-d`, invokes `git push --delete` instead removing local branches.
  `--force` also be passed to `git push` if `-f` is specified.

As default, this command selects "gone" branches.

### Relative Date Format

Some option in `git-stale` can accept relative date.

- "1mo 2days" will be 1 month and 2 days.
- "3y4w" will be 3 years and 1 month.

Syntax in BNF is roughly described as below:

```
<period> ::= [<digits> <year-suffix>] [<digits> <month-suffix>] [<digits> <week-suffix>] [<digits> <day-suffix>]
<year-suffix> ::= "y" | "yr" | "yrs" | "year" | "years"
<month-suffix> ::= "mo" | "month" | "months"
<week-suffix> ::= "w" | "week" | "weeks"
<day-suffix> ::= "d" | "day" | "days"
```
