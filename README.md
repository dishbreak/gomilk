# gomilk

[![Build Status](https://travis-ci.org/dishbreak/gomilk.svg?branch=master)](https://travis-ci.org/dishbreak/gomilk)

gomilk is a CLI tool for managing your tasks, written in Go!

# Getting Started

So you've cloned the repo. Cool! Here's what you need to do first.

## Setting up API Tokens

You'll need to [apply with RTM](https://www.rememberthemilk.com/services/api/keys.rtm) for an API token. You'll need to configure credentials in your environment like so:

```
GOMILK_API_KEY=yourapikeyhere
GOMILK_SHARED_SECRET=yoursharedsecret

The terms of the RTM API key require me to protect my app's credentials. Sorry!

## Building and Installing

Once that's done, you can use `make` to build the binary and `make install` to install it on your path. Whenever you make changes and you want to try them out, run `make install`.

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
