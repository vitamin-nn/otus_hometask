package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	I   = interface{}
	In  = <-chan I
	Out = In
	Bi  = chan I
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var out Out

	for _, stage := range stages {
		out = stage(proxyStage(in, done))
		in = out
	}

	return out
}

func proxyStage(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			case <-done:
				return
			}
		}
	}()
	return out
}
