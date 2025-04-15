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


