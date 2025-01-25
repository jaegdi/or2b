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

	outputFileName := flag.String("output", "", "Output file name")
	flag.StringVar(outputFileName, "o", "", "Output file name (shorthand)")

	flag.Parse()

<<<<<<< Updated upstream
	// Get the output filename from the flag
	outputFile := *outputFileName

	// Define the clusters
	clusters := []string{"dev-scp0", "cid-scp0", "ppr-scp0", "pro-scp0", "pro-scp1"}

	// Start creating the HTML document
=======
	cmd := exec.Command("bash", "-c", "oc whoami --show-server")

>>>>>>> Stashed changes
	var buffer bytes.Buffer
	buffer.WriteString(`<!DOCTYPE NETSCAPE-Bookmark-file-1>
<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">
<TITLE>Bookmarks</TITLE>
<H1>Bookmarks</H1>
<DL><p>
`)

	for _, cluster := range clusters {
		buffer.WriteString(fmt.Sprintf(`<DT><H3 ADD_DATE="%d" LAST_MODIFIED="%d" TITLE="%s">
<DL><p>
`, time.Now().Unix(), time.Now().Unix(), cluster))

		// Execute the command to get the routes for each cluster
		cmd := exec.Command("bash", "-c", fmt.Sprintf("oc get routes --all-namespaces --no-headers --context=default/api-%s-sf-rz-de:6443/jaegdi", cluster))
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error fetching routes for cluster", cluster, ":", err)
			continue
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
			if *filterPattern != "" && (strings.Contains(name, *filterPattern) || strings.Contains(host, *filterPattern)) {
				buffer.WriteString(fmt.Sprintf(`<DT><A HREF="http://%s" ADD_DATE="%d">%s in %s</A>
`, host, time.Now().Unix(), name, namespace))
			}
		}

		buffer.WriteString(`</DL><p>
`)
	}

	// End the HTML document
	buffer.WriteString(`</DL><p>
`)

<<<<<<< Updated upstream
	// Write to the output file
	err := os.WriteFile(outputFile, []byte(buffer.String()), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
=======
	// Write to output file or STDOUT
	if *outputFileName == "" {
		// Write to STDOUT
		fmt.Print(buffer.String())
	} else {
		// Write to file
		err = os.WriteFile(*outputFileName, []byte(buffer.String()), 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Printf("Bookmarks file created: %s\n", *outputFileName)
>>>>>>> Stashed changes
	}
}
