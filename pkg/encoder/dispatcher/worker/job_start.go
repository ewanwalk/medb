package worker

import (
	"encoder-backend/pkg/models"
	"time"
)

// onJobStart
func (w *Worker) onJobStart() error {

	checksum, err := w.file.CurrentChecksum()
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	w.file.StatusEncoder = models.FileEncodeStatusRunning
	w.file.Encodes = []models.Encode{
		{
			QualityProfileID: w.file.Path.QualityProfileID,
			FileID:           w.file.ID,
			TimeStart:        &now,
			ChecksumAtStart:  checksum,
			NameAtStart:      w.file.Name,
			SizeAtStart:      w.file.Size,
			Status:           models.EncodeRunning,
		},
	}

	return w.db.Save(w.file).Error
}
