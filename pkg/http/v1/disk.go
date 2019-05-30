package v1

import (
	"encoder-backend/pkg/models"
	"github.com/ewanwalk/respond"
	"github.com/shirou/gopsutil/disk"
	"net/http"
	"strings"
)

func getDiskUsage(w http.ResponseWriter, r *http.Request) {

	stats, err := disk.Partitions(true)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	usage := make([]interface{}, 0)
	devices := make(map[string]disk.PartitionStat)
	paths := make([]models.Path, 0)

	if err := db.Scopes(models.PathEnabled).Find(&paths).Error; err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	for _, path := range paths {
		for _, stat := range stats {
			if stat.Device == "rootfs" || !strings.Contains(path.Directory, stat.Mountpoint) {
				continue
			}

			devices[stat.Mountpoint] = stat
			break
		}
	}

	for _, stat := range devices {

		used, err := disk.Usage(stat.Mountpoint)
		if err != nil {

			continue
		}

		if used.Used == 0 {
			used.Used = used.Total - used.Free
			used.UsedPercent = (float64(used.Used) / float64(used.Total)) * 100
		}

		usage = append(usage, map[string]interface{}{
			"device": stat,
			"usage":  used,
		})
	}

	respond.With(w, r, http.StatusOK, usage)
}
