// Referenced from: https://github.com/mraron/njudge/
package memory

import (
	"fmt"
	"strconv"
)

type Memory int64

const (
	Byte Memory = 1
	KB          = 1000 * Byte
	KiB         = 1024 * Byte
	MB          = 1000 * KB
	MiB         = 1024 * KiB
	GB          = 1000 * MB
	GiB         = 1024 * MiB
)

func (m *Memory) MarshalJSON() ([]byte, error) {
	// return []byte(fmt.Sprintf("\"%dB\"", *m)), nil
	return fmt.Appendf(nil, "\"%dB\"", *m), nil
}

func (m *Memory) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	// Not a string => `val` in KB -> multiple by KiB
	if b[0] != '"' {
		tmp, err := strconv.Atoi(string(b))
		if err != nil {
			return err
		}
		*m = Memory(tmp) * KiB
		return nil
	}

	// Is a string => `val`B in Bytes
	tmp, err := strconv.Atoi(string(b[1 : len(b)-2]))
	if err != nil {
		return err
	}
	*m = Memory(tmp)
	return nil
}
