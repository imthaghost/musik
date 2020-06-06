package joiner

import (
	"os"
	"sync"
	"time"
)

// Joiner ...
type Joiner struct {
	l      sync.Mutex
	blocks map[int][]byte
	file   *os.File
	name   string
}

// New ...
func New(name string) (*Joiner, error) {
	f, err := os.OpenFile("assets/music/"+name, os.O_CREATE|os.O_TRUNC|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return nil, err
	}

	joiner := &Joiner{
		blocks: map[int][]byte{},
		file:   f,
		name:   name,
	}

	return joiner, nil

}

// Join ...
func (j *Joiner) Join(id int, block []byte) {
	j.l.Lock()
	j.blocks[id] = block
	j.l.Unlock()
}

// Run ...
func (j *Joiner) Run(count int) error {
	var index = 0
	for index < count {
		j.l.Lock()
		block := j.blocks[index]
		j.l.Unlock()
		if block != nil {
			_, err := j.file.Write(block)
			if err != nil {
				return err
			}
			j.l.Lock()
			delete(j.blocks, index)
			j.l.Unlock()
			index++
		} else {
			time.Sleep(time.Millisecond * 10)
		}
	}

	return j.file.Close()
}

// Name ...
func (j Joiner) Name() string {
	return j.name
}
