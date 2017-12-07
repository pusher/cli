# Here be dragons! Alpha-quality software! It may eat your dog!

# Pusher CLI [(pusher.com)](https://pusher.com)

This is a tool that allows developers access to their Pusher accounts via a command line interface. 

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
1. Copy `pusher` to your `$GOPATH/bin` or just use it as is.

### Homebrew

You can install this package via Homebrew by pasting the following into a terminal.

```
brew install pusher/brew/pusher
```

## Hacking On It

1. Clone this repository;
1. Create a new branch by running `git checkout -b <YOUR_BRANCH_NAME_HERE> master`
1. Pull dependencies with `dep ensure`;
1. Ready to hack.
