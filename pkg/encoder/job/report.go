package job

import (
	"encoder-backend/pkg/config"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

type Report struct {
	Percentage float32       `json:"percent"`
	FPS        float32       `json:"fps"`
	AvgFPS     float32       `json:"fps_avg"`
	Estimate   time.Duration `json:"eta"`
}

var (
	replacer = strings.NewReplacer(
		"Encoding: task 1 of 1, ", "",
		" % ", ", ",
		"(", "",
		")", "",
		"ETA ", "",
		" fps", "",
		"avg ", "",
	)
)

// Parse
// take an output line from handbrake and convert it into a format parsed into a mini-report
func (r *Report) Stream() chan<- string {

	logging := false
	if env := os.Getenv(config.EnvEncoderReportLogs); len(env) != 0 {
		logging = true
	}

	receiver := make(chan string)

	go func(recv <-chan string) {
		for line := range recv {

			if !strings.Contains(line, "Encoding: task 1 of 1") {
				continue
			}

			r.parse(line)

			if logging {
				log.WithField("report", *r).Debug("encode.report.stream: progress tick")
			}

		}
	}(receiver)

	return receiver
}

// parse
// take an output line from handbrake and convert it into a format parsed into a mini-report
func (r *Report) parse(line string) {

	formatted := strings.Split(replacer.Replace(line), ", ")
	if len(formatted) != 4 {
		return
	}

	/*
	 * FORMAT
	 * 0 - PERCENTAGE
	 * 1 - FPS
	 * 2 - AvgFPS
	 * 3 - ETA
	 */

	percent, err := strconv.ParseFloat(formatted[0], 32)
	if err != nil {
		return
	}
	r.Percentage = float32(percent)

	fps, err := strconv.ParseFloat(formatted[1], 32)
	if err != nil {
		return
	}
	r.FPS = float32(fps)

	avgfps, err := strconv.ParseFloat(formatted[2], 32)
	if err != nil {
		return
	}
	r.AvgFPS = float32(avgfps)

	eta, err := time.ParseDuration(formatted[3])
	if err != nil {
		return
	}
	r.Estimate = eta

}
