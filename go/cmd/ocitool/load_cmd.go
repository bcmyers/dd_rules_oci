package main

import (
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

func LoadCmd(c *cli.Context) error {
	tarPath := c.String("file")

	f, err := os.Open(tarPath)
	if err != nil {
		return err
	}
	defer f.Close()

	cmd := exec.Command("docker", "image", "load")
	cmd.Stdin = f
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// TODO: Capture stdout and then do the tag?

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
