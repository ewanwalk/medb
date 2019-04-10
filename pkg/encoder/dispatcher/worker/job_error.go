package worker

import (
	"encoder-backend/pkg/models"
	"time"
)

// onJobError
func (w *Worker) onJobError(err error) error {

	w.file.StatusEncoder = models.FileEncodeStatusErrored
	w.file.Encodes[0].Status = models.EncodeErrored
	w.file.Encodes[0].TimeEnd = time.Now().UTC()
	w.file.Encodes[0].Duration = int64(time.Since(w.file.Encodes[0].TimeStart) / time.Millisecond)
	w.file.Encodes[0].Error = err.Error()

	return w.db.Save(w.file).Error
}
