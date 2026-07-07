# mdwi

A static Markdown Wiki generator.

## Usage

To use `mdwi` just run it in the directory where your Markdown files are located.

Additional command line flags:

    Usage: mdwi [options]
    Options:
      -v, --version            Print version information and exit
      -h, --help               Print this message and exit
      -s, --standalone <file>  Create a standalone HTML file

## The Problem

1. Want to take notes
1. Want to do it in markdown
1. Want to be able to use a wiki style links
1. Want to export to HTML

There are many great personal note taking apps out there such as [Obsidian](https://obsidian.md/) that do 1-3 extremely well, but lock 4 behind a paywall.

This is an simple tool that does 1-4, albeit simply.

## How it Works

- Write notes in markdown files
- Make `index.md` your main entry point file
- Keep everything in a single folder (including image files)

```
    notes/
      |
      +---- index.md
      |
      +---- foo.md
      |
      +---- foo.png
```

- To link from `index.md` to `foo.md` use the following wiki-like syntax `{{foo}}`
- Run `mdwi` in your notes folder to generate HTML
- HTML files will be generated in a `_site` subdirectory, all image files will be copied there too

### Wiki Style Links

The files are linked together using wiki style links. If you have a file called `foo.md` and you want to link to it from another file, you can use the following syntax:

    {{foo}}

This will be seamlessly converted to:

    <a href="foo.html">foo</a>

### Using Images

If you want to use images, just dump them in the same directory as your markdown files. Link them using standard markdown syntax:

    ![alt text](image.png)

Once you run `mdwi` all the images will be copied to the `_site` directory.

Mdwi is oppinionated. It will generate a basic `style.css` file for you for styling. You can change it afterwards.

### Standalone Mode

In standalone mode, `mdwi` takes in a file name as an argument, and generates a single `index.html` file in the `_site` subdirectory as an output.

The output file contains:

- Inlined css stylesheet
- Inlined SVG favicon

As of 0.4.3, if the file contained standard markdown image tags, these images will be converted to base64 and inlined as well.

As such, this file is completely self contained.

Note: wiki style links will be converted to HTML links, but the linked files will not be converted or inlined.

## Installation


There are few different ways:

### Platform Independent

 Install via `go`:
 
    go install github.com/maciakl/mdwi@latest

### Linux

Use [grab](https://github.com/maciakl/grab):

    grab maciakl/mdwi

### macOS

Use `homebrew`:

    brew tap maciakl/tap
    brew trust maciakl/tap
    brew install mdwi

### Windows

Use `scoop` (see [scoop.sh](https://scoop.sh)).

    scoop bucket add maciak https://github.com/maciakl/bucket
    scoop update
    scoop install mdwi


