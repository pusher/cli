# Pusher CLI [(pusher.com)](https://pusher.com)

This is a tool that allows developers access to their Pusher accounts via a command line interface. 

This is an alpha release. We recommend using it for non-production apps. It may eat your laundry! We'd appreciate your feedback.

## Usage

Before attempting to use the Pusher CLI, You should log into your [dashboard](https://dashboard.pusher.com/accounts/edit) and generate a new API key. Next, follow the installation instructions [below](#installing), and then run `pusher login`.

## Installing

There's multiple ways you can get the Pusher CLI onto your machine:

### Downloading the Binary

You can download the latest release from [here](https://github.com/pusher/cli/releases) and add it to your path.

### Building from Source

1. Clone this repository;
1. Pull dependencies with `dep ensure`;
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
1. Run `go test` to fetch dependencies and run tests for the first time.
1. Ready to hack.

We [publish binaries on GitHub](https://github.com/pusher/cli/releases)

1. Get [fpm](https://github.com/jordansissel/fpm)
1. Get [goreleaser](https://goreleaser.com/)
1. Get `rpmbuild`, e.g. `brew install rpm` on MacOS
1. [Generate a GitHub personal access token](https://github.com/settings/tokens)
   with the `repo` scope selected.
   Set this as env var `GITHUB_TOKEN`.
1. From this directory, run `goreleaser` (or [follow these instructions](https://goreleaser.com/#releasing))
