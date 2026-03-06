package main

import (
	"context"
	"flag"
	"log"

	tfprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/srmullaney/terraform-provider-duo/internal/provider"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/srmullaney/duo",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), func() tfprovider.Provider {
		return provider.New("1.0.0")
	}, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
