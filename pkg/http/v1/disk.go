package v1

import (
	"encoder-backend/pkg/models"
	"github.com/ewanwalk/respond"
	"github.com/shirou/gopsutil/disk"
	"net/http"
	"sort"
	"strings"
)

type DiskUsageReport struct {
	Device disk.PartitionStat `json:"device"`
	Usage  disk.UsageStat     `json:"usage"`
}

func getDiskUsage(w http.ResponseWriter, r *http.Request) {

	stats, err := disk.Partitions(true)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	usage := make([]DiskUsageReport, 0)
	devices := make(map[string]disk.PartitionStat)
	paths := make([]models.Path, 0)

	if err := db.Where("status = ? OR type = ?", models.PathStatusEnabled, models.PathTypePseudo).Find(&paths).Error; err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	for _, path := range paths {

		best := disk.PartitionStat{}
		last := 0

		for _, stat := range stats {
			if stat.Device == "rootfs" || stat.Device == "/dev/root" {
				continue
			}

			if !strings.Contains(path.Directory, stat.Mountpoint) {
				continue
			}

			if len(stat.Mountpoint) > last {
				last = len(stat.Mountpoint)
				best = stat
			}
		}

		devices[best.Mountpoint] = best
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

		usage = append(usage, DiskUsageReport{
			Device: stat,
			Usage:  *used,
		})
	}

	sort.Slice(usage, func(i, j int) bool {
		return usage[i].Device.Mountpoint > usage[j].Device.Mountpoint
	})

	respond.With(w, r, http.StatusOK, usage)
}

func getDiskUsageDebug(w http.ResponseWriter, r *http.Request) {

	stats, err := disk.Partitions(true)
	if err != nil {
		respond.With(w, r, http.StatusInternalServerError, err)
		return
	}

	usage := make([]disk.UsageStat, len(stats))

	for _, stat := range stats {
		used, err := disk.Usage(stat.Mountpoint)
		if err != nil || len(used.Path) == 0 || used.Total == 0 {
			continue
		}

		usage = append(usage, *used)
	}

	respond.With(w, r, http.StatusOK, map[string]interface{}{
		"_":     stats,
		"usage": usage,
	})
}
