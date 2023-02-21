package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func get_workspaces(binary_path string) []string {
	// Execute terraform command for getting all workspaces
	cmd := exec.Command(string(binary_path), "workspace", "list")
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	lines := strings.Split(string(stdout), "\n")

	// Remove the * character, strip whitespace from the lines, and skip blank lines.
	var workspaces []string
	for _, line := range lines {
		if strings.Contains(line, "*") {
			line = strings.ReplaceAll(line, "*", "")
		}
		line = strings.TrimSpace(line)
		if line != "" {
			workspaces = append(workspaces, line)
		}
	}

	// Return list of workspaces
	return workspaces
}

func deployment_data(binary_path string, workspaces []string) [][]string {
	// Store workspace, hostname and IP
	var machines [][]string

	for _, workspace := range workspaces {
		// Change workspace
		workspace_cmd := exec.Command(string(binary_path), "workspace", "select", workspace)
		_, err := workspace_cmd.Output()

		output_cmd := exec.Command(string(binary_path), "show", "-json")
		stdout, err := output_cmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		// Parse the JSON
		var data map[string]interface{}
		err = json.Unmarshal(stdout, &data)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if len(data) == 1 {
			machines = append(machines, []string{workspace, "", ""})
		} else {
			for _, child := range data["values"].(map[string]interface{})["root_module"].(map[string]interface{})["child_modules"].([]interface{})[1].(map[string]interface{})["child_modules"].([]interface{})[0].(map[string]interface{})["child_modules"].([]interface{}) {
				for _, resource := range child.(map[string]interface{})["resources"].([]interface{}) {
					if values, ok := resource.(map[string]interface{})["values"].(map[string]interface{}); ok {
						if clone, ok := values["clone"].([]interface{}); ok && len(clone) > 0 {
							if customize, ok := clone[0].(map[string]interface{})["customize"].([]interface{}); ok && len(customize) > 0 {
								if linuxOptions, ok := customize[0].(map[string]interface{})["linux_options"].([]interface{}); ok && len(linuxOptions) > 0 {
									if hostName, ok := linuxOptions[0].(map[string]interface{})["host_name"].(string); ok {
										if networkInterface, ok := customize[0].(map[string]interface{})["network_interface"].([]interface{}); ok && len(networkInterface) > 0 {
											if ipv4Address, ok := networkInterface[0].(map[string]interface{})["ipv4_address"].(string); ok {
												machines = append(machines, []string{workspace, hostName, ipv4Address})
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return machines
}

func main() {

	// Must set ENV for path to terraform binary
	terraform_binary_path := os.Getenv("TFB_PATH")
	if terraform_binary_path == "" {
		fmt.Println("Must set TFB_PATH Environment variable for path to terraform binary.")
	}

	// Get all workspaces
	workspaces := get_workspaces(terraform_binary_path)
	if workspaces == nil {
		fmt.Println("Could not get workspaces from terraform.")
		return
	}

	// Populate machines data
	machines := deployment_data(terraform_binary_path, workspaces)
	if machines == nil {
		fmt.Println("Something went wrong getting workspace data.")
		return
	}

	// Build table for data
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Workspace", "HostName", "IP Address"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	for _, v := range machines {
		table.Append(v)
	}

	table.Render()

}
