package aurgo

func NewRepoVisitor(repo Repo) VisitorAdapater {
	return VisitorAdapater{
		repo: repo,
	}
}

type VisitorAdapater struct {
	repo Repo
}

func (v VisitorAdapater) Visit(pkg string) ([]string, error) {
	err := v.repo.Sync(pkg)
	if err != nil {
		return nil, err
	}

	return v.repo.GetDeps(pkg)
}
