# Scriptprox2
Basically mitmproxy for FiveM servers

This was one of the first Go projects I've ever worked on & is definitely one I've learnt a lot from. I decided to publish this because no one is really using this anymore & maybe someone will need this/learn from this :)

## Usage
1. Build the project with `go build` or `make debug`
	- Either works, I just wanted a kind of complete Makefile for this project
2. (optional) Create your own resource.rpf
	- This entire project revolves around injecting a specific resource & [a sample resource](https://git.patriick.dev/git/patrick/modmod) is already included (see resource.rpf here)
	- The easiest way to get your own is to setup your own FiveM server, load the resource you're trying to use & going to `server-data\cache\files\<name of the resource>` and using that resource.rpf
3. Connect to a server using the built-in server list