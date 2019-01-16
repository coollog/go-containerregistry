// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rocker

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
		"github.com/google/go-containerregistry/pkg/v1"
	"archive/tar"
	"io"
		"errors"
			"os"
	"fmt"
)

func init() { Root.AddCommand(NewCmdPull()) }

func NewCmdPull() *cobra.Command {
	return &cobra.Command{
		Use:   "pull [reference]",
		Short: "Pull something from a registry",
		Args:  cobra.ExactArgs(1),
		Run:   pull,
	}
}

func pull(_ *cobra.Command, args []string) {
	src := args[0]

	ref, err := name.ParseReference(src, name.WeakValidation)
	if err != nil {
		log.Fatalf("parsing tag %q: %v", src, err)
	}
	log.Printf("Pulling %v", ref)

	i, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		log.Fatalf("reading image %q: %v", ref, err)
	}

	if err := writeLayerData(i); err != nil {
		log.Fatalf("writing layer data: %v", err)
	}
}

func writeLayerData(image v1.Image) error {
	layers, err := image.Layers()
	if err != nil {
		return err
	}

	if len(layers) != 1 {
		return errors.New(fmt.Sprintf("expected 1 layer, but got %d", len(layers)))
	}

	dataLayer := layers[0]
	layerReader, err := dataLayer.Uncompressed()
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(layerReader)
	header, err := tarReader.Next()
	if err != nil {
		return err
	}
	if header.Name != "data" {
		return errors.New(fmt.Sprintf("expected /data, but got %s", header.Name))
	}

	if _, err := io.Copy(os.Stdout, tarReader); err != nil {
		return err
	}
	return nil
}
