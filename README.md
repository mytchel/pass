# Pass
## The Unix Password Manager

Pass Uses chained aes encryption to store passwords (or small notes) in a file
that allows for easy modification and viewing of the stored passwords.

Passwords can be organised into directory trees like in a file system.

## Usage

`go get -u github.com/mytchel/pass` to install then run `pass`.

The first time it is run it will prompt for a password and initialise 
a new encrypted file at `$HOME/.secstore`. After that it will drop to 
a repl and allow you to add, edit, delete, move, and rename passwords.

You can also give the commands as arguments eg: 
```
pass show a-password
```

Is the same as running `pass` then `show a-password`

#### Requirements

Go and github.com/peterh/liner for repl support.

## Origins

This all came about due to me switching to OpenBSD and discovering that 
[the standard unix password manager](http://www.passwordstore.org/) didn't work.
Portability didn't seem to be a concern. I remade the script to work much more
simpily. Then I decided that it was kind of strange needing to have public and 
private keys just to encrypt some files with a password. That and my friend told 
me it was fucking weird to be using gpg for that. 

Now you have this.
