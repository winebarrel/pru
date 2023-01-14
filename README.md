# pru

pru is a CLI that updates the pull request branch that contains the specified file.

cf. https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#update-a-pull-request-branch

## Usage

```
Usage: pru [OPTION] OWNER/REPO PATTERNS...
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

# Installation

```
brew install winebarrel/pru/pru
```
