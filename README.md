# Ingram Micro Senior Backend Developer test

This task challenges you to implement in Go an in-memory queue satisfying a specific interface. Our intention behind this task is to have something to talk about in an interview, so it does not matter if you get stuck and cannot complete it: just deliver whatever you got.

## Setup

- [Install Go](https://golang.org/doc/install#install) if you do not have it already
- Clone this repository outside your GOPATH (because it is using [go modules](https://blog.golang.org/using-go-modules))

## Develop

You should write your code in the queue package, modifying the worker.go file to complete the New method and maybe in other .go files. You should probably leave the interfaces.go file as it is.

As specified below, the delivery format is a git repository. It would be nice that the git history showed your progress.

### Run your code

The main.go file in the top-level directory of this repo will use your implementation to approximate pi. Before delivering your task you should try your best to make it work. You can run it with:

```Bash
go run main.go
```

### Tips

* Channels are not a silver bullet in go
* You will probably need to get familiar with the [context package](https://golang.org/pkg/context/). The following blog post can be an introductory read: http://p.agnihotry.com/post/understanding_the_context_package_in_golang/

### Bonus tracks

If you really want to impress us by going the extra mile you could do one of the following:

- Add tests for your code
- Dockerize the running of the main.go file
- Implement a persisted (file, redis or something else -based) implementation of the queue, maybe in another package, that satisfies the interface.

## Deliver

Create a public github repository and push your code there. Then send us an email with a link to it.
