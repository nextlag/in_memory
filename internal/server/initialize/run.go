package initialize

import (
	"context"
	"os"
	"time"
)

func (i *Initialize) Run(ctx context.Context) {

	sigs := make(chan os.Signal, 1)
	defer close(sigs)

	i.wg.Add(1)
	go func() {
		defer i.wg.Done()
		if err := i.srv.LaunchServer(ctx, func(ctx context.Context, query []byte, count int) (response string) {
			message := string(query[:count])
			response = i.uc.HandleQuery(ctx, message)
			return response
		}); err != nil {
			sigs <- os.Interrupt
		}
	}()

	<-ctx.Done()

	sCtx, sCansel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer sCansel()

	if err := i.srv.Close(sCtx); err != nil {
		i.log.Error("error closing server", "err", sCtx.Err().Error())
	}

	i.wg.Wait()
	i.log.Info("Close completed")
}
