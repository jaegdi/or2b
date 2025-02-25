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

func loginCluster(cluster string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("ocl %s", cluster))
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error logging into cluster", cluster, ":", err, "\n", output)
		return err
	}
	return nil
}

func initClusters(clusters *[]string, login bool) {
	// Trim spaces from cluster names
	for i := range *clusters {
		(*clusters)[i] = strings.TrimSpace((*clusters)[i])
	}

	if login {
		currentcluster := os.Getenv("CLUSTER")
		for _, cluster := range *clusters {
			if err := loginCluster(cluster); err != nil {
				continue
			}
			fmt.Println("Logged into cluster", cluster)
		}
		// Switch back to the original cluster
		if err := loginCluster(currentcluster); err != nil {
			fmt.Println("Error switching back to the original cluster", currentcluster, ":", err)
		} else {
			fmt.Println("Switched back to cluster", currentcluster)
		}
	}
}

func main() {
	// Parse command-line arguments
	filterPattern := flag.String("pattern", "", "Filter pattern for route names or hosts")
	flag.StringVar(filterPattern, "p", "", "Filter pattern for route names or hosts (shorthand)")

	denyPattern := flag.String("denypattern", "", "Deny-Filter pattern for route names, namespaces or hosts")
	flag.StringVar(denyPattern, "d", "", "Deny-Filter pattern for route names, namespaces or hosts (shorthand)")

	outputFileName := flag.String("output", "", "Output file name")
	flag.StringVar(outputFileName, "o", "", "Output file name (shorthand)")

	login := flag.Bool("login", false, "Login into clusters")
	flag.BoolVar(login, "l", false, "Login into clusters (shorthand)")

	// Add new parameter for clusters
	clustersFlag := flag.String("clusters", "dev-scp0,cid-scp0,ppr-scp0,pro-scp0,pro-scp1", "Comma-separated list of clusters")
	flag.StringVar(clustersFlag, "c", "dev-scp0,cid-scp0,ppr-scp0,pro-scp0,pro-scp1", "Comma-separated list of clusters (shorthand)")

	flag.Parse()

	// Split clusters string into slice
	clusters := strings.Split(*clustersFlag, ",")
	initClusters(&clusters, *login)

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf(`<!DOCTYPE NETSCAPE-Bookmark-file-1>` + "\n"))
	buffer.WriteString(fmt.Sprintf(`<META HTTP-EQUIV="Content-Type" CONTENT="text/html; charset=UTF-8">` + "\n"))
	buffer.WriteString(fmt.Sprintf(`<TITLE>Bookmarks %s</TITLE>`+"\n", *filterPattern))
	buffer.WriteString(fmt.Sprintf(`<H1>Bookmarks %s</H1>`+"\n", *filterPattern))
	buffer.WriteString(fmt.Sprintf(`<DL><p>` + "\n"))

	for _, cluster := range clusters {
		// Execute the command to get the routes for each cluster
		cmd := exec.Command("bash", "-c", fmt.Sprintf("oc get routes --all-namespaces --no-headers --context=default/api-%s-sf-rz-de:6443/jaegdi", cluster))
		output, err := cmd.Output()
		if err != nil {
			fmt.Println("Error fetching routes for cluster", cluster, ":", err)
			continue
		}

		buffer.WriteString(fmt.Sprintf(`	<DT><H3 ADD_DATE="%d" LAST_MODIFIED="%d">%s</H3>`+"\n", time.Now().Unix(), time.Now().Unix(), cluster))

		buffer.WriteString(fmt.Sprintf(`	<DL><p>` + "\n"))

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
			if *filterPattern != "" &&
				(strings.Contains(name, *filterPattern) || strings.Contains(host, *filterPattern)) &&
				!(strings.Contains(name, *denyPattern) || strings.Contains(namespace, *denyPattern) || strings.Contains(host, *denyPattern)) {
				buffer.WriteString(fmt.Sprintf(`		<DT><A HREF="http://%s" ADD_DATE="%d">%s -- %s</A>`+"\n", host, time.Now().Unix(), namespace, name))
			}
		}
	}
	// End the HTML document
	buffer.WriteString(fmt.Sprintf(`</DL><p>` + "\n"))

	// Write to output file or STDOUT
	if *outputFileName == "" {
		// Write to STDOUT
		fmt.Print(buffer.String())
	} else {
		// Write to file
		err := os.WriteFile(*outputFileName, []byte(buffer.String()), 0644)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Printf("Bookmarks file created: %s\n", *outputFileName)
	}
}
