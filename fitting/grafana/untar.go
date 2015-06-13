package grafana

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"github.com/nathanielc/morgoth/Godeps/_workspace/src/github.com/golang/glog"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Untar a `sourceTar` into `dir`
// `sourceTar` may be gzipped if it ends in a .gz extension
// `dir` will be created if it doesn't exist
func untar(sourceTar, dir string) error {

	if len(sourceTar) == 0 {
		return errors.New("Must pass valid sourceTar.")
	}

	if len(dir) == 0 {
		return errors.New("Must pass valid dir.")
	}

	file, err := os.Open(sourceTar)

	if err != nil {
		return err
	}

	defer file.Close()

	var fileReader io.ReadCloser = file

	// Handle gzipped tar
	if strings.HasSuffix(sourceTar, ".gz") {
		if fileReader, err = gzip.NewReader(file); err != nil {
			return err
		}
		defer fileReader.Close()
	}

	tarBallReader := tar.NewReader(fileReader)

	// Extracting tarred files
	glog.V(2).Infof("Extracting tar: %s into %s", sourceTar, dir)
	for {
		header, err := tarBallReader.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		// Remove parent dir from tar and replace with `dir`
		fileparts := strings.Split(header.Name, string(filepath.Separator))
		filename := path.Join(dir, path.Join(fileparts[1:]...))

		switch header.Typeflag {
		case tar.TypeDir:
			// Handle directory
			glog.V(3).Info("Creating directory :", filename)
			err = os.MkdirAll(filename, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

		case tar.TypeReg:
			// Handle normal file
			glog.V(3).Info("Untarring :", filename)
			filedir := path.Dir(filename)
			err = os.MkdirAll(filedir, 0755)
			if err != nil {
				return err
			}
			writer, err := os.Create(filename)
			if err != nil {
				return err
			}

			io.Copy(writer, tarBallReader)

			err = os.Chmod(filename, os.FileMode(header.Mode))

			if err != nil {
				return err
			}
			writer.Close()

		default:
			glog.Error("Unable to untar type : %c in file %s", header.Typeflag, filename)
		}
	}
	return nil
}
