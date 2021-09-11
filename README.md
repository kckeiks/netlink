# Netlink

Netlink implementation in GO. 

My motivation for creating this library was that I could not find one that had the following properties:

* Does not use `unsafe` package.
* Does not use deprecated `syscall` package.
* Allows control for setting endian order.
