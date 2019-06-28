# gomilk

gomilk is a CLI tool for managing your tasks, written in Go!

# Getting Started

## Fetching Dependencies

*Note: at some point I will figure out how to do this automatically.*

Run the following commands to set up dependencies for yourself.

```
$ go get github.com/pkg/browser
$ go get github.com/sirupsen/logrus
$ go get github.com/urfave/cli
```

## Fetching Source Code

Cool! Once you're done with that, grab the code using this command.

```
$ go get github.com/dishbreak/gomilk
```

It will fail. I'm very sorry. Here's why...

## Setting up API Tokens

You'll need to [apply with RTM](https://www.rememberthemilk.com/services/api/keys.rtm) for an API token. Once you have one, you can create a file under `api/secrets.go`, with contents like so.

```
package api

const (
	// APIKey is the API key provided to us by remember the milk
	APIKey = "YOUR API KEY"
	// SharedSecret is the secret used to sign API requests
	SharedSecret = "YOUR SHARED SECRET"
)
```

The terms of the RTM API key require me to protect my app's credentials. Sorry!

Once that's done, you can use `go build && go install` to build the binary and install it on your path.

# Using gomilk

## Logging In

The first thing you'll want to do is login.

```
$ gomilk login
```

This will open a browser window and ask you to login and authorize the application. Once you're done, come back to the terminal and press `[Enter]`.

## Adding Tasks

Probably the thing you'll do the most is add tasks. Gomilk lets you use [Smart Add syntax](https://www.rememberthemilk.com/help/answer/basics-smartadd-howdoiuse) to add your tasks.

```
$ gomilk add "fix the bug tomorrow #gomilk" 
Created task 'fix the bug'
```

## Listing Tasks

The `ls` command will list your tasks.

```
$ gomilk ls
[0] Fix the flux capacitor (due: Today)
[1] Hide the macguffin (due: Today)
[2] Order Photos for minon of the month (due: Today)
[3] fix the bug (due: Tomorrow)
[4] Call back mad scientist (due: Jul 11)
[5] Review doom ray plans (due: Jul 17)
[6] Pay pet license for Mr. Fluffles (due: Jul 20)
[7] Rule the world! (due: Jul 25)
```

If you pass a [search query](https://www.rememberthemilk.com/help/answer/basics-search-advanced), `ls` will show only matching results.

```
$ gomilk ls tag:gomilk
[0] fix the bug (due: Tomorrow)
```

Those numbers are important. Let's see how.

## Completing Tasks

The `complete` command will complete a task that you've listed. 

```
$ gomilk ls tag:gomilk
[0] fix the bug (due: Tomorrow)
$ gomilk complete 0
Completed task 'fix the bug'
```
You can specify ranges, like so:

```
$ gomilk ls
[0] Fix the flux capacitor (due: Today)
[1] Hide the macguffin (due: Today)
[2] Order Photos for minon of the month (due: Today)
[3] fix the bug (due: Tomorrow)
[4] Call back mad scientist (due: Jul 11)
[5] Review doom ray plans (due: Jul 17)
[6] Pay pet license for Mr. Fluffles (due: Jul 20)
[7] Rule the world! (due: Jul 25)
$ gomilk complete 0 1-3 7
```
