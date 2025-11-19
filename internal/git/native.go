package git

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

// isGitInstalled indicates if the git command is installed.
// TODO: check that version supports git archive
func isGitInstalled() bool {
	_, err := exec.LookPath("git")

	// TODO: check version / capabilities and cache result
	return err == nil
}

func (r *Repository) nativeExtractGitArchive(ctx context.Context, w io.Writer, file string, selectedRef *Ref) (err error) {
	// attention credential auth etc
	/*
		git archive --remote=$REPO_URL HEAD:path/to -- file.xz |
		tar xO > /where/you/want/to/have.it
	*/
	hash := selectedRef.Hash()
	args := []string{"archive",
		"--format=tgz",
		fmt.Sprintf("--remote=%v", r.repoURL),
		fmt.Sprintf("%s:%s", hash, file),
	}
	r.debug("running git %s", strings.Join(args, " "))
	cmd := exec.CommandContext(ctx, "git", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	r.debug("got stdout pipe")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		r.debug("cmd error: %v", err)

		return err
	}

	defer func() {
		r.debug("closing command")
		const maxErrSize = 2048
		var b bytes.Buffer
		// drain stderr and capture
		_, _ = io.CopyN(&b, stderr, maxErrSize)
		_, _ = io.Copy(io.Discard, stderr)

		if err != nil {
			r.debug("early exit with error: %v", err)
			// drain command output on early exit
			_, _ = io.Copy(io.Discard, stdout)
		}

		errCommand := cmd.Wait()
		switch {
		case err == nil && errCommand != nil:
			err = errCommand
			if b.Len() > 0 {
				err = errors.Join(errCommand, errors.New(b.String()))
			}
		case err != nil && errCommand == nil:
			if b.Len() > 0 {
				err = errors.Join(err, errors.New(b.String()))
			}
		case err != nil && errCommand != nil:
			err = errors.Join(err, errCommand)
			if b.Len() > 0 {
				err = errors.Join(err, errors.New(b.String()))
			}
		case err == nil && errCommand == nil:
			fallthrough
		default:
		}
		log.Printf("DEBUG: %s", b.String())
	}()
	r.debug("cmd running in the background")

	// copy piped tgz to writer
	gzipReader, err := gzip.NewReader(stdout)
	if err != nil {
		return err
	}
	defer func() {
		_ = gzipReader.Close()
	}()
	r.debug("got gzip reader")

	tarReader := tar.NewReader(gzipReader)
	r.debug("got tar reader")

	r.debug("starting command")

	r.debug("reading tar")
	for {
		_, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			r.debug("tar read error: %v", err)
			break
		}

		_, err = io.Copy(w, tarReader)
		if err != nil {
			break
		}
	}

	r.debug("end of reading err=%v", err)
	return err
}
