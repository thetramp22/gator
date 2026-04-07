# gator

gator is a command line blog aggregator. It is a program written in Go using Postrgres.  This project is a guided project from Boot.Dev and has been created as a learning project.

## Requirements

To build or install this program you will need Postgres and Go installed locally.

## Installation

Once Postgres and Go are installed, gator can be installed by running the command:

```go install github.com/thetramp22/gator@latest```

## Usage

```gator register <user>```

Register a user.

```gator login <user>```

Login to gator as a registered user.

```gator addfeed <name> <url>```

Add a new RSS feed to your database of feeds.  Adding a feed also follows the feed.

```gator feeds```

Lists all feeds added by any user.

```gator follow <url>``` 

Follows a feed added by another user.

```gator unfollow <url>```

Unfollows a feed.

```gator following```

List all feeds followed by currently logged in user.

```gator agg <interval>```

Runs the gator aggregator.  This process runs continually in the background, pulling in posts from any followed feed at the given interval of time, and saving them to your database of feeds.  Ctrl + C to end the process.
