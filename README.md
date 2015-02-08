[![Build Status](https://travis-ci.org/julienc91/heygo.png)](https://travis-ci.org/julienc91/heygo)

HEYGO
=====

HEYGO is a web application that offers an easy way to share media
within a small circle of users. HEYGO meets all your expectations
about privacy:

* *you* control the media you share
* *you* decide who can be part of your circle
* *you* establish the access rules of your media


## This is an Alpha version

**HEYGO is currently in Alpha version**. Many things are not implemented yet, there are still a lot of bugs, probably some **security flaws** too, and the design of the web application is not even closed to be finished. Be aware of this before deploying it on your server !

What you can do however is helping me to continue this project by [reporting bugs or sharing ideas](https://github.com/julienc91/heygo/issues), and I will try to do my best to make HEYGO the web application we want it to be !

## Getting started

### Dependencies

HEYGO requires an installation of [Go](http://golang.org/doc/install) (tested with the 1.3 version, but it should work with Go 1.2 too) and [Bower](http://bower.io/).

It is also recommended to use [Grunt](http://gruntjs.com/) for the translation-making process, and a web proxy server such as [nginx](http://nginx.org/) for deployment.

### First time installation

Clone the repository on your server:

    git clone https://github.com/julienc91/heygo
    cd heygo

Then install the Javascript dependencies:

    bower install

Now you can run the server-side application:

    go build heygo
    ./heygo

### Usage

Currently, HEYGO uses a sqlite database to store its configuration. By default, the server application listens on `localhost:8080` and creates a user with admin rights `admin:admin`. You can change both the login and the password in the administration panel.



## Web proxy configuration

If you want to use port 80 for HEYGO, you should definitely not run HEYGO as root. Instead, use a web proxy such as [nginx](http://nginx.org/). Here is a basic nginx configuration file to access HEYGO on example.com while HEYGO is set to run on `localhost:8080`.

    $> cat /etc/nginx/sites-enabled/heygo.conf
    server {
        listen 80;
        server_name example.com;
        
        location / {
			    proxy_pass http://127.0.0.1:8080;
			    proxy_set_header Host $host;
        }
    }

## What will come next

This is a non-exhaustive list of what might be implemented in future versions of HEYGO:

* Videos are good, but pictures are nice to share too !
* Mass import of media from the administration panel
* Using Grunt for a lot of other cool stuff
* Sorting media by groups
* Support of `https`, because confidentiality is essential !
* HEYGO binaries for multiple platforms
* ...

## Licence

HEYGO is published under the GPL public license.

## Author

Julien CHAUMONT
https://julienc.io


