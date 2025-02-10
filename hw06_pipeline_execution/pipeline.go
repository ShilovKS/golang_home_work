package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := processStage(in, done)
	for _, stage := range stages {
		if stage == nil {
			continue
		}
		out = processStage(stage(out), done)
	}
	return out
}

func processStage(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			for range in {
				_ = struct{}{} // Пустая операция
			}
		}()
		defer close(out)
		for {
			select {
			case val, ok := <-in:
				if ok {
					out <- val
				} else {
					return
				}
			case <-done:
				return
			}
		}
	}()
	return out
}
