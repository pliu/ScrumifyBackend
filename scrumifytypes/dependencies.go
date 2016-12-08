package scrumifytypes

import (
    "database/sql/driver"
    "bytes"
    "encoding/gob"
    "ScrumifyBackend/utils"
)

type Dependencies struct {
    Dependencies []int64
    Valid bool
}

func (d *Dependencies)Scan(value interface{}) error {
    if value == nil {
        d.Dependencies = []int64{}
        d.Valid = false
    } else {
        if check := fromBytes(value.([]byte)); check != nil {
            d.Dependencies = check
            d.Valid = true
        } else {
            d.Dependencies = []int64{}
            d.Valid = false
        }
    }
    return nil
}

func (d Dependencies)Value() (driver.Value, error) {
    if d.Valid {
        return d.toBytes(), nil
    } else {
        return nil, nil
    }
}

func (d Dependencies)toBytes() []byte {
    b := bytes.Buffer{}
    e := gob.NewEncoder(&b)
    if err := e.Encode(d.Dependencies); err != nil {
        utils.PrintErr(err, "Could not serialize Dependencies")
        return nil
    }
    return b.Bytes()
}

func fromBytes(b []byte) []int64 {
    out := []int64{}
    buffer := bytes.NewBuffer(b)
    decoder := gob.NewDecoder(buffer)
    if err := decoder.Decode(&out); err == nil {
        return out
    } else {
        utils.PrintErr(err, "Could not deserialize Dependencies")
        return nil
    }
}

func (d Dependencies)IsValid() bool {
    return !d.Valid || len(d.Dependencies) < 1000
}
