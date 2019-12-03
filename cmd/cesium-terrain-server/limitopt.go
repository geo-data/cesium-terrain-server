package main

import (
	"errors"
	"fmt"
	"github.com/nmccready/cesium-terrain-server/handlers"
	"strconv"
)

// Adapted from <https://golang.org/doc/effective_go.html#constants>.
type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
)

func (b ByteSize) String() string {
	switch {
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fkB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}

func ParseByteSize(size string) (bytes ByteSize, err error) {
	defer func() {
		if bytes < 0 {
			err = errors.New("size cannot be negative")
		}
	}()

	val, err := strconv.ParseFloat(size, 64)
	if err == nil {
		bytes = ByteSize(val)
		return
	}

	if len(size) < 3 {
		err = errors.New("the size must be specified as a suffix e.g 5MB")
		return
	}

	val, err = strconv.ParseFloat(size[:len(size)-2], 64)
	if err != nil {
		return
	}
	bytes = ByteSize(val)

	suffix := size[len(size)-2:]
	switch suffix {
	case "TB":
		bytes *= TB
	case "GB":
		bytes *= GB
	case "MB":
		bytes *= MB
	case "KB":
		bytes *= KB
	default:
		err = errors.New("bad size suffix: " + suffix)
	}
	return
}

type LimitOpt struct {
	Value handlers.Bytes
}

func NewLimitOpt() *LimitOpt {
	return &LimitOpt{}
}

func (this *LimitOpt) String() string {
	return ByteSize(this.Value).String()
}

func (this *LimitOpt) Set(size string) error {
	byteSize, err := ParseByteSize(size)
	if err != nil {
		return err
	}

	this.Value = handlers.Bytes(byteSize)

	return nil
}
