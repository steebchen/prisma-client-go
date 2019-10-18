package main

import (
	"github.com/prisma/photongo/jsonrpc"
)

// NewManifest generates a manifest for the prisma CLI
func NewManifest() jsonrpc.ManifestResponse {
	return jsonrpc.ManifestResponse{
		Manifest: jsonrpc.Manifest{
			DefaultOutput:      "./client_gen.go",
			PrettyName:         "Photon Go",
			Denylist:           []string{},
			RequiresGenerators: []string{},
			RequiresEngines:    []string{"queryEngine"},
		},
	}
}
