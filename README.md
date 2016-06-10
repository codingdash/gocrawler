# gocrawler
A Basic crawler written in Golang.


### Setting up localbox
This program is tested on Ubuntu 14.04, assuming you have your golang env  setup already if not you can follow [this article](https://www.digitalocean.com/community/tutorials/how-to-install-go-1-6-on-ubuntu-14-04) to setup your box locally.


### Downloading the source
To get the source, simply type the ``go get`` command-

```
go get github.com/codingdash/gocrawler
```

This will download the source in your GOPATH src directory. You can run the binary from ``$GOPATH/bin`` directory by firing below command-

```
$ $GOPATH/bin/gocrawler codingdash.com
```

This will crawl the website to your current directory.


### Improvements
Following things can be improved in the program-
 
- Passing the level to crawl the website usign command line, currently it's hard coded and set to 5.
- Make processing async for the document.
- Update the link to relative paths in the html document.