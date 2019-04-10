package job

import (
	"context"
	"encoder-backend/pkg/models"
	"errors"
	"io"
	"path/filepath"
	"syscall"
)

type Encode struct {
	*options
	file   *models.File
	report Report
}

var (
	ErrCancelled = errors.New("cancelled encode job")
)

// New
// create a new encode job
func New(file *models.File) (*Encode, error) {

	e := &Encode{
		options: &options{},
		report:  Report{},
		file:    file,
	}

	opts := []option{
		withProfile(file.Path.QualityProfile),
	}

	for _, opt := range opts {
		err := opt(e.options)
		if err != nil {
			return nil, err
		}
	}

	return e, nil
}

// Run
// executes the jobs runtime command
func (e *Encode) Run(ctx context.Context) error {

	cmd, err := e.handbrake.Get(ctx, filepath.Join(e.file.Source, e.file.Name))
	if err != nil {
		return err
	}

	// note that we need this so that a call to ctrl+c (which kills the process group) will make the command
	// also exit instead of it being handled via. our own runtime

	/*if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			CreationFlags: 0x00000008,
		}
	}*/
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	read, write := io.Pipe()
	cmd.Stdout = write
	cmd.Stderr = write

	//defer write.Close()
	defer read.Close()

	// run the command
	err = cmd.Start()
	if err != nil {
		return err
	}

	done := make(chan error)

	go e.scan(ctx, read)

	go func() {
		done <- cmd.Wait()
	}()

	for {
		select {
		case <-ctx.Done():
			write.Write([]byte(scanQuit))
			return ErrCancelled
		case err := <-done:
			write.Write([]byte(scanQuit))
			return err
		}
	}
}

// Report
// obtain the current status report of the job
func (e *Encode) Report() Report {
	return e.report
}
