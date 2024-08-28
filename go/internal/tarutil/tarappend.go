package tarutil

import (
	"archive/tar"
	"io"
	"os"
	"os/user"
	"strconv"
	"strings"
	"time"
)

// AppendFileToTarWriter appends a file (given as a filepath) to a tarfile
// through the tarfile interface.
func AppendFileToTarWriter(filePath string, loc string, tw *tar.Writer) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}

	hdr.ChangeTime = time.Time{}
	hdr.ModTime = time.Time{}
	hdr.AccessTime = time.Time{}

	hdr.Name = loc

	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := io.Copy(tw, f); err != nil {
		return err
	}

	return nil
}

func AppendStringToTarWriter(s string, loc string, tw *tar.Writer) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		return err
	}

	hdr := &tar.Header{
		Name:     loc,
		Size:     int64(len(s)),
		Typeflag: tar.TypeReg,

		Gid:  gid,
		Mode: 0660,
		Uid:  uid,

		AccessTime: time.Time{},
		ChangeTime: time.Time{},
		ModTime:    time.Time{},
	}

	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}

	if _, err := io.Copy(tw, strings.NewReader(s)); err != nil {
		return err
	}

	return nil
}
