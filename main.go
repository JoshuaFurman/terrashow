package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// Define structs to match the JSON structure
type TerraformState struct {
	Resources []Resource `json:"resources"`
}

type Resource struct {
	Instances []Instance `json:"instances"`
}

type Instance struct {
	Attributes VMAttributes `json:"attributes"`
}

type VMAttributes struct {
	Clone []Clone `json:"clone"`
}

type Clone struct {
	Customize []Customize `json:"customize"`
}

type Customize struct {
	LinuxOptions     []LinuxOptions     `json:"linux_options"`
	NetworkInterface []NetworkInterface `json:"network_interface"`
}

type LinuxOptions struct {
	HostName string `json:"host_name"`
}

type NetworkInterface struct {
	IPv4Address string `json:"ipv4_address"`
	IPv6Address string `json:"ipv6_address"`
}

// VMData holds the combined information of a VM
type VMData struct {
	HostName    string
	IPv4Address string
	IPv6Address string
}

// TODO:
// - Figure out the ordering in the table...

func main() {
	full_data := make(map[string][]VMData)
	var wrkspace_buffer string
	var state_location string

	workspaces, err := getWorkspaces()
	if err != nil {
		fmt.Println("No workspace for:", err)
	}

	for _, wrkspace := range workspaces {
		if wrkspace == "default" || wrkspace == "* default" {
			vmData, err := parseTerraformState("terraform.tfstate")
			if err != nil {
				fmt.Println("Error parsing file:", err)
				continue
			}

			full_data[wrkspace] = vmData
		} else {
			if strings.Contains(wrkspace, "* ") {
				wrkspace_buffer = strings.TrimLeft(wrkspace, "* ")
				state_location = "terraform.tfstate.d/" + wrkspace_buffer + "/terraform.tfstate"
			} else {
				state_location = "terraform.tfstate.d/" + wrkspace + "/terraform.tfstate"
			}
			vmData, err := parseTerraformState(state_location)
			if err != nil {
				fmt.Println("Error parsing file:", err)
				continue
			}

			full_data[wrkspace] = vmData
		}
	}

	// Build table for data
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Workspace", "HostName", "IPv4 Address", "IPv6 Address"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)

	// for k, v := range full_data {
	// 	if len(v) == 0 {
	// 		row := []string{k, "", "", ""}
	// 		table.Append(row)
	// 	}
	// 	for _, info := range v {
	// 		row := []string{k, info.HostName, info.IPv4Address, info.IPv6Address}
	// 		table.Append(row)
	// 	}
	// }

	for _, k := range workspaces {
		if len(full_data[k]) == 0 {
			row := []string{k, "", "", ""}
			table.Append(row)
		}
		for _, info := range full_data[k] {
			row := []string{k, info.HostName, info.IPv4Address, info.IPv6Address}
			table.Append(row)
		}
	}
	table.Render()
}

func parseTerraformState(filename string) ([]VMData, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var state TerraformState
	err = json.Unmarshal(file, &state)
	if err != nil {
		return nil, err
	}

	var vms []VMData
	for _, resource := range state.Resources {
		for _, instance := range resource.Instances {
			for _, clone := range instance.Attributes.Clone {
				for _, customize := range clone.Customize {
					if len(customize.LinuxOptions) > 0 && len(customize.NetworkInterface) > 0 {
						vm := VMData{
							HostName:    customize.LinuxOptions[0].HostName,
							IPv4Address: customize.NetworkInterface[0].IPv4Address,
							IPv6Address: customize.NetworkInterface[0].IPv6Address,
						}
						vms = append(vms, vm)
					}
				}
			}
		}
	}

	return vms, nil
}

func listSubdirectories(directory string) ([]string, error) {
	dirEntries, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var subdirectories []string
	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			subdirectories = append(subdirectories, dirEntry.Name())
		}
	}

	return subdirectories, nil
}

func getWorkspaces() ([]string, error) {
	var workspaces []string
	// Get current workspace
	file, err := ioutil.ReadFile(".terraform/environment")
	if err != nil {
		return nil, err
	}
	current_workspace := string(file)

	// Get list of all workspaces
	wrkspace_buffer, err := listSubdirectories("terraform.tfstate.d")
	if err != nil {
		return nil, err
	}
	wrkspace_buffer = append([]string{"default"}, wrkspace_buffer...)

	for _, wrkspace := range wrkspace_buffer {
		if wrkspace == current_workspace {
			current_workspace = "* " + current_workspace
			workspaces = append(workspaces, current_workspace)
		} else {
			workspaces = append(workspaces, wrkspace)
		}
	}

	return workspaces, nil
}
