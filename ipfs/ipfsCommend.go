package ipfs

import "main/utils"

func RunIPFSCommend(args []string) error {
	err := utils.RunCommend("ipfs", args, "")
	if err != nil {
		return err
	}
	return nil
}
