// Copyright © 2019 Máté Birkás <gadfly16@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gradient-images/teflon/internal/teflon"
	// "gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Reads Teflon metadata",
	Long: `'teflon get' prints the metadata belonging to the targets.  If no <target>
is specified it will run for '.'. As a side effect the metadata access the
command creates the meta files if they are not already exist or refreshes
them if they are not up-to-date.`,
	Run: getRun,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func getRun(cmd *cobra.Command, args []string) {
	log.Println("'get' command called")
	if len(args) == 0 {
		args = append(args, ".")
		log.Println("No targets given, running for '.' .")
	}
	for _, target := range args {
		o, err := teflon.NewObject(target)
		if err != nil {
			log.Fatalln("Couldn't create object:", err)
		}
		err = o.InitMeta()
		if err != nil {
			log.Fatalln("Couldn't get metadata.", err)
		}
		d, err := json.MarshalIndent(&o, "", "  ")
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		fmt.Println(string(d))
		// for k, v := range m.UserData {
		// 	fmt.Println(k+":", v)
		// }
	}
}