package rffmpeg

import (
	"github.com/xfrr/goffmpeg/transcoder"
	"io"
	"io/ioutil"
	"sync"
)

type Transcoder struct {
	InputPath  string
	OutputPath string

	Input  io.Reader
	Output io.Writer

	Rate int
}

func (t *Transcoder) SetOutputPath(outputPath string) *Transcoder {
	t.OutputPath = outputPath
	return t
}

func (t *Transcoder) SetInputPath(inputPath string) *Transcoder {
	t.InputPath = inputPath
	return t
}

func (t *Transcoder) SetInput(input io.Reader) *Transcoder {
	t.Input = input
	return t
}

func (t *Transcoder) SetOutput(output io.Writer) *Transcoder {
	t.Output = output
	return t
}

func (t *Transcoder) SetRate(rate int) *Transcoder {
	t.Rate = rate
	return t
}

func (t *Transcoder) To(format string) ([]byte, error) {
	trans := new(transcoder.Transcoder)
	err := trans.InitializeEmptyTranscoder()
	if err != nil {
		return nil, err
	}
	err = trans.SetInputPath(t.InputPath)
	if err != nil {
		return nil, err
	}
	r, err := trans.CreateOutputPipe(format)
	if err != nil {
		return nil, err
	}
	defer r.Close()


	if 0 < t.Rate {
		trans.MediaFile().SetAudioRate(t.Rate)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	var data []byte
	go func() {
		defer wg.Done()
		data, _ = ioutil.ReadAll(r)
	}()

	done := trans.Run(true)

	err = <-done
	if err != nil {
		return nil, err
	}
	wg.Wait()
	return data, err

}
