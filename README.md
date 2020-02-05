# Wrapper for Microsoft Application Insights library

This library wraps the official Go library for Application Insights, and adds Kubernetes meta data to the telemetry. It also sets the Cloud Role and Cloud Role Instance values to the deployment name/pod name. Setting these values makes the applications available in the Application Map view.

## Install

`$ go get github.com/ourstudio-se/go-appinsights-k8s`

## Usage

```go
package main

import (
	"os"

    "github.com/Microsoft/ApplicationInsights-Go/appinsights"
    "github.com/ourstudio-se/go-appinsights-k8s"
)

func main() {
	client := appink8s.NewTelemetryClient(os.Getenv("INSTRUMENTATION_KEY"))

	tm := appinsights.NewEventTelemetry("An event occurred!")
	client.Track(tm)
}
```

# License

MIT