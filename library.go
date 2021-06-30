package graal

import "fmt"

type Library interface {
	Name() string
	Install(api *ApiPrototype) error
	Init(api Api) error
	Finish(api Api)
}

var libraries = []Library{&core{}}

func UseLibrary(library Library) {
	libraries = append(libraries, library)
}

func installLibraries() (Api, error) {
	if api != nil {
		return nil, fmt.Errorf("libraries are already installed")
	}
	proto := ApiPrototype{}
	for _, lib := range libraries {
		fmt.Printf("Installing library %s\n", lib.Name())
		if err := lib.Install(&proto); err != nil {
			return nil, err
		}
	}
	return &apiAdapter{proto, false}, nil
}

func initLibraries(api Api) error {
	return api.TryInvoke(func(api Api) error {
		initialized := []Library{}
		var err error
		for _, lib := range libraries {
			if err = lib.Init(api); err != nil {
				break
			}
			initialized = append([]Library{lib}, initialized...)
		}
		if err != nil {
			for _, lib := range initialized {
				lib.Finish(api)
			}
			return err
		}
		return nil
	})
}

func finishLibraries(api Api) {
	api.Invoke(func(api Api) {
		for i := len(libraries) - 1; i >= 0; i-- {
			libraries[i].Finish(api)
		}
	})
}
