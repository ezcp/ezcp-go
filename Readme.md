# What's EZCP

*EZCP works like a clipboard*: you put one thing in it from machine A, and you can get it back from machine B.

It works accross the internet, so you can put something in it on one machine,
and get it back on another machine. Easy Copy. Hence the name.

To make sure your clipboard is private, it is associated to a Token.
You can create a new token by visiting [https://ezcp.io](https://ezcp.io).
Copy the token showing on the page, and you're good to go (others get a different token).

## Install

There are two ways to get the EZCP command line tool: Go, and Node.js.

Some prefer a single binary compiled for their platform, others like Node.js tools better.

### You can install the following *Go* binaries

- for linux, [https://get.ezcp.io/linux/ezcp](https://get.ezcp.io/linux/ezcp)
- for osx, [https://get.ezcp.io/osx/ezcp](https://get.ezcp.io/osx/ezcp)
- for windows, [https://get.ezcp.io/windows/ezcp](https://get.ezcp.io/windows/ezcp)

Open the link in a navigator, or use `wget` or `curl`, then copy the file somewhere 
in your path.

### Or the Node.js version

The Node.js install is even easier if you already have Node.js installed.

`npm install -g ezcp`

Then you can use the `ezcp` command.

# Starting

Say you have a file you want to send from machine A to machine B.

On machine A, type `ezcp [filename] [token]`. It copies the file contents to the internet.

On machine B, simply type `ezcp [token] [filename]` and ezcp will re-create the file.

Of course, both machine A and machine B must have the internet.

*Note:*

- Once you've put something in a Token, you can't reuse the token to put something else inside.
- Once you've gottent something out of a Token, you can't reuse it to get the file again.
