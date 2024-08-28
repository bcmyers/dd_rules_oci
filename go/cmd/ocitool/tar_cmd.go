package main

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/DataDog/rules_oci/go/internal/tarutil"
	"github.com/urfave/cli/v2"
)

func TarCmd(c *cli.Context) error {
	blobIndexPath := c.String("blob-index")
	descriptorPath := c.String("descriptor-file")
	paths := c.StringSlice("file")
	isGzip := c.String("gzip")
	outPath := c.String("out")

	fmt.Println("TODO:", isGzip)

	// Create a tar file

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("Error creating file %s: %w", outPath, err)
	}
	defer f.Close()

	gzipWriter := gzip.NewWriter(f)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	// Append oci-layout to tar

	err = tarutil.AppendStringToTarWriter(
		"{\"imageLayoutVersion\": \"1.0.0\"}",
		"oci-layout",
		tarWriter,
	)
	if err != nil {
		return fmt.Errorf("Error appending oci-layout to tar: %w", err)
	}

	// Append index.json to tar

	indexDotJson, err := createIndexDotJson(descriptorPath)
	if err != nil {
		return err
	}
	err = tarutil.AppendStringToTarWriter(
		indexDotJson,
		"index.json",
		tarWriter,
	)
	if err != nil {
		return fmt.Errorf("Error appending index.json to tar: %w", err)
	}

	// Append blobs to tar

	blobs, err := getBlobs(blobIndexPath)
	if err != nil {
		return fmt.Errorf("TODO: %w", err)
	}

	done := make(map[string]struct{})
	for _, b := range blobs {
		newPath := fmt.Sprintf("blobs/%s", b.digest)
		err = tarutil.AppendFileToTarWriter(b.path, newPath, tarWriter)
		if err != nil {
			return fmt.Errorf("Error appending %s to tar: %w", newPath, err)
		}
		done[b.path] = struct{}{}
	}

	// Append more blobs to tar

	blobs = []blob_{}
	for _, p := range paths {
		if _, ok := done[p]; ok {
			continue
		}
		digest, err := getSha256(p)
		if err != nil {
			return fmt.Errorf("Error getting sha256 of %s: %w", p, err)
		}
		blobs = append(
			blobs,
			blob_{
				path: p, digest: digest,
			},
		)
	}

	for _, b := range blobs {
		newPath := fmt.Sprintf("blobs/sha256/%s", b.digest)
		err = tarutil.AppendFileToTarWriter(b.path, newPath, tarWriter)
		if err != nil {
			return fmt.Errorf("Error appending %s to tar: %w", newPath, err)
		}
	}

	return nil
}

func createIndexDotJson(descriptorPath string) (string, error) {
	b, err := os.ReadFile(descriptorPath)
	if err != nil {
		return "", err
	}

	var i interface{}
	err = json.Unmarshal(b, &i)
	if err != nil {
		return "", err
	}

	index := &struct {
		SchemaVersion int    `json:"schemaVersion"`
		MediaType     string `json:"mediaType"`
		Manifests     []interface{}
	}{
		SchemaVersion: 2,
		MediaType:     "application/vnd.oci.image.index.v1+json",
		Manifests:     []interface{}{i},
	}
	b, err = json.Marshal(index)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

type blob_ struct {
	path   string
	digest string
}

func getBlobs(blobIndexPath string) ([]blob_, error) {
	f, err := os.Open(blobIndexPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	type blobIndex struct {
		Blobs map[string]string `json:"Blobs"`
	}
	var bi blobIndex
	err = json.Unmarshal(b, &bi)
	if err != nil {
		return nil, err
	}

	var blobs []blob_
	for digest, path := range bi.Blobs {
		blobs = append(blobs, blob_{path, digest})
	}

	return blobs, nil
}

func getSha256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := sha256.New()

	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	b := hash.Sum(nil)

	s := fmt.Sprintf("%x", b)

	return s, nil
}
