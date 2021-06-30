package graal

type Prefab interface {
	Resource
	Spawn() (Handle, error)
}

type PrefabLoader = func(api *Api, prefab Prefab) (Prefab, error)
