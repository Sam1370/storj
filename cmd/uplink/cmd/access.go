// Copyright (C) 2020 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	libuplink "storj.io/storj/lib/uplink"
	"storj.io/storj/pkg/cfgstruct"
	"storj.io/storj/pkg/process"
)

var inspectCfg AccessConfig

func init() {
	// We skip the use of addCmd here because we only want the configuration options listed
	// above, and addCmd adds a whole lot more than we want.
	accessCmd := &cobra.Command{
		Use:   "access",
		Short: "Set of commands to manage access.",
	}

	inspectCmd := &cobra.Command{
		Use:   "inspect [ACCESS]",
		Short: "Inspect allows you to explode a serialized access into it's constituent parts.",
		RunE:  accessInspect,
		Args:  cobra.MaximumNArgs(1),
	}
	RootCmd.AddCommand(accessCmd)
	accessCmd.AddCommand(inspectCmd)

	process.Bind(inspectCmd, &inspectCfg, defaults, cfgstruct.ConfDir(getConfDir()))
}

func accessInspect(cmd *cobra.Command, args []string) (err error) {
	var access *libuplink.Scope
	if len(args) == 0 {
		access, err = inspectCfg.GetAccess()
		if err != nil {
			return err
		}
	} else {
		firstArg := args[0]

		access, err = inspectCfg.GetNamedAccess(firstArg)
		if err != nil {
			return err
		}

		if access == nil {
			if access, err = libuplink.ParseScope(firstArg); err != nil {
				return err
			}
		}
	}

	serializedAPIKey := access.APIKey.Serialize()
	serializedEncAccess, err := access.EncryptionAccess.Serialize()
	if err != nil {
		return err
	}

	fmt.Println("=========== ACCESS INFO ==================================================================")
	fmt.Println("Satellite        :", access.SatelliteAddr)
	fmt.Println("API Key          :", serializedAPIKey)
	fmt.Println("Encryption Access:", serializedEncAccess)
	return nil
}
