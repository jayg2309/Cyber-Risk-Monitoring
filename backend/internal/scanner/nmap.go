package scanner

import (
	"encoding/xml"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ScanResult represents a single port scan result
type ScanResult struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	State    string `json:"state"`
	Service  string `json:"service"`
	Version  string `json:"version"`
	Banner   string `json:"banner"`
}

// Scanner handles nmap scanning operations
type Scanner struct {
	timeout time.Duration
}

// NewScanner creates a new Scanner instance
func NewScanner(timeout time.Duration) *Scanner {
	return &Scanner{
		timeout: timeout,
	}
}

// ScanTarget performs an nmap scan on the specified target
func (s *Scanner) ScanTarget(target string) ([]ScanResult, error) {
	// Validate target format (basic validation)
	if target == "" {
		return nil, fmt.Errorf("target cannot be empty")
	}

	// Build nmap command: -sS (SYN scan), -sV (version detection), -p 1-1000 (port range), -oX - (XML output to stdout)
	args := []string{
		"-sS",          // SYN scan
		"-sV",          // Version detection
		"-p", "1-1000", // Port range
		"-oX", "-", // XML output to stdout
		"--host-timeout", fmt.Sprintf("%ds", int(s.timeout.Seconds())),
		target,
	}

	// Execute nmap command
	cmd := exec.Command("nmap", args...)

	// Set timeout for the command
	done := make(chan error, 1)
	var output []byte
	var err error

	go func() {
		output, err = cmd.Output()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return nil, fmt.Errorf("nmap scan failed: %v", err)
		}
	case <-time.After(s.timeout):
		cmd.Process.Kill()
		return nil, fmt.Errorf("nmap scan timed out after %v", s.timeout)
	}

	// Parse XML output
	results, err := s.parseNmapXML(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse nmap output: %v", err)
	}

	return results, nil
}

// NmapRun represents the root XML structure from nmap output
type NmapRun struct {
	XMLName xml.Name   `xml:"nmaprun"`
	Hosts   []NmapHost `xml:"host"`
}

// NmapHost represents a host in the nmap XML output
type NmapHost struct {
	XMLName xml.Name   `xml:"host"`
	Ports   NmapPorts  `xml:"ports"`
	Status  NmapStatus `xml:"status"`
}

// NmapStatus represents host status
type NmapStatus struct {
	XMLName xml.Name `xml:"status"`
	State   string   `xml:"state,attr"`
}

// NmapPorts represents the ports section
type NmapPorts struct {
	XMLName xml.Name   `xml:"ports"`
	Ports   []NmapPort `xml:"port"`
}

// NmapPort represents a single port in the XML output
type NmapPort struct {
	XMLName  xml.Name    `xml:"port"`
	Protocol string      `xml:"protocol,attr"`
	PortID   int         `xml:"portid,attr"`
	State    NmapState   `xml:"state"`
	Service  NmapService `xml:"service"`
}

// NmapState represents port state
type NmapState struct {
	XMLName xml.Name `xml:"state"`
	State   string   `xml:"state,attr"`
}

// NmapService represents service information
type NmapService struct {
	XMLName xml.Name `xml:"service"`
	Name    string   `xml:"name,attr"`
	Product string   `xml:"product,attr"`
	Version string   `xml:"version,attr"`
	Banner  string   `xml:"banner,attr"`
}

// parseNmapXML parses the XML output from nmap and returns structured results
func (s *Scanner) parseNmapXML(xmlData []byte) ([]ScanResult, error) {
	var nmapRun NmapRun

	// Parse XML
	err := xml.Unmarshal(xmlData, &nmapRun)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %v", err)
	}

	var results []ScanResult

	// Process each host
	for _, host := range nmapRun.Hosts {
		// Skip hosts that are not up
		if host.Status.State != "up" {
			continue
		}

		// Process each port
		for _, port := range host.Ports.Ports {
			// Only include open ports
			if port.State.State == "open" {
				version := port.Service.Version
				if port.Service.Product != "" {
					if version != "" {
						version = fmt.Sprintf("%s %s", port.Service.Product, version)
					} else {
						version = port.Service.Product
					}
				}

				result := ScanResult{
					Port:     port.PortID,
					Protocol: port.Protocol,
					State:    port.State.State,
					Service:  port.Service.Name,
					Version:  version,
					Banner:   port.Service.Banner,
				}

				results = append(results, result)
			}
		}
	}

	return results, nil
}

// ValidateTarget performs basic validation on the target string
func ValidateTarget(target string) error {
	if target == "" {
		return fmt.Errorf("target cannot be empty")
	}

	// Basic validation - check for obvious invalid characters
	if strings.ContainsAny(target, ";|&`$(){}[]<>") {
		return fmt.Errorf("target contains invalid characters")
	}

	return nil
}
