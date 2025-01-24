package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	// Parse command-line arguments
	filterPattern := flag.String("pattern", "", "Filter pattern for route names or hosts")
	flag.StringVar(filterPattern, "p", "", "Filter pattern (shorthand)")

	outputFileName := flag.String("output", "chrome_bookmarks.html", "Output file name")
	flag.StringVar(outputFileName, "o", "chrome_bookmarks.html", "Output file name (shorthand)")

	flag.Parse()

	// Check if a filter pattern was provided
	if *filterPattern == "" {
		fmt.Println("Please provide a filter pattern using the -pattern flag.")
		return
	}

	// Get the output filename from the flag
	outputFile := *outputFileName

	// Start creating the HTML document
	var buffer bytes.Buffer
	buffer.WriteString(`<!DOCTYPE NETSCAPE-Bookmark-file-1>
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
<DT><H3 ADD_DATE="1609459200" LAST_MODIFIED="1651017600" PERSONAL_TOOLBAR_FOLDER="true" PROTECTION="weak" TITLE="OpenShift Routes">
    <A HREF="http://example.com" ADD_DATE="1609459200">OpenShift Routes</A>
<DL><p>
`)

	// Execute the command to get the routes
	cmd := exec.Command("oc", "get", "routes", "--all-namespaces", "--no-headers")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error fetching routes:", err)
		return
	}

	// Iterate over each route
	for _, line := range strings.Split(string(output), "\n") {
		if line == "" {
			continue // Skip empty lines
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue // Skip lines that do not have sufficient data
		}

		namespace, name, host := parts[0], parts[1], parts[2]

		// Filter by pattern
		if strings.Contains(name, *filterPattern) || strings.Contains(host, *filterPattern) {
			buffer.WriteString(fmt.Sprintf(`<DT><A HREF="http://%s" ADD_DATE="%d">%s in %s</A>
`, host, time.Now().Unix(), name, namespace))
		}
	}

	// End the HTML document
	buffer.WriteString(`</DL><p>
</DL><p>
`)

	// Write to the output file
	err = os.WriteFile(outputFile, []byte(buffer.String()), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("Bookmarks file created: %s\n", outputFile)
}
