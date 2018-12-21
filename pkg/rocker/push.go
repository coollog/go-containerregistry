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
	  "github.com/google/go-containerregistry/pkg/v1/tarball"
  "github.com/google/go-containerregistry/pkg/v1/remote"
  "net/http"
    "os"
  "bytes"
      "io"
    "io/ioutil"
  "github.com/google/go-containerregistry/pkg/v1/empty"
  "github.com/google/go-containerregistry/pkg/v1/mutate"
  "archive/tar"
  "fmt"
  "errors"
)

func init() { Root.AddCommand(NewCmdPush()) }

func NewCmdPush() *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "Push something to a registry",
		Args:  cobra.ExactArgs(1),
		Run:   push,
	}
}

func push(_ *cobra.Command, args []string) {
	dst := args[0]
	t, err := name.NewTag(dst, name.WeakValidation)
	if err != nil {
		log.Fatalf("parsing tag %q: %v", dst, err)
	}
	log.Printf("Pushing %v", t)

	auth, err := authn.DefaultKeychain.Resolve(t.Registry)
	if err != nil {
		log.Fatalf("getting creds for %q: %v", t, err)
	}

	dataLayerBuf, err := makeLayerTar()
	if err != nil {
    log.Fatal("", err)
  }

  // Append the new layer.
  dataLayerBytes := dataLayerBuf.Bytes()
  dataLayer, err := tarball.LayerFromOpener(func() (io.ReadCloser, error) {
    return ioutil.NopCloser(bytes.NewBuffer(dataLayerBytes)), nil
  })
  if err != nil {
    log.Fatal("", err)
  }

  // Augment the base image with our data layer.
  i, err := mutate.AppendLayers(empty.Image, dataLayer)
  if err != nil {
    log.Fatal("", err)
  }

	if err := remote.Write(t, i, auth, http.DefaultTransport); err != nil {
		log.Fatalf("writing image %q: %v", t, err)
	}
}

// Makes a layer tarball with dataFile at path /data.
func makeLayerTar() (*bytes.Buffer, error) {
  // Save stdin as temp file
  file, err := ioutil.TempFile("", "")
  if err != nil {
    return nil, err
  }
  defer os.Remove(file.Name())

  size, err := io.Copy(file, os.Stdin)
  if err != nil {
    return nil, err
  }

  // Write layer with one file /data
  dataLayerBuf := bytes.NewBuffer(nil)
  tw := tar.NewWriter(dataLayerBuf)
  defer tw.Close()

  // Copy the file into the image tarball.
  if err := tw.WriteHeader(&tar.Header{
   Name:     "data",
   Size:     size,
   Typeflag: tar.TypeReg,
   // Use a fixed Mode, so that this isn't sensitive to the directory and umask
   // under which it was created. Additionally, windows can only set 0222,
   // 0444, or 0666, none of which are executable.
   Mode: 0555,
  }); err != nil {
   return nil, err
  }

  file.Seek(0, 0)
  writtenSize, err := io.Copy(tw, file)
  if err != nil {
    return nil, err
  }
  if writtenSize != size {
    return nil, errors.New(fmt.Sprintf("expected to write %d bytes but wrote %d bytes", size, writtenSize))
  }
  return dataLayerBuf, nil
}
