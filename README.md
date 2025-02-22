# Waypoint Plugin SDK

This repository is a Go library that enables users to write custom [Waypoint](https://waypointproject.io) plugins.
Waypoint supports plugins for builders, deployment platforms, release managers, and more. Waypoint plugins enable
Waypoint to utilize any popular service providers as well as custom in-house solutions.

Plugins in Waypoint are separate binaries which communicate with the Waypoint application; the plugin communicates using
gRPC, and while it is theoretically possible to build a plugin in any language supported by the gRPC framework. We
recommend that the developers leverage the [Waypoint SDK](https://github.com/hashicorp/waypoint-plugin-sdk).

## Simple Plugin Overview

To initialize the plugin SDK you use the `Main` function in the `sdk` package and register Components which provide
callbacks for Waypoint during the various parts of the lifecycle.

```go
package main

import (
  sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

func main() {

  sdk.Main(sdk.WithComponents(
  	// Comment out any components which are not
  	// required for your plugin
  	&builder.Builder{},
  ))

}
```

Components are Go structs which implement the various lifecycle interfaces in the Waypoint SDK. The following example
shows a plugin which Waypoint can use for the build lifecycle. Creating a `build` plugin is as simple as implementing
the interface `component.Builder` on a struct as shown in the following example.

```go
package builder

import (
  "context"
	"fmt"

  "github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

type Builder struct {}

func (b *Builder) BuildFunc() interface{} {
	// return a function which will be called by Waypoint
	return func(ctx context.Context, ui terminal.UI) (*Binary, error) {
	u := ui.Status()
	defer u.Close()
	u.Update("Building application")

	return &Binary{}, nil
  }
}
```

To pass values from one Waypoint component to another Waypoint component you return structs which implement
`proto.Message`, generally these structs are not manually created but generated by defining a Protocol Buffer message
template and then using the `protoc` and associated Go plugin to generate the Go structs.

The following example shows the Protocol Buffer message template which is used to define the `Binary` struct which is
returned from the previous example.

```go
syntax = "proto3";

package platform;

option go_package = "github.com/hashicorp/waypoint-plugin-examples/template/builder";

message Binary {
  string location = 1;
}
```

For more information on Protocol Buffers and Go, please see the documentation on the Google Developers website.

[https://developers.google.com/protocol-buffers/docs/gotutorial](https://developers.google.com/protocol-buffers/docs/gotutorial)

For full walkthrough for creating a Waypoint Plugin and reference documentation, please see the
[Extending Waypoint](https://www.waypointproject.io/docs/extending-waypoint) section of the Waypoint website.


## Example Plugins

Please see the following Plugins for examples of real world implementations of the Waypoint SDK.

### Build
[Docker](https://github.com/hashicorp/waypoint/tree/main/builtin/docker/builder.go)  
[Build Packs](https://github.com/hashicorp/waypoint/tree/main/builtin/pack/builder.go)

### Registry
[Docker](https://github.com/hashicorp/waypoint/tree/main/builtin/docker/registry.go)  
[Files](https://github.com/hashicorp/waypoint/tree/main/builtin/files/registry.go)

### Deploy
[Nomad](https://github.com/hashicorp/waypoint/tree/main/builtin/nomad/platform.go)  
[Kubernetes](https://github.com/hashicorp/waypoint/tree/main/builtin/k8s/platform.go)  
[Docker](https://github.com/hashicorp/waypoint/tree/main/builtin/docker/platform.go)  
[Azure Container Interface](https://github.com/hashicorp/waypoint/tree/main/builtin/azure/aci/platform.go)  
[Google Cloud Run](https://github.com/hashicorp/waypoint/tree/main/builtin/google/cloudrun/platform.go)  
[Netlify](https://github.com/hashicorp/waypoint/tree/main/builtin/netlify/platform.go)  
[Amazon EC2](https://github.com/hashicorp/waypoint/tree/main/builtin/aws/ec2/platform.go)

### Release
[Kubernetes](https://github.com/hashicorp/waypoint/tree/main/builtin/k8s/releaser.go)  
[Google Cloud Run](https://github.com/hashicorp/waypoint/tree/main/builtin/google/cloudrun/releaser.go)  
[Amazon ALB](https://github.com/hashicorp/waypoint/tree/main/builtin/aws/alb/releaser.go)

## Contributing

Thank you for your interest in contributing! Please refer to [CONTRIBUTING.md](https://github.com/hashicorp/waypoint-plugin-sdk/blob/master/.github/CONTRIBUTING.md) for guidance.
