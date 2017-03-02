package aurgo

type Repo interface {
	Sync(pkg string) error
	GetDeps(pkg string) ([]string, error)
	List() ([]string, error)
	Remove(pkg string) error
}
