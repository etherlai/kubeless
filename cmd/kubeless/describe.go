/*
Copyright (c) 2016-2017 Bitnami

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

package main

import (
	"encoding/json"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/gosuri/uitable"
	"github.com/kubeless/kubeless/pkg/spec"
	"github.com/kubeless/kubeless/pkg/utils"
	"github.com/spf13/cobra"
	"k8s.io/client-go/pkg/api"
)

var describeCmd = &cobra.Command{
	Use:     "describe FLAG",
	Aliases: []string{"ls"},
	Short:   "describe a function deployed to Kubeless",
	Long:    `describe a function deployed to Kubeless`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			logrus.Fatal("Need exactly one argument - function name")
		}
		funcName := args[0]

		ns, err := cmd.Flags().GetString("namespace")
		if err != nil {
			logrus.Fatalf("Can not describe function: %v", err)
		}

		output, err := cmd.Flags().GetString("out")
		if err != nil {
			logrus.Fatalf("Can not describe function: %v", err)
		}

		f, err := utils.GetFunction(funcName, ns)
		if err != nil {
			logrus.Fatalf("Can not describe function: %v", err)
		}

		print(f, funcName, output)
	},
}

func init() {
	describeCmd.Flags().StringP("out", "o", "", "Output format. One of: json|yaml")
	describeCmd.Flags().StringP("namespace", "", api.NamespaceDefault, "Specify namespace for the function")
}

func print(f spec.Function, name, output string) {
	switch output {
	case "":
		table := uitable.New()
		table.MaxColWidth = 80
		table.Wrap = true
		label, _ := json.Marshal(f.Metadata.Labels)
		env, _ := json.Marshal(f.Spec.Template.Spec.Containers[0].Env)
		table.AddRow("Name:", name)
		table.AddRow("Namespace:", fmt.Sprintf(f.Metadata.Namespace))
		table.AddRow("Handler:", fmt.Sprintf(f.Spec.Handler))
		table.AddRow("Runtime:", fmt.Sprintf(f.Spec.Runtime))
		table.AddRow("Type:", fmt.Sprintf(f.Spec.Type))
		table.AddRow("Topic:", fmt.Sprintf(f.Spec.Topic))
		table.AddRow("Label:", fmt.Sprintf(string(label)))
		table.AddRow("Envvar:", fmt.Sprintf(string(env)))
		table.AddRow("Memory:", fmt.Sprintf(f.Spec.Template.Spec.Containers[0].Resources.Requests.Memory().String()))
		table.AddRow("Dependencies:", fmt.Sprintf(f.Spec.Deps))
		fmt.Println(table)
	case "json":
		b, _ := json.MarshalIndent(f, "", "  ")
		fmt.Println(string(b))
	case "yaml":
		b, _ := yaml.Marshal(f)
		fmt.Println(string(b))
	default:
		fmt.Println("Wrong output format. Please use only json|yaml")
	}
}
