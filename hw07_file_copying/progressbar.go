package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

var (
	progressLength float64 = 50
	percent        float64 = 100
)

type ProgressBar struct {
	io.Reader
	out        io.Writer
	bytesCount int64
	current    int64
}

func NewProxyReader(r io.Reader, size int64) *ProgressBar {
	p := &ProgressBar{
		Reader:     r,
		out:        os.Stdout,
		bytesCount: size,
		current:    0,
	}
	return p
}

// Read: судя по реализации CopyN - данный метод будет вызываться при копирования каждых 32 Кб
// поэтому для больших файлов имеет смысл отрисовывать прогрессбар реже
// например, на каждое N-ое обращение. Либо асинхронно по тикеру как в рекомендованной библиотеке. */
func (p *ProgressBar) Read(r []byte) (n int, err error) {
	n, err = p.Reader.Read(r)
	p.current += int64(n)
	p.write()

	return n, err
}

func (p *ProgressBar) write() {
	s := p.render()
	fmt.Fprintf(p.out, "\r%s", s)
}

func (p *ProgressBar) render() string {
	current := p.current
	var b strings.Builder
	b.WriteString("[")
	ready := float64(current) / float64(p.bytesCount)
	readyPerc := math.Ceil(ready * percent)
	complLen := math.Ceil(ready * progressLength)
	leftLen := progressLength - complLen
	b.WriteString(strings.Repeat("█", int(complLen)))
	b.WriteString(strings.Repeat(" ", int(leftLen)))
	b.WriteString("] ")
	b.WriteString(fmt.Sprintf("%.0f%%", readyPerc))

	return b.String()
}

func (p *ProgressBar) Finish() {
	fmt.Fprintln(p.out, "\nCompleted!")
}
