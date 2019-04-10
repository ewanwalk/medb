package worker

import (
	"encoder-backend/pkg/models"
	"time"
)

// onJobCancel
func (w *Worker) onJobCancel(err error) error {

	w.file.StatusEncoder = models.FileEncodeStatusNotDone
	w.file.Encodes[0].Status = models.EncodeCancelled
	w.file.Encodes[0].TimeEnd = time.Now().UTC()
	w.file.Encodes[0].Duration = int64(time.Since(w.file.Encodes[0].TimeStart) / time.Millisecond)
	w.file.Encodes[0].Error = err.Error()

	return w.db.Save(w.file).Error
}
