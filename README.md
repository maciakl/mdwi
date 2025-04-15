# mdwi

A static Markdown Wiki generator.

## Usage

To use `mdwi` just run it in the directory where your Markdown files are located.

Additional command line flags:

    Usage: mdwi [options]
    Options:
      -v, --version    Print version information and exit
      -h, --help       Print this message and exit

## The Problem

1. Want to take notes
1. Want to do it in markdown
1. Want to be able to use a wiki style links
1. Want to export to HTML

There are many great personal note taking apps out there such as [Obsidian](https://obsidian.md/) that do 1-3 extremely well, but lock 4 behind a paywall.

This is an simple tool that does 1-4, albeit simply.

## How it Works

Write notes in Markdown files.

The files are linked together using wiki style links. For example, if you have a file called `foo.md` and you want to link to it from another file, you can use the following syntax:

    {{foo}}

Everything else is just plain markdown. Keep a flat file structure (no folders). If you want to use images, just dump them in the same directory. Link them using standard markdown syntax:

    ![alt text](image.png)

When done run:

    mdwi

This will create a directory named `_site` and generate HTML version of your notes there.

Mdwi is oppinionated. It will generate a basic `style.css` file for you for styling. You can change it afterwards.


## Dependencies

Mdwi uses Pandoc to convert markdown to HTML because it exists, it's mature and works well.

- [Pandoc](https://pandoc.org/) - for converting markdown to HTML


## Installation


There are few different ways:

### Platform Independent

 Install via `go`:
 
    go install github.com/maciakl/mdwi@latest

### Linux

On Linux (requires `wget` & `unzip`, installs to `/usr/local/bin`):

    p="mdwi" && wget -qN "https://github.com/maciakl/${p}/releases/latest/download/${p}_lin.zip" && unzip -oq ${p}_lin.zip && rm -f ${p}_lin.zip && chmod +x ${p} && sudo mv ${p} /usr/local/bin

To uninstall, simply delete it:

    rm -f /usr/local/bin/mdwi

### Windows

On Windows, this tool is distributed via `scoop` (see [scoop.sh](https://scoop.sh)).

 First, you need to add my bucket:

    scoop bucket add maciak https://github.com/maciakl/bucket
    scoop update

 Next simply run:
 
    scoop install mdwi

If you don't want to use `scoop` you can simply download the executable from the release page and extract it somewhere in your path.


