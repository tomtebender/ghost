# Ghost
## What is Ghost?
Ghost is a light-weight HTTP/HTTPS server that runs directories as websites.

## How do I use Ghost?
You need a directory containing a file called "index.html", and that's it.

If you want to use TLS, create a directory called "secret" within the
website's directory, containing the files "key.pem" and "cert.pem", and
Ghost will automatically recognize the certificate.

## Usage
```
usage: ghost [--ip <ip or domain>] [--port <port>] [<directory>]

<directory> | The directory of the webpage (default: $PWD)
OPTIONS:
	--ip <ip or domain> | The ip or domain listened on (default: localhost)
	--port <port> | The port listened on (default: auto)
	--help | Show this text
```
