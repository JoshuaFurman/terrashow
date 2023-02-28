# Terrashow

## Overview

This tool is able to display the posture of your terraform deployments when invoked from the directory holding your terraform state. This will display workspaces, hostnames and IPs of machines managed by terraform.

## Before You Use

To use this tool you need to set the following environment variable:

- `TFB_PATH`: The full path to your terraform binary.

To build: `go build -o terrashow`

## Example

```bash
[jfurman@spoc.local@jump-linux-01 vsphere]$ terrashow
+-----------+------------------------+---------------+
| WORKSPACE |        HOSTNAME        |  IP ADDRESS   |
+-----------+------------------------+---------------+
| default   | istio-mgmt-lb-1        | 172.23.41.110 |
+           +------------------------+---------------+
|           | istio-mgmt-master-1    | 172.23.41.111 |
+           +------------------------+---------------+
|           | istio-mgmt-worker-1    | 172.23.41.112 |
+           +------------------------+---------------+
|           | istio-mgmt-worker-2    | 172.23.41.109 |
+           +------------------------+---------------+
|           | istio-mgmt-worker-3    | 172.23.41.113 |
+-----------+------------------------+---------------+
| primary   | istio-primary-master-1 | 44.130.4.6    |
+           +------------------------+---------------+
|           | istio-primary-worker-1 | 44.130.4.7    |
+           +------------------------+---------------+
|           | istio-primary-worker-2 | 44.130.4.8    |
+           +------------------------+---------------+
|           | istio-primary-worker-3 | 44.130.4.13   |
+-----------+------------------------+---------------+
| remote    |                        |               |
+-----------+------------------------+---------------+
```
