package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	I   = interface{}
	In  = <-chan I
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var prevOut = in
	for _, stage := range stages {
		if stage != nil {
			prevOut = createPipe(stage, prevOut, done)
		}
	}

	return prevOut
}

func createPipe(stage Stage, in In, done In) Out {
	var out = make(chan I, cap(in))

	go func() {
		defer close(out)

		if done == nil {
			for item := range stage(in) {
				out <- item
			}
		} else {
			for {
				select {
				case item, ok := <-stage(in):
					if !ok {
						return
					}
					out <- item
				case <-done:
					return
				}
			}
		}
	}()

	return out
}
