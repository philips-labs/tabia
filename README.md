# Tabia

[![Go CI](https://github.com/philips-labs/tabia/workflows/Go%20CI/badge.svg)](https://github.com/philips-labs/tabia/actions)
[![codecov](https://codecov.io/gh/philips-labs/tabia/branch/develop/graph/badge.svg?token=K2R9WOXNBm)](https://codecov.io/gh/philips-labs/tabia)

Tabia means characteristic in Swahili. Tabia is giving us insights on the characteristics of our code bases.

## Setup

Copy `.env.example` to `.env` and fill out the credentials. This environment variable is read by the CLI and tests. Also vscode will read the variable when running tests or starting debugger.

```bash
cp .env.example .env
source .env
env | grep TABIA
```

Install required tools:

```bash
make install-tools
```

Add your $GO_PATH/bin folder to your path to make the installed Go tools globally available.

```bash
export GO_PATH="$HOME/go"
export PATH="$GO_PATH/bin:$PATH"
```

## Build

To build the CLI you can make use of the `build` target using `make`.

```bash
make build
```

## Test

To run tests you can make use of the `test` target using `make`.

```bash
make test
```

## Run

### Bitbucket

To interact with Bitbucket `tabia` makes use of the [Bitbucket 1.0 Rest API](https://docs.atlassian.com/bitbucket-server/rest/7.3.0/bitbucket-rest.html).

```bash
bin/tabia bitbucket --help
bin/tabia bitbucket projects --help
bin/tabia bitbucket repositories --help
```

### Github

To interact with Github `tabia` makes use of the [Github graphql API](https://api.github.com/graphql).

```bash
bin/tabia github --help
bin/tabia github repositories --help
```

#### Authentication

Please note when using Github Authentication there are 2 options to authenticate.

1. Authenticate as a Github App (your app will have to be installed in the organization)
   - integration-id
   - private-key
2. Authenticate using a Personal Access Token
   - token

> :warning: When authenticating as a *GitHub App* please be informed you can only fetch information from **one** organization at a time as the client will be bound to that organizations App installation. To support multiple organizations we require a refactor using a Github client per organization.

### Output - Grimoirelab

To expose the repositories in [Grimoirelab projects.json](https://github.com/chaoss/grimoirelab-sirmordred#projectsjson-) format, you can optionally provide a json file to map repositories to projects. By default the project will be mapped to the owner of the repository. Anything not matching the rules will fall back to this default.

E.g.:

```bash
bin/tabia -O philips-labs -M github-projects.json -F grimoirelab > projects.json
```

Regexes should be defined in the [following format](https://golang.org/pkg/regexp/syntax/).

```json
{
  "rules": {
    "One Codebase": { "url": "tabia|varys|garo|^code\\-chars$" },
    "HSDP": { "url": "(?i)hsdp" },
    "iX": { "url": "(?i)ix\\-" },
    "Licensing Entitlement": { "url": "(?i)lem\\-" },
    "Code Signing": { "url": "(?i)^code\\-signing$|notary" }
  }
}
```

#### Output - using template

To generate the output for example in a markdown format you can use the option for a templated output format. This requires you to provide the path to a template file as well. Templates can be defined using the following [template/text package syntax](https://golang.org/pkg/text/template/).

E.g.:

```md markdown.tmpl
# Our repositories

Our repository overview. Private/Internal repositories are marked with a __*__

{{range .}}* [{{ .Name}}]({{ .URL }}) {{if .IsPrivate() }}__*__{{end}}
{{end}}
```

Using above template we can now easily generate a markdown file with this unordered list of repository links.

```bash
bin/tabia -O philips-labs -F templated -T markdown.tmpl > repositories.md
```

#### Filter

##### Repositories

The following repository fields can be filtered on.

* ID
* Name
* Description
* URL
* SSHURL
* Owner
* Visibility
* CreatedAt
* UpdatedAt
* PushedAt
* Topics

The following functions are available.

* `func (RepositoryFilterEnv) Contains(s, substr string) bool`
* `func (Repository) IsPublic() bool`
* `func (Repository) IsInternal() bool`
* `func (Repository) IsPrivate() bool`
* `func (Repository) HasTopic(topic string) bool`
* `func (Repository) CreatedSince(date string) bool`
* `func (Repository) UpdatedSince(date string) bool`
* `func (Repository) PushedSince(date string) bool`

```bash
$ bin/tabia github repositories -O philips-labs -f '{ !.IsPrivate() && !.IsInternal() && !Contains(.Name, "terraform") }'
      Name                                  Owner         Visibility  Clone
0001  helm2cf                               philips-labs  Public      https://github.com/philips-labs/helm2cf
0002  dct-notary-admin                      philips-labs  Public      https://github.com/philips-labs/dct-notary-admin
0003  notary                                philips-labs  Public      https://github.com/philips-labs/notary
0004  about-this-organization               philips-labs  Public      https://github.com/philips-labs/about-this-organization
0005  sonar-scanner-action                  philips-labs  Public      https://github.com/philips-labs/sonar-scanner-action
0006  medical-delivery-drone                philips-labs  Public      https://github.com/philips-labs/medical-delivery-drone
0007  dangerous-dave                        philips-labs  Public      https://github.com/philips-labs/dangerous-dave
0008  varys                                 philips-labs  Public      https://github.com/philips-labs/varys
0009  garo                                  philips-labs  Public      https://github.com/philips-labs/garo
..........
...........
........
```

##### Members

The following member fields can be filtered on.

* ID
* Login
* Name
* Organization
* SamlIdentity
  * ID

The following functions are available.

* `func (MemberFilterEnv) Contains(s, substr string) bool`

#### Download contents

```bash
bin/tabia github contents --repo philips-labs/tabia --file README.md --output downloads/tabia/README.md
$ cat downloads/tabia/README.md
# Tabia

[![Go CI](https://github.com/philips-labs/tabia/workflows/Go%20CI/badge.svg)](https://github.com/philips-labs/tabia/actions)
[![codecov](https://codecov.io/gh/philips-labs/tabia/branch/develop/graph/badge.svg?token=K2R9WOXNBm)](https://codecov.io/gh/philips-labs/tabia)
...
...
..
```
