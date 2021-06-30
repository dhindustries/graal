package graal

import "fmt"

func Run(app interface{}) error {
	if api == nil {
		if err := Init(); err != nil {
			return err
		}
		defer Finish()
	}
	if err := api.initWindow(); err != nil {
		return err
	}
	defer api.finishWindow()
	if err := api.initRenderer(); err != nil {
		return err
	}
	defer api.finishRenderer()
	return start(app)
}

func Init() error {
	if api != nil {
		return fmt.Errorf("engine is already initialized")
	}
	var localApi Api
	var err error
	LockMainThread()
	if localApi, err = installLibraries(); err != nil {
		UnlockMainThread()
		return err
	}
	if err := initLibraries(localApi); err != nil {
		UnlockMainThread()
		return err
	}
	api = localApi
	return nil
}

func Finish() {
	if api != nil {
		finishLibraries(api)
		UnlockMainThread()
		api = nil
	}
}
