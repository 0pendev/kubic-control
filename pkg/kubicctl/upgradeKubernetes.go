// Copyright 2019 Thorsten Kukuk
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubicctl

import (
	"context"
	"time"
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	pb "github.com/thkukuk/kubic-control/api"
)

func UpgradeKubernetesCmd() *cobra.Command {
        var subCmd = &cobra.Command {
                Use:   "upgrade",
                Short: "Upgrade Kubernetes Cluster",
                Run: upgradeKubernetes,
		Args: cobra.ExactArgs(0),
	}

	return subCmd
}

func upgradeKubernetes(cmd *cobra.Command, args []string) {
	// Set up a connection to the server.

	conn, err := CreateConnection()
	if err != nil {
		return
	}
	defer conn.Close()

	c := pb.NewKubeadmClient(conn)

	var deadlineMin = flag.Int("deadline_min", 20, "Default deadline in minutes.")
	clientDeadline := time.Now().Add(time.Duration(*deadlineMin) * time.Minute)
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Print ("Upgrading kubernetes can take a very long time, please be patient.\n")
	r, err := c.UpgradeKubernetes(ctx, &pb.Version{Version: ""})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not upgrade: %v", err)
		return
	}
	if r.Success {
		fmt.Printf("Kubernetes cluster was successfully upgraded to version %s", r.Message)
	} else {
		fmt.Fprintf(os.Stderr, "Upgrading kubernetes failed: %s", r.Message)
	}
}