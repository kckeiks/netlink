# Netlink

This project provides the building blocks for linux netlink. Please see [RFC359](https://datatracker.ietf.org/doc/html/rfc3549) and [netlink(7)-manpage](https://man7.org/linux/man-pages/man7/netlink.7.html).

My motivation for creating this library was that I could not find one that had the following properties:

* Does not directly depend on the `unsafe` package and the deprecated `syscall` package.
* Allows users to explicitly set the endian order.

## Features

* The `netlink` package provides some serializers, constructors and other parsers for handling netlink messages. 

* The `sock_diag` package has serializers and request constructors for the netlink family, `NETLINK_SOCK_DIAG`. Please see [sock_diag](https://man7.org/linux/man-pages/man7/sock_diag.7.html).

## Example

In the playground directory, you can find different examples on how these packages are used.

## Work In Progress

Most of the heavy lifting has been done and the package `sock_diag` is close to been complete. The most important thing left to do is to implement some proper error handling.

