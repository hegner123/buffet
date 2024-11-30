package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
)

func main() {
    fmt.Println("check for defined editor")
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim" // default editor
	}

    fmt.Println("create temp file")
	// 2. Create a temporary file
	tempFile, err := os.CreateTemp("", "command-*")
	if err != nil {
		log.Fatal(err)
	}
    fmt.Println("defore closing the file")
	defer os.Remove(tempFile.Name()) // Clean up the temp file afterward

    fmt.Println("write init content to temp file")
	// 3. Write initial content to the temporary file (optional)
	initialContent := `
# Save and quit the editor to output your command to the prompt.
# Lines starting with '#' will be ignored. 
# A blank prompt aborts buffet.
`
    fmt.Println("write to tempfile")
	if _, err := tempFile.WriteString(initialContent); err != nil {
		log.Fatal(err)
	}
    fmt.Println("close tmp file")
	tempFile.Close() // Close the file so the editor can open it

    fmt.Println("define args")
	var args string
    fmt.Println("define cmd")
	var cmd *exec.Cmd


    fmt.Println("open editor")
	switch {
	case editor == "nvim":
		args = "+set filetype=sh"
		cmd = exec.Command(editor, args, tempFile.Name())
	default:
		cmd = exec.Command(editor, tempFile.Name())
	}
    
    fmt.Println("setting cmd inputs and outputs")

    fmt.Println("check terminal redirection input")
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Println("input is not terminal")
	} else {
		fmt.Println("input is terminal")
        cmd.Stdin = os.Stdin
	}
    fmt.Println("check terminal redirection output")
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Println("output is not terminal")
	} else {
		fmt.Println("output is terminal")
        cmd.Stdout = os.Stdout
	}
	cmd.Stderr= os.Stderr

    
    fmt.Println("run command")

	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run editor: %v", err)
	}

    fmt.Println("read contents")
	// 5. Read and process the contents
	contentBytes, err := os.ReadFile(tempFile.Name())
	if err != nil {
		log.Fatal(err)
	}
	content := string(contentBytes)

    fmt.Println("process content")
	// 6. Process the content
	var processedLines []string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" || strings.HasPrefix(trimmedLine, "#") {
			continue // Ignore empty lines and comments
		}
		processedLines = append(processedLines, line)
	}

    fmt.Println("join finalMessage")
	finalMessage := strings.Join(processedLines, "\n")
	if strings.TrimSpace(finalMessage) == "" {
		log.Fatal("Aborting due to empty message.")
	}

    fmt.Println("write to stdout")
	os.Stdout.Write([]byte(finalMessage))
}
