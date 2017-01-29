package blamewarrior

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Rating struct {
	RepoName string
	Val      int
}

func (rating Rating) Value() (driver.Value, error) {
	j, err := json.Marshal(rating)
	return j, err
}

func (rating *Rating) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	err := json.Unmarshal(source, &rating)
	if err != nil {
		return err
	}

	return nil
}
