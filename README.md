# go-mobile-algorand-sdk

## Summary

This repo makes a subset of the [Go Algorand SDK](https://github.com/algorand/go-algorand-sdk) available for use as an iOS and Android library. This is achieved using the [Go Mobile](https://pkg.go.dev/golang.org/x/mobile) project.

There are many limitations to what can be exposed through Go Mobile. This repo pulls in the official Go Algorand SDK as a dependency and only exposes certain functions with a limited set of parameters and return types.

## :warning: Guarantees :warning:

This repo and the Go Mobile project are **experimental**. Because of this we cannot make any guarantees about its behavior, future support, or suitability for use in a production system.

## Building

The command `make install-go-mobile` can be used to install the `gomobile` CLI. Then `make ios` and `make android` can be used to build the iOS and Android bindings for this library.

### Dependencies

To run these commands you will need a compatible version of Go (we currently use 1.17), as well as the necessary mobile SDK dependencies. For iOS, this means XCode, and for Android this means the Android SDK and NDK.

As a convenience, this repo has a Github Action which builds the iOS version automatically.
