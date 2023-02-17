package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/olekukonko/tablewriter"
)

func main() {
	// // Read the file
	// json_data, err := ioutil.ReadFile("show_tenant")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// Must set ENV for path to terraform binary
	terraform_binary_path := os.Getenv("TFB_PATH")
	if terraform_binary_path == "" {
		fmt.Println("Must set TFB_PATH Environment variable for path to terraform binary.")
	}

	cmd := exec.Command(string(terraform_binary_path), "show", "-json")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Parse the JSON
	var data map[string]interface{}
	err = json.Unmarshal(stdout, &data)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Store hostname and IP
	var machines [][]string

	for _, child := range data["values"].(map[string]interface{})["root_module"].(map[string]interface{})["child_modules"].([]interface{})[1].(map[string]interface{})["child_modules"].([]interface{})[0].(map[string]interface{})["child_modules"].([]interface{}) {
		for _, resource := range child.(map[string]interface{})["resources"].([]interface{}) {
			if values, ok := resource.(map[string]interface{})["values"].(map[string]interface{}); ok {
				if clone, ok := values["clone"].([]interface{}); ok && len(clone) > 0 {
					if customize, ok := clone[0].(map[string]interface{})["customize"].([]interface{}); ok && len(customize) > 0 {
						if linuxOptions, ok := customize[0].(map[string]interface{})["linux_options"].([]interface{}); ok && len(linuxOptions) > 0 {
							if hostName, ok := linuxOptions[0].(map[string]interface{})["host_name"].(string); ok {
								if networkInterface, ok := customize[0].(map[string]interface{})["network_interface"].([]interface{}); ok && len(networkInterface) > 0 {
									if ipv4Address, ok := networkInterface[0].(map[string]interface{})["ipv4_address"].(string); ok {
										machines = append(machines, []string{hostName, ipv4Address})
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// Print the hostname and IPs
	if len(machines) == 0 {
		fmt.Println("No data found")
		return
	}

	// Build table for data
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"HostName", "IP Address"})
	table.SetRowLine(true)

	for _, v := range machines {
		table.Append(v)
	}

	table.Render()

	// workingDir := "/home/jfurman/istio/caas/terraform/config/vsphere"
	// execPath := "/home/jfurman/bin/terraform1.1.7"
	// tf, err := tfexec.NewTerraform(workingDir, execPath)
	// if err != nil {
	// 	fmt.Printf("error running NewTerraform: %s\n", err)
	// }

	// state, err := tf.Show(context.Background())
	// if err != nil {
	// 	fmt.Printf("error running Show: %s", err)
	// }

	// fmt.Println(state.Values.RootModule.Resources)
	// fmt.Printf("%T\n", state.Values)

}
