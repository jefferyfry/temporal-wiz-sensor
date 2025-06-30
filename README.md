# Temporal Wiz Sensor Activity

Reusable Wiz Sensor activity for k8s to be used in your Temporal workflows.

## Features
- Install and validate Wiz Sensor on Kubernetes clusters using Temporal Workflows.
- Can be used in Temporal workflows to ensure the sensor is installed before proceeding with other activities.
- Packaged as a Go Module

## Installation
To use the Wiz Sensor in your Temporal workflows, you can install it as a Go module. Run the following command in your project directory:

```bash
go get github.com/jefferyfry/temporal-wiz-sensor@v0.1.0
```

## Usage
See main-reg.go for an example of how to register the workflow and activity with Temporal.
See main-client.go for an example of how to execute the workflow and pass in the required parameters.