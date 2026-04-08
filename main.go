package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Finding struct {
	Name     string `json:"name"`
	Severity string `json:"severity"`
	Host     string `json:"host"`
	Matched  string `json:"matched,omitempty"`
	Details  struct {
		IP               string `json:"ip,omitempty"`
		Cloud            string `json:"cloud,omitempty"`
		HijackableDomain string `json:"hijackable_domain,omitempty"`
		Type             string `json:"type,omitempty"`
	} `json:"details,omitempty"`
}

type ScanCompletion struct {
	Status       string  `json:"status"`
	ScanDuration float64 `json:"scan_duration"`
}

func main() {
	filePtr := flag.String("l", "", "File containing subdomains")
	keyPtr := flag.String("SUBPIPE_API_KEY", "", "SubPipe API Key")
	flag.Parse()

	apiKey := *keyPtr
	if apiKey == "" {
		apiKey = os.Getenv("SUBPIPE_API_KEY")
	}

	if apiKey == "" {
		color.Red("Error: API Key is required.")
		color.Yellow("Use --SUBPIPE_API_KEY flag or set the SUBPIPE_API_KEY environment variable.")
		os.Exit(1)
	}

	// Environment Detection: Default to Prod, switch to Local if DEV is set
	apiURL := "https://api.subpipe.run/v1/scan"
	if os.Getenv("SUBPIPE_DEV") == "true" || os.Getenv("SUBPIPE_ENV") == "local" {
		apiURL = "http://localhost:3000/v1/scan"
		color.Magenta("🛠️  Development Mode: Connecting to %s", apiURL)
	}

	var subdomains []string
	if *filePtr != "" {
		content, err := os.ReadFile(*filePtr)
		if err != nil {
			color.Red("Error reading file: %v", err)
			os.Exit(1)
		}
		subdomains = strings.Split(string(content), "\n")
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				subdomains = append(subdomains, scanner.Text())
			}
		}
	}

	if len(subdomains) == 0 {
		color.Yellow("Usage: subfinder -d target.com | ./subpipe")
		os.Exit(1)
	}

	// Layer 3 Defense: Local Truncation
	if len(subdomains) > 10000 {
		color.Yellow("⚠️  Warning: Payload truncated to first 10,000 targets.")
		subdomains = subdomains[:10000]
	}

	startScan(apiURL, subdomains, apiKey)
}

func startScan(url string, targets []string, apiKey string) {
	var cleanTargets []string
	for _, t := range targets {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			cleanTargets = append(cleanTargets, trimmed)
		}
	}

	// IMPORTANT: Key must be "targets" to match Backend Validator
	payload, _ := json.Marshal(map[string][]string{"targets": cleanTargets})

	client := &http.Client{Timeout: 0}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		color.Red("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		color.Red("Connection Error: Could not reach SubPipe Engine at %s", url)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		
		// Detailed Error Parsing
		var errResp struct {
			Error string `json:"error"`
		}
		
		fmt.Println(strings.Repeat("-", 70))
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
			color.Red("❌ Server Error %d: %s", resp.StatusCode, errResp.Error)
		} else {
			// Fallback: Show raw body if it's not JSON
			rawMsg := strings.TrimSpace(string(body))
			if rawMsg == "" {
				rawMsg = http.StatusText(resp.StatusCode)
			}
			color.Red("❌ Server Error %d: %s", resp.StatusCode, rawMsg)
		}
		fmt.Println(strings.Repeat("-", 70))
		return
	}

	color.Cyan("🚀 SubPipe Analysis Started: %d targets sent to %s", len(cleanTargets), url)
	fmt.Println(strings.Repeat("-", 70))

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF { break }
			color.Red("Stream Interrupted: %v", err)
			break
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			if strings.Contains(data, "\"status\":\"completed\"") {
				var complete ScanCompletion
				json.Unmarshal([]byte(data), &complete)
				fmt.Println(strings.Repeat("-", 70))
				color.HiGreen("✅ Scan Finished in %.2fs", complete.ScanDuration)
				continue
			}

			var f Finding
			if err := json.Unmarshal([]byte(data), &f); err == nil {
				printFinding(f)
			}
		}
	}
}

func printFinding(f Finding) {
	sevColor := color.New(color.FgWhite).SprintFunc()
	switch strings.ToLower(f.Severity) {
	case "critical": sevColor = color.New(color.FgHiRed, color.Bold).SprintFunc()
	case "high": sevColor = color.New(color.FgRed).SprintFunc()
	case "medium": sevColor = color.New(color.FgHiYellow).SprintFunc()
	case "low": sevColor = color.New(color.FgCyan).SprintFunc()
	case "info": sevColor = color.New(color.FgBlue).SprintFunc()
	}

	timestamp := color.New(color.FgWhite).Sprintf("[%s]", time.Now().Format("15:04:05"))
	
	detailStr := ""
	if f.Matched != "" {
		detailStr = fmt.Sprintf(" (%s)", color.HiBlueString(f.Matched))
	} else if f.Details.IP != "" {
		if f.Details.Cloud != "" {
			detailStr = fmt.Sprintf(" [%s - %s]", color.CyanString(f.Details.Cloud), color.WhiteString(f.Details.IP))
		} else {
			detailStr = fmt.Sprintf(" [%s]", color.WhiteString(f.Details.IP))
		}
	} else if f.Details.HijackableDomain != "" {
		detailStr = fmt.Sprintf(" [Takeover: %s]", color.MagentaString(f.Details.HijackableDomain))
	}

	fmt.Printf("%s %s %s: %s%s\n",
		timestamp,
		sevColor(fmt.Sprintf("%-8s", strings.ToUpper(f.Severity))),
		color.WhiteString(f.Name),
		color.GreenString(f.Host),
		detailStr,
	)
}