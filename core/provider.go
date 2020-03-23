package core

import (
	"github.com/dhindustries/graal"
	"github.com/dhindustries/graal/queue"
)

type Provider struct{}

func (Provider) Provide(api *graal.Api) error {
	rm := newResourceManager()
	mem := memory{}
	rt := runtime{queue.NewInvoker(&queue.Main)}
	rm.setLoader(api, "file/*", loadFileResource)
	api.Logf = logf
	api.Schedule = rt.schedule
	api.Invoke = rt.invoke
	api.NewHandle = newHandle
	api.Handle = mem.handle
	api.Cleanup = func(api *graal.Api) {
		rm.cleanup(api)
		mem.cleanup(api)
	}
	// api.GetParent = getParent
	// api.SetParent = setParent
	api.NewResource = rm.newResource
	api.LoadResource = rm.loadResource
	api.SetResourceLoader = rm.setLoader
	api.GetRelativePath = relativePath
	api.LoadFile = loadFile
	return nil
}
