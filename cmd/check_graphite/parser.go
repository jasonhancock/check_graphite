package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
)

type Target struct {
	Datapoints []point `json:"datapoints"`
	Target     string  `json:"target"`
}

type point []*float64

func parseGraphiteResponse(r io.Reader, metric string) (float64, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	if len(body) == 0 {
		return 0, errors.New("no data read from reader")
	}

	var targets []Target

	err = json.Unmarshal(body, &targets)
	if err != nil {
		return 0, err
	}

	for _, v := range targets {
		if v.Target != metric {
			continue
		}

		for i := len(v.Datapoints) - 1; i >= 0; i-- {
			if v.Datapoints[i][0] != nil {
				return *v.Datapoints[i][0], nil
			}
		}

		return 0, errors.New("unable to determine a value for metric")
	}

	return 0, errors.New("metric not found in response")
}
