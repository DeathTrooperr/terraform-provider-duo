package main

import (
	"context"
	"flag"
	"log"

	"github.com/DeathTrooperr/terraform-provider-duo/internal/provider"
	tfprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// these will be set by the build process
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/DeathTrooperr/duo",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), func() tfprovider.Provider {
		return provider.New(version)
	}, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
