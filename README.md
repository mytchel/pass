# Pass
## The Unix Password Manager

Uses chained aes encryption to store passwords (or small notes) in a file
that allows for easy modification and viewing of the stored passwords.

Passwords can be organised into directory trees like in a file system.

There is a repl to make multiple changes not require typing your password
half a dozen times.

## origins

This all came about due to me switching to OpenBSD and discovering that 
[the standard unix password manager](http://www.passwordstore.org/) didn't work.
Portability didn't seem to be a concern. So I remade the script to work, much more
simpily. Then I decided that it was kind of strange needing to have public and 
private keys just to encrypt some files with a password. That and my friend told 
me it was fucking weird to be using gpg for that (I somewhat disagree still). So
I made aes. 

Now you have this.
