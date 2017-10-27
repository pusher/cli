# Pusher CLI

This is a tool that allows developers access to their Pusher accounts via a command line interface. 

## Usage

Before attempting to use the Pusher CLI, You should log into your [dashboard](https://dashboard.pusher.com/accounts/edit) 
and generate a new API key. Next, follow the setup instructions below, and then run `pusher login`.

```
pusher [command]

Available Commands:
  apps        Manage your Pusher Apps
  generate    Generate a Pusher client, server, or Authorisation server
  help        Help about any command
  login       Enter and store Pusher account credentials
  logout      Remove Pusher account credentials from this computer
```

## Installing

There's multiple ways you can get the Pusher CLI onto your machine:

### Building from Source

1. Clone this repository;
1. Pull dependencies with `govendor sync`;
1. Build with `go build -o pusher`;
1. Copy `pusher` to your `$GOPATH/bin` or just use it as is.

### Downloading the Binary

You can download the latest release from [here](https://github.com/pusher/cli/releases) and add it to your path.

### Homebrew

You can install this package via Homebrew by pasting the following into a terminal.

```
brew install pusher/brew/pusher
```

## Hacking on it

1. Clone this repository;
1. Create a new branch based on `master` by running `git checkout -b <YOUR_BRANCH_NAME> master`;
1. Pull dependencies with `govendor sync` - This will modify vendor.json. Don't commit this file;
1. Hack on it.
