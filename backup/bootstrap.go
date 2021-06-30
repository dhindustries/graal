package graal

import (
	"github.com/dhindustries/graal/queue"
)

type Bootstrapper interface {
	Run(app interface{}) error
}

type bootstrapper struct {
	api Api
}

func Bootstrap(providers ...ApiProvider) Bootstrapper {
	var bs bootstrapper
	for _, p := range providers {
		if err := p.Provide(&bs.api); err != nil {
			panic(err)
		}
	}
	return bs
}

func (bs bootstrapper) Run(app interface{}) error {
	ctx := context{
		app: app,
		api: bs.api,
		rnr: queue.NewRunner(&queue.Main),
	}
	return ctx.run()
}
