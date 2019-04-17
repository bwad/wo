# wo

`wo` is a rudimentary tool for running tests when watched files change.

The project was motivated by a slight addiction to my `guard-rspec` workflow.
It also seemed like a reasonable sized first golang project. 

## Watched Files

Currenly `wo` watches the `go` files in the directory in which it's launched
(i.e. current working directory).

## Test Files

When a watched file changes (.go or \_test.go), we run the associated 
test file with the following command:

    go test <file>.go <file>_test.go

If there is no test file, nothing happens.

## Future 

This project scratches a personal itch and as such will most likely evolve
as I find other places that I can't reach with standard development tools.


