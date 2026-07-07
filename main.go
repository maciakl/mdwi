package main

import (
	"fmt"
	"net/url"
	"strings"
	"os"
	"path/filepath"
	"regexp"
	"encoding/base64"

	cp "github.com/otiai10/copy"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

)

const version = "0.4.3"

func main() {


	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-v", "--version":
			Version()
		case "-h", "--help":
			Usage()
		case "-s", "--standalone":
			if len(os.Args) < 3 {
				fmt.Fprintln(os.Stderr, "Error: no input file specified for standalone mode.")
				os.Exit(1)
			}
			inputFile := os.Args[2]
			if _, err := os.Stat(inputFile); os.IsNotExist(err) {
				fmt.Fprintln(os.Stderr, "Error: input file does not exist:", inputFile)
				os.Exit(1)
			}
			generateStandaloneFile(inputFile)
		default:
			Usage()
		}
	} else {
		generateWiki()
	}
}

func Version() {
	fmt.Println(filepath.Base(os.Args[0]), "version", version)
	os.Exit(0)
}

func Usage() {
	fmt.Println("Usage:", filepath.Base(os.Args[0]), "[options]")
	fmt.Println("Options:")
	fmt.Println("  -v, --version    		Print version information and exit")
	fmt.Println("  -h, --help       		Print this message and exit")
	fmt.Println("  -s, --standalone <file>	Create a standalone HTML file")
	os.Exit(0)
}

// check if a directory exists, if not create it, if it does exist, delete it and re-create it
func makeDir(path string) {

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// create the directory
		err := os.Mkdir(path, 0755)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error (mkdir):", err)
			os.Exit(1)
		}
		fmt.Println("Created", path, "directory")
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Error (dir check):", err)
		os.Exit(1)
	} else {
		// delete the directory and re-create it
		err := os.RemoveAll(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error (dir remove):", err)
			os.Exit(1)
		}
		err = os.Mkdir(path, 0755)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error (mkdir):", err)
			os.Exit(1)
		} else {
			fmt.Println("Removed and re-created", path, "directory")
		}
	}
}

func removeDir(path string) {

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return
	}

	_ = os.RemoveAll(path)
	if err == nil {
		fmt.Println("Removed", path, "directory")
	}
}


func generateWiki() {

	fmt.Println("Generating wiki using mdwi version", version, "...")

	makeDir("_site")  // create _site directory
	makeDir("_tmp")   // create _tmp directory

	// write the default stylesheet to _site/style.css
	stylesheet := generateStylesheetString()
	writeFile("_site/style.css", stylesheet, "Created _site/style.css", "css write")

	// write the default favicon to _site/favicon.svg
	favicon := generateFavicon()
	writeFile("_site/favicon.svg", favicon, "Created _site/favicon.svg", "favicon write")

	// remove file list.md if it exists
	_ = os.Remove("_list.md")
	fmt.Println("Removed _list.md")

	// find all markdown files in the current directory
	files, err := filepath.Glob("*.md")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error (md find):", err)
		os.Exit(1)
	}

	var list_builder strings.Builder

	list_builder.WriteString("# List of Pages\n\n")

	// iterate over the files and convert each markdown file to HTML
	for _, file := range files {

		outputFile := file[:len(file)-3] + ".html"
		outputPath := filepath.Join("_site", outputFile)

		// convert the markdown file to HTML and write it to the _site directory
		markdownFile(file, outputPath, false)

		// add the file to the list
		if file != "index.md" {
			fmt.Fprintf(&list_builder, "- [%s](%s)\n", file[:len(file)-3], url.PathEscape(outputFile))
		}
	}

	// write the list to list.md
	listInputPath := filepath.Join("_tmp", "list.md")
	writeFile(listInputPath, list_builder.String(), "Created _tmp/list.md", "list write")

	// convert list.md to HTML
	listOutputFile := "list.html"
	listOutputPath := filepath.Join("_site", listOutputFile)
	markdownFile(listInputPath, listOutputPath, false)

	// copy all the image files to the _site directory
	copyFiles("*.png")
	copyFiles("*.jpg")

	removeDir("_tmp") // remove _tmp directory

	fmt.Println("Done!")
}

// generate standalone html file with an inline stylesheet
func generateStandaloneFile(input_file string) {

		makeDir("_site")  // create _site directory

		output_file := filepath.Join("_site", "index.html")

		fmt.Println("Generating standalone HTML file:", output_file)
		markdownFile(input_file, output_file, true)
}

func writeFile(path string, content string, success_msg string, error_msg string) {

	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error (", error_msg, "):", err)
		os.Exit(1)
	} else {
		fmt.Println(success_msg)
	}
}

func copyFiles(filetype string) {

	// copy all the image files to the _site directory
	imgFiles, err := filepath.Glob(filetype)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error (", filetype, "find):", err)
		os.Exit(1)
	}

	for _, file := range imgFiles {
		src := file
		dst := filepath.Join("_site", file)
		err := cp.Copy(src, dst)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error (", filetype, "copy):", err)
			os.Exit(1)
		} else {
			fmt.Println("Copied", file, "to", dst)
		}
	}
}

func markdownFile(inputPath string, outputPath string, inline bool) {

	// Read the markdown file
	input, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error (md read):", err)
		os.Exit(1)
	}

	// Create a new markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// Parse the markdown content
	doc := markdown.Parse(input, p)

	// Create an HTML renderer with options
	opts := html.RendererOptions{
		Flags: html.CommonFlags | html.TOC | html.CompletePage,
	}
	renderer := html.NewRenderer(opts)

	// Render the markdown to HTML
	output := markdown.Render(doc, renderer)

	contentStr := string(output)

	// inject stylesheet before </head>
	if inline {
		// inline stylesheet
		contentStr = injectStylesheetInline(contentStr)
	} else {
		// link to external stylesheet
		re := regexp.MustCompile(`(?i)</head>`)
		contentStr = re.ReplaceAllString(contentStr, `<link rel="stylesheet" href="style.css">`+`$0`)
	}


	// find all instances of {{Name}} and replace them with <a href="Name.html">Name</a>
	reg := regexp.MustCompile(`\{\{([a-zA-Z0-9_ ]+)\}\}`)

	contentStr = reg.ReplaceAllStringFunc(contentStr, func(match string) string {
		name := match[2 : len(match)-2]
		link := fmt.Sprintf("<a href=\"%s.html\">%s</a>", url.PathEscape(name), name)
		return link
	})

	// inject custom HTML into the page
	if inline {
		contentStr = injectFaviconInline(contentStr) // inline svg favicon
	} else {
		contentStr = injectFavicon(contentStr) // link to external svg favicon file
	}

	// inject navigation links if not inline
	if !inline {
		contentStr = injectNav(contentStr)
	}

	contentStr = injectFooter(contentStr)  // footer


	if inline {
		// inline images
		contentStr = inlineImages(contentStr)
	}

	//convert outputStr back to []byte
	output = []byte(contentStr)

	// Write the HTML output to the specified file
	err = os.WriteFile(outputPath, output, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error (html write):", err)
		os.Exit(1)
	} else {
		fmt.Println("Converted", inputPath, "to", outputPath)
	}

}

func inlineImages(content string) string {

	fmt.Println("Inlining images...")

	// inline images by converting them to base64 and replacing the src attribute
	regImg := regexp.MustCompile(`(?i)<img\s+[^>]*src="([^"]+)"[^>]*>`)

	contentStr := regImg.ReplaceAllStringFunc(content, func(match string) string {
		// extract the src attribute value
		submatches := regImg.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match // no src found, return original match
		}
		imgPath := submatches[1]

		fmt.Println("Found image:", imgPath)

		// read the image file
		imgData, err := os.ReadFile(imgPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error (image read):", err)
			return match // return original match if error occurs
		}

		// determine the image MIME type based on the file extension
		var mimeType string
		switch strings.ToLower(filepath.Ext(imgPath)) {
		case ".png":
			mimeType = "image/png"
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".gif":
			mimeType = "image/gif"
		case ".svg":
			mimeType = "image/svg+xml"
		default:
			mimeType = "application/octet-stream" // fallback MIME type
		}

		// convert the image data to base64
		base64Data := base64.StdEncoding.EncodeToString(imgData)

		// create the new src attribute value
		newSrc := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)

		// replace the src attribute in the original match with the new src
		newMatch := strings.Replace(match, imgPath, newSrc, 1)
		return newMatch
	})

	return contentStr
}

func injectNav(content string) string {
	// Define the SVG icon as a string
	homeIconSVG := `
    <div class="links">
        <ul>
           <li><a href="index.html">🏠 Home</a></li>
           <li><a href="list.html">📁 List</a></li>
       </ul>
    </div>

    <h4>Table of Contents</h4>`

	// Use a regex to find the <body> tag
	re := regexp.MustCompile(`(?i)<nav[^>]*>`)
	// Replace it with the <body> tag and the SVG icon
	return re.ReplaceAllString(content, `$0`+homeIconSVG)
}

func injectFavicon(content string) string {
	// Define the favicon link tag
	favicon := `<link rel="icon" href="favicon.svg" type="image/svg+xml">`

	// Use a regex to find the <head> tag
	re := regexp.MustCompile(`(?i)<head[^>]*>`)
	// Replace it with the <head> tag and the favicon link
	return re.ReplaceAllString(content, `$0`+favicon)
}

func injectFaviconInline(content string) string {
	// Define the favicon link tag with inline SVG
	favicon := `<link rel="icon" href="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0OCIgaGVpZ2h0PSI0OCIgdmlld0JveD0iMCAwIDIwIDE2Ij48dGV4dCB4PSIwIiB5PSIxNCI+8J+TmjwvdGV4dD48L3N2Zz4=">`

	// Use a regex to find the <head> tag
	re := regexp.MustCompile(`(?i)<head[^>]*>`)
	// Replace it with the <head> tag and the favicon link
	return re.ReplaceAllString(content, `$0`+favicon)
}

func injectStylesheetInline(content string) string {

	// Define the stylesheet link tag with inline CSS
	stylesheet := `<style>` + generateStylesheetString() + `</style>`

	// inject stylesheet before </head>
	re := regexp.MustCompile(`(?i)</head>`)
	return re.ReplaceAllString(content, stylesheet+`$0`)
}

func injectFooter(content string) string {
	// Define the footer content
	footer := `
    <footer>
    <p>generated by <a href="https://github.com/maciakl/mdwi">mdwi</a> <small>%s</small></p>
    </footer>`

	// inject verson
	footer = fmt.Sprintf(footer, version)

	// Use a regex to find the closing </body> tag
	re := regexp.MustCompile(`(?i)</body>`)
	// Replace it with the closing </body> tag and the footer
	return re.ReplaceAllString(content, footer+`$0`)
}

func generateFavicon() string {
	// create a default favicon
	favicon := `<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 20 16"><text x="0" y="14">📚</text></svg>`
	return favicon
}

func generateStylesheetString() string {
	// create a default stylesheet
	stylesheet := `

body {
    font-family: "Avenir Next", Helvetica, Arial, sans-serif;
    padding:1em;
    margin:auto;
    max-width:42em;
    background:#fefefe;
}

h1, h2, h3, h4, h5, h6 {
	font-weight: bold;
}

h1 {
	color: #000000;
	font-size: 28pt;
    border-bottom: 1px solid gray;
}

h2 {
	border-bottom: 1px solid #CCCCCC;
	color: #000000;
	font-size: 24px;
}

h3 {
	font-size: 18px;
	border-bottom: 1px solid #CCCCCC;
}

h4, h5, h6 {
	text-decoration: underline;
}

h4 {
	font-size: 16px;
}

h5 {
	font-size: 14px;
}

h6 {
	color: #777777;
	background-color: inherit;
	font-size: 14px;
}

hr {
	height: 0.2em;
	border: 0;
	color: #CCCCCC;
	background-color: #CCCCCC;
}

p, blockquote, ul, ol, dl, li, table, pre {
	margin: 15px 0;
}

a, a:visited {
	color: #4183C4;
	background-color: inherit;
	text-decoration: none;
}

#message {
	border-radius: 6px;
	border: 1px solid #ccc;
	display:block;
	width:100%;
	height:60px;
	margin:6px 0px;
}

button, #ws {
	font-size: 10pt;
	padding: 4px 6px;
	border-radius: 5px;
	border: 1px solid #bbb;
	background-color: #eee;
}

code, pre, #ws, #message {
	font-family: Monaco;
	font-size: 10pt;
	border-radius: 3px;
	background-color: #F8F8F8;
	color: inherit;
}

code {
	border: 1px solid #EAEAEA;
	margin: 0 2px;
	padding: 0 5px;
}

pre {
	border: 1px solid #CCCCCC;
	overflow: auto;
	padding: 4px 8px;
}

pre > code {
	border: 0;
	margin: 0;
	padding: 0;
}

img {
    padding: 20px;
    max-width: 80%;
    height: auto;
    width: auto\9;
}

td {
    border: 1px solid lightGray;
    padding-left: 10px;
    padding-right: 10px;
    min-width: 150px;
}

th {
    border-bottom: 1px solid black;
    padding-left: 10px;
}

del {
    color: gray;
}

em {
    color: #088A85;
}

figure {
    border: 1px solid #CCCCCC;
    padding: 10px;
    background-color: #F8F8F8;
    margin: 10px;
}

figcaption {
    font-style: italic;
    font-size: 12px;
    color: darkGray;
}

footer {
    font-size: 10px;
    margin-top: 10em;
    border-top: 1px solid gray;
    text-align: right;
}

#ws { background-color: #f8f8f8; }

.send { color:#77bb77; }
.server { color:#7799bb; }
.error { color:#AA0000; }

nav {
     margin-top: 2em;
     position: absolute;
     left: 50px;
     width: 200px;
     font-size: 16px;
}

nav li {
    padding: 0;
    margin: 0
}

nav ul {
    margin-top: 0;
    margin-bottom: 0;
    padding-left: 15px;

}

@media print {
    nav {
        display: none !important;
    }

    h2 {
        page-break-before: auto;
        page-break-after: avoid;
    }

    h2, h3, h4 {
        page-break-after: avoid;
    }

    img {
        display: block;
        margin-left: auto;
        margin-right: auto;
        width: 4.5in;
        page-break-before: auto;
        page-break-after: auto;
        page-break-inside: avoid;
    }

    table {
        page-break-before: auto;
        page-break-after: auto;
        page-break-inside: avoid;
    }

   a:link, a:visited {
        text-decoration: underline
   }

   a:link:after, a:visited:after {
       content: " (" attr(href) ") ";
       font-size: 90%;
   }

}

@media (max-width: 1100px) { 
    nav {
        margin-top: 2em;
        left: 0;
        position: relative;
    }
}`

	return stylesheet
}
