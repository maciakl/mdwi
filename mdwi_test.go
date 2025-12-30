package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testDir      = "test_run"
	mdwiBinary   = "mdwi"
	indexMd      = "index.md"
	anotherMd    = "another.md"
	imagePng     = "image.png"
	expectedSite = "_site"
)

var mdwiBinaryAbsPath string

// setup and teardown
func TestMain(m *testing.M) {
	// Compile the mdwi binary
	cmd := exec.Command("go", "build", "-o", mdwiBinary)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Failed to compile mdwi binary: %v\n", err)
		os.Exit(1)
	}

	// Get the absolute path of the binary
	mdwiBinaryAbsPath, err = filepath.Abs(mdwiBinary)
	if err != nil {
		fmt.Printf("Failed to get absolute path of mdwi binary: %v\n", err)
		os.Exit(1)
	}

	// Run the tests
	exitCode := m.Run()

	// Clean up the compiled binary
	os.Remove(mdwiBinary)

	os.Exit(exitCode)
}

func setupTestDir(t *testing.T) {
	// Create a test directory
	err := os.Mkdir(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Create dummy markdown and image files
	createDummyFile(t, filepath.Join(testDir, indexMd), "# Index\n\nThis is the index page with a link to {{another}}.")
	createDummyFile(t, filepath.Join(testDir, anotherMd), "# Another Page\n\nThis is another page.")
	createDummyFile(t, filepath.Join(testDir, imagePng), "dummy image content")
}

func teardownTestDir(t *testing.T) {
	// Remove the test directory
	err := os.RemoveAll(testDir)
	if err != nil {
		t.Fatalf("Failed to remove test directory: %v", err)
	}
}

func createDummyFile(t *testing.T, path, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create dummy file %s: %v", path, err)
	}
}

// Test Functions
func TestNoArgs(t *testing.T) {
	setupTestDir(t)
	defer teardownTestDir(t)

	cmd := exec.Command(mdwiBinaryAbsPath)
	cmd.Dir = testDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("mdwi command failed: %v\nOutput: %s", err, string(output))
	}

	// check if the _site directory was created
	siteDir := filepath.Join(testDir, expectedSite)
	if _, err := os.Stat(siteDir); os.IsNotExist(err) {
		t.Errorf("_site directory was not created")
	}

	// check for expected files in _site
	expectedFiles := []string{
		"index.html",
		"another.html",
		"list.html",
		"style.css",
		"favicon.svg",
		"image.png",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(siteDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created in _site", file)
		}
	}

	// check if the link was replaced in index.html
	indexContent, err := os.ReadFile(filepath.Join(siteDir, "index.html"))
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}
	if !strings.Contains(string(indexContent), `<a href="another.html">another</a>`) {
		t.Errorf("Link was not replaced in index.html")
	}
}

func TestVersionFlag(t *testing.T) {
	cmd := exec.Command(mdwiBinaryAbsPath, "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("mdwi -v command failed: %v\nOutput: %s", err, string(output))
	}

	if !strings.Contains(string(output), "mdwi version") {
		t.Errorf("Expected version information, but got: %s", string(output))
	}
}

func TestHelpFlag(t *testing.T) {
	cmd := exec.Command(mdwiBinaryAbsPath, "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("mdwi -h command failed: %v\nOutput: %s", err, string(output))
	}

	if !strings.Contains(string(output), "Usage:") {
		t.Errorf("Expected usage information, but got: %s", string(output))
	}
}

func TestStandaloneFlag(t *testing.T) {
	setupTestDir(t)
	defer teardownTestDir(t)

	outputFile := filepath.Join(testDir, "another.html")

	cmd := exec.Command(mdwiBinaryAbsPath, "-s", anotherMd)
	cmd.Dir = testDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("mdwi -s command failed: %v\nOutput: %s", err, string(output))
	}

	// check if the output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("standalone file %s was not created", outputFile)
	}

	// check if the stylesheet is inlined
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read standalone file: %v", err)
	}
	if !strings.Contains(string(content), "<style>") {
		t.Errorf("Stylesheet was not inlined in standalone file")
	}
}

func TestPandocNotInstalled(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "test-pandoc")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a dummy mdwi binary in the temp dir
	originalBinaryPath, err := filepath.Abs(mdwiBinary)
	if err != nil {
		t.Fatalf("Failed to get absolute path of mdwi binary: %v", err)
	}
	dummyBinaryPath := filepath.Join(tmpDir, mdwiBinary)
	copyFile(t, originalBinaryPath, dummyBinaryPath)
	os.Chmod(dummyBinaryPath, 0755) // make it executable

	// Run the command with a modified PATH that doesn't include pandoc
	cmd := exec.Command(mdwiBinaryAbsPath)
	cmd.Dir = tmpDir
	cmd.Env = []string{"PATH="} // Empty PATH
	output, err := cmd.CombinedOutput()

	if err == nil {
		t.Fatalf("mdwi command should have failed without pandoc, but it didn't")
	}
	if !strings.Contains(string(output), "Error: pandoc is not installed.") {
		t.Errorf("Expected error message about pandoc not being installed, but got: %s", string(output))
	}
}

func TestInternalLinks(t *testing.T) {
	// Create a temporary directory for the test
	tmpDir, err := os.MkdirTemp("", "test-internal-links")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create dummy markdown files with internal links
	page1Md := "page1.md"
	page2Md := "page2.md"
	createDummyFile(t, filepath.Join(tmpDir, page1Md), "# Page 1\n\nLink to {{page2}}.")
	createDummyFile(t, filepath.Join(tmpDir, page2Md), "# Page 2\n\nLink to {{page1}}.")

	// Run mdwi in the temporary directory
	cmd := exec.Command(mdwiBinaryAbsPath)
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("mdwi command failed: %v\nOutput: %s", err, string(output))
	}

	// Check if the HTML files were created
	siteDir := filepath.Join(tmpDir, expectedSite)
	page1Html := filepath.Join(siteDir, "page1.html")
	page2Html := filepath.Join(siteDir, "page2.html")

	if _, err := os.Stat(page1Html); os.IsNotExist(err) {
		t.Errorf("page1.html was not created")
	}
	if _, err := os.Stat(page2Html); os.IsNotExist(err) {
		t.Errorf("page2.html was not created")
	}

	// Check if the links were replaced correctly
	page1Content, err := os.ReadFile(page1Html)
	if err != nil {
		t.Fatalf("Failed to read page1.html: %v", err)
	}
	if !strings.Contains(string(page1Content), `<a href="page2.html">page2</a>`) {
		t.Errorf("Internal link in page1.html is incorrect")
	}

	page2Content, err := os.ReadFile(page2Html)
	if err != nil {
		t.Fatalf("Failed to read page2.html: %v", err)
	}
	if !strings.Contains(string(page2Content), `<a href="page1.html">page1</a>`) {
		t.Errorf("Internal link in page2.html is incorrect")
	}
}


func copyFile(t *testing.T, src, dst string) {
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", src, err)
	}
	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write file %s: %v", dst, err)
	}
}