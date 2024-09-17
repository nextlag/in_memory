package initialize

import (
	"context"
	"os"
)

func (i *Initialize) Run(ctx context.Context) {

	sigs := make(chan os.Signal, 1)
	defer close(sigs)

	i.wg.Add(1)
	go func() {
		defer i.wg.Done()
		if err := i.srv.LaunchServer(ctx, func(ctx context.Context, query []byte) []byte {
			response := i.uc.HandleQuery(ctx, string(query))
			return []byte(response)
		}); err != nil {
			sigs <- os.Interrupt
		}
	}()

	<-ctx.Done()

	i.srv.Close()

	i.wg.Wait()
	i.log.Info("Close completed")
}
