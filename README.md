# pru

pru is a tool that updates pull requests branch from the base branch that contains specified files.

cf. https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#update-a-pull-request-branch

```mermaid
gitGraph
  commit
  commit
  branch pr/foo
  commit
  commit
  commit
  checkout main
  branch pr/bar
  commit
  commit
  commit
  checkout main
  merge pr/foo tag:"merge pull request"
  commit
  checkout pr/bar
  merge main tag:"auto update by pru"
  checkout main
  commit
```

## Usage

```
Usage: pru [OPTION] OWNER/REPO PATTERNS...
  -bases value
    	base branches to update (default "main,master")
  -dry-run
    	dry run
  -ignore-labels value
    	labels for pull requests that do not update
  -token string
    	GitHub access token. use $GITHUB_TOKEN env
  -version
    	print version
```

```
$ pru my-owner/my-repo '**/*.go'
update https://github.com/my-owner/my-repo/pull/123
update https://github.com/my-owner/my-repo/pull/125
update https://github.com/my-owner/my-repo/pull/127
...
```

## Installation

```
brew install winebarrel/pru/pru
```

## GitHub action

see https://github.com/winebarrel/pru-action

### Example

see https://github.com/winebarrel/pru-example
