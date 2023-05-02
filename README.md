# Pusher CLI [(pusher.com)](https://pusher.com)

This is a tool that allows developers access to their Pusher accounts via a command line interface. 

This is a beta release. We recommend using it for non-production apps unless otherwise advised. We'd appreciate your feedback.

## Usage

Before attempting to use the Pusher CLI, You should log into your [dashboard](https://dashboard.pusher.com/accounts/edit) and generate a new API key. Next, follow the installation instructions [below](#installing), and then run `pusher login`.

## Installing

There's multiple ways you can get the Pusher CLI onto your machine:

### Downloading the Binary

You can download the latest release from [here](https://github.com/pusher/cli/releases) and add it to your path.

### Building from Source

1. Clone this repository;
1. Build with `go build -o pusher`;
1. Copy `pusher` to somewhere in your `$PATH` or just use it as is.

### Homebrew

You can install this package via Homebrew by pasting the following into a terminal.

```
brew install pusher/brew/pusher
```

## Hacking On It

1. Clone this repository;
1. Create a new branch by running `git checkout -b <YOUR_BRANCH_NAME_HERE> master`
1. Run `go build` to fetch dependencies and run tests for the first time.
1. Ready to hack.

We [publish binaries on GitHub](https://github.com/pusher/cli/releases) and we use a github action to release for multiple platforms. To create a release just tag

1. `git tag -a v0.14 -m "v0.14"`
1. `git push origin v0.14`

### Configuration

`pusher login` creates a file `~/.config/pusher.json` (or updates it if it already exists).
If you need to point the Pusher CLI to different servers (e.g. when testing), you can change the `endpoint` value and add new name/value pairs as necessary:
```JSON
{
  "endpoint": "https://cli.another.domain.com",
  "token": "my-secret-api-key",
  "apihost": "api-mycluster.another.domain.com",
  "httphost": "sockjs-mycluster.another.domain.com",
  "wshost": "ws-mycluster.another.domain.com"
}
```
