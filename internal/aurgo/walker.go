package aurgo

type Visitor interface {
	Visit(pkg string) ([]string, error)
}

func NewVisitingDepWalker(visitor Visitor) VisitingDepWalker {
	return VisitingDepWalker{
		visitor: visitor,
	}
}

type VisitingDepWalker struct {
	visitor Visitor
}

func (v VisitingDepWalker) Walk(pkgs []string) ([]string, error) {
	sw := statefulWalker{
		visitor: v.visitor,
	}

	err := sw.walk(pkgs)
	if err != nil {
		return nil, err
	}

	return sw.visited, nil
}

type statefulWalker struct {
	visited []string
	visitor Visitor
}

func (s *statefulWalker) walk(pkgs []string) error {
	for _, pkg := range pkgs {
		if s.didVisit(pkg) {
			continue
		}

		deps, err := s.visitor.Visit(pkg)
		if err != nil {
			return err
		}
		s.visited = append(s.visited, pkg)

		err = s.walk(deps)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *statefulWalker) didVisit(pkg string) bool {
	for _, visitedPkg := range s.visited {
		if pkg == visitedPkg {
			return true
		}
	}

	return false
}

type PkgManager interface {
	ListAvailable() ([]string, error)
}

func NewFilteringVisitor(visitor Visitor, pkgManager PkgManager) *FilteringVisitor {
	return &FilteringVisitor{
		visitor:    visitor,
		pkgManager: pkgManager,
	}
}

type FilteringVisitor struct {
	visitor    Visitor
	pkgManager PkgManager
	seen       map[string]bool
}

func (f *FilteringVisitor) Visit(pkg string) ([]string, error) {
	err := f.ensureInitialized()
	if err != nil {
		return nil, err
	}

	deps, err := f.visitor.Visit(pkg)
	if err != nil {
		return nil, err
	}

	unseen := []string{}
	for _, pkg := range deps {
		if !f.seen[pkg] {
			unseen = append(unseen, pkg)
		}
	}
	return unseen, nil
}

func (f *FilteringVisitor) ensureInitialized() error {
	if f.seen != nil {
		return nil
	}

	availablePkgs, err := f.pkgManager.ListAvailable()
	if err != nil {
		return err
	}

	f.seen = make(map[string]bool)
	for _, pkg := range availablePkgs {
		f.seen[pkg] = true
	}

	return nil
}
