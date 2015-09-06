# pass and aes

aes is a simple go program that uses the crypto/aes library to encrypt and decrypt 
files for you. It makes it as easy as `aes -e file-to-encrypt`.

pass is a shell script that uses aes to encrypt files in a directory (`~/.secstore`)
and make it easy for you to get them. It is intended as a very simple password 
manager.

## origins

This all came about due to me switching to OpenBSD and discovering that 
[the standard unix password manager](http://www.passwordstore.org/) didn't work.
Portability didn't seem to be a concern. So I remade the script to work, much more
simpily. Then I decided that it was kind of strange needing to have public and 
private keys just to encrypt some files with a password. That and my friend told 
me it was fucking weird to be using gpg for that (I somewhat disagree still). So
I made aes. 

Now you have this.
