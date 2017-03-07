package models

import (
	"crypto/rand"
	"fmt"
	"hash/fnv"
	"io"
	"reflect"
	"time"
)

type SyrosService struct {
	Id          string            `gorethink:"id,omitempty" json:"id"`
	Hostname    string            `gorethink:"hostname" json:"hostname"`
	Type        string            `gorethink:"type" json:"type"`
	Config      map[string]string `gorethink:"config" json:"config"`
	Environment string            `gorethink:"environment" json:"environment"`
	Collected   time.Time         `gorethink:"collected" json:"collected"`
}

// ConfigToMap converts a config struct to a map using the m tags
func ConfigToMap(in interface{}, tag string) (map[string]string, error) {
	out := make(map[string]string)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ConfigToMap only accepts structs got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		if mtag := fi.Tag.Get(tag); mtag != "" {
			out[mtag] = fmt.Sprint(v.Field(i).Interface())
		}
	}
	return out, nil
}

func Hash(val string) string {
	h := fnv.New32a()
	h.Write([]byte(val))
	return fmt.Sprint(h.Sum32())
}

func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
