# Terrashow

## Overview

This tool is able to display the posture of your terraform deployments when invoked from the directory holding your terraform state. This will display workspaces, hostnames and IPs of machines managed by terraform.

## Before You Use

To build: `go build -o terrashow`

## Example

```bash
[jfurman@spoc.local@jump-linux-01 vsphere]$ terrashow
+-----------+------------------------------+--------------+--------------+
| WORKSPACE |           HOSTNAME           | IPV4 ADDRESS | IPV6 ADDRESS |
+-----------+------------------------------+--------------+--------------+
| * default | caas3-jf-mgmt-lb-1           | 44.130.4.7   |              |
+           +------------------------------+--------------+--------------+
|           | caas3-jf-mgmt-master-1       | 44.130.4.8   |              |
+           +------------------------------+--------------+--------------+
|           | caas3-jf-mgmt-worker-1       | 44.130.4.13  |              |
+           +------------------------------+--------------+--------------+
|           | caas3-jf-mgmt-worker-2       | 44.130.4.14  |              |
+           +------------------------------+--------------+--------------+
|           | caas3-jf-mgmt-worker-3       | 44.130.4.15  |              |
+-----------+------------------------------+--------------+--------------+
| tenant1   | caas3-openfaas-test-master-1 | 44.130.4.221 |              |
+           +------------------------------+--------------+--------------+
|           | caas3-openfaas-test-worker-1 | 44.130.4.222 |              |
+           +------------------------------+--------------+--------------+
|           | caas3-openfaas-test-worker-2 | 44.130.4.223 |              |
+           +------------------------------+--------------+--------------+
|           | caas3-openfaas-test-worker-3 | 44.130.4.224 |              |
+-----------+------------------------------+--------------+--------------+
| tenant2   | caas3-istio-test-master-1    | 44.130.4.225 |              |
+           +------------------------------+--------------+--------------+
|           | caas3-istio-test-worker-1    | 44.130.4.19  |              |
+           +------------------------------+--------------+--------------+
|           | caas3-istio-test-worker-2    | 44.130.4.20  |              |
+           +------------------------------+--------------+--------------+
|           | caas3-istio-test-worker-3    | 44.130.4.21  |              |
+-----------+------------------------------+--------------+--------------+
```
