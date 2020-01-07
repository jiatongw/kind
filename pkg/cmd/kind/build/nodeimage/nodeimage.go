/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nodeimage

import (
	"github.com/spf13/cobra"

	"sigs.k8s.io/kind/pkg/build/node"
	"sigs.k8s.io/kind/pkg/cmd"
	"sigs.k8s.io/kind/pkg/errors"
	"sigs.k8s.io/kind/pkg/log"
)

type flagpole struct {
	Source    string
	BuildType string
	Image     string
	BaseImage string
	KubeRoot  string
}

// NewCommand returns a new cobra.Command for building the node image
func NewCommand(logger log.Logger, streams cmd.IOStreams) *cobra.Command {
	flags := &flagpole{}
	cmd := &cobra.Command{
		Args: cobra.NoArgs,
		// TODO(bentheelder): more detailed usage
		Use:   "node-image",
		Short: "build the node image",
		Long:  "build the node image which contains kubernetes build artifacts and other kind requirements",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runE(logger, flags)
		},
	}
	cmd.Flags().StringVar(
		&flags.BuildType, "type",
		"docker", "build type, one of [bazel, docker]",
	)
	cmd.Flags().StringVar(
		&flags.Image, "image",
		node.DefaultImage,
		"name:tag of the resulting image to be built",
	)
	cmd.Flags().StringVar(
		&flags.KubeRoot, "kube-root",
		"",
		"path to the Kubernetes source directory (if empty, the path is autodetected)",
	)
	cmd.Flags().StringVar(
		&flags.BaseImage, "base-image",
		node.DefaultBaseImage,
		"name:tag of the base image to use for the build",
	)
	return cmd
}

func runE(logger log.Logger, flags *flagpole) error {
	// TODO(bentheelder): inject logger down the chain
	ctx, err := node.NewBuildContext(
		node.WithMode(flags.BuildType),
		node.WithImage(flags.Image),
		node.WithBaseImage(flags.BaseImage),
		node.WithKuberoot(flags.KubeRoot),
		node.WithLogger(logger),
	)
	if err != nil {
		return errors.Wrap(err, "error creating build context")
	}
	if err := ctx.Build(); err != nil {
		return errors.Wrap(err, "error building node image")
	}
	return nil
}
