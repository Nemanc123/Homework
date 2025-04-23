package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	signal := make(Bi)
	temp := runStages(in, stages...)
	go func() {
		defer close(signal)
		for {
			select {
			case <-done:
				return
			case v, ok := <-temp:
				if !ok {
					return
				}
				signal <- v
			}
		}
	}()
	return signal
}
func runStages(in In, stages ...Stage) Out {
	for _, k := range stages {
		in = k(in)
	}
	return in
}
