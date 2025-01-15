# WheresMyLift

The project aims to consolodate all Irish public transport operators that provide realtime information on where their vehicles into one website/API.

## Packages

### API

Source is located in [api](api). Information on this package can be found in the README under that directory.

### Web-UI

Source is located in [web-ui](web-ui). Information on this package can be found in the README under that directory.

## Contributing

### Using Git

There is a git hook that has been set up to help you create messaging and won't let you commit unless you do it right. All you need to do is execute the following script once after you clone the repo:
```
    .githooks/git-hooks-config.sh
```

### Conventional Commits
This repo makes use of the 'conventional commits' convention. This means that all commits should be formatted using this convention. 

For example:
    feat(router): added events route

An example of an incorrect commit message would be:
    i fixed an error that caused the events route to not work lol

All commit message should follow the follow standard:
```
    <type>: <subject>
        or
    <type(<scope>):> <subject>
```

A type can be: build, chore, ci, docs, feat, fix, perf, refactor, revert, style, test, or wip.
A scope should encompass the change in one to two words. It would usually refer to the location of the change, e.x. the controller, model, or router.
The subject should be limited to 50 characters and can describe your change. If you need to add more detail to the commit message, you can use a line break (\n).

## Why am I here?
I don't know. While you're here though, check out this [cool video](https://www.youtube.com/watch?v=dQw4w9WgXcQ).

![XKCD Santa Sudo Meme](https://imgs.xkcd.com/comics/incident.png "He sees you when you're sleeping, he knows when you're awake, he's copied on /var/spool/mail/root, so be good for goodness' sake.")