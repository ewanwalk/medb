package worker

import (
	"encoder-backend/pkg/models"
	"os"
	"path/filepath"
	"time"
)

// onJobEnd
func (w *Worker) onJobEnd() error {

	absFile := filepath.Join(w.file.Source, w.file.Name)
	file, err := os.Stat(absFile)
	if err != nil {
		return err
	}

	checksum, err := w.file.CurrentChecksum()
	if err != nil {
		return err
	}

	w.file.StatusEncoder = models.FileEncodeStatusDone
	w.file.Size = file.Size()
	w.file.Checksum = checksum

	w.file.Encodes[0].Status = models.EncodeDone
	w.file.Encodes[0].ChecksumAtEnd = checksum
	w.file.Encodes[0].NameAtEnd = w.file.Name
	w.file.Encodes[0].SizeAtEnd = file.Size()
	w.file.Encodes[0].TimeEnd = time.Now().UTC()
	w.file.Encodes[0].Duration = int64(time.Since(w.file.Encodes[0].TimeStart) / time.Millisecond)

	return w.db.Save(w.file).Error
}
