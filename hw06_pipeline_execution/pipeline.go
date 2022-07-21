package hw06pipelineexecution

import "log"

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	outCh := in
	for _, stage := range stages {
		outCh = func(in In) (out Out) {
			bindCh := make(Bi)

			go func() {
				defer close(bindCh)

				for {
					select {
					case <-done:
						log.Println("graceful shutdown")
						return
					case v, ok := <-in:
						if !ok {
							return
						}
						select {
						case <-done:
						case bindCh <- v:
							log.Println("move the value to the link channel")
						}
					}
				}
			}()

			log.Println("calling the transformer on the link channel")

			return stage(bindCh)
		}(outCh)
	}

	return outCh
}
