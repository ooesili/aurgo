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

func NewFilteringVisitor(visitor Visitor, seen []string) FilteringVisitor {
	seenMap := make(map[string]bool)
	for _, pkg := range seen {
		seenMap[pkg] = true
	}

	return FilteringVisitor{
		visitor: visitor,
		seen:    seenMap,
	}
}

type FilteringVisitor struct {
	visitor Visitor
	seen    map[string]bool
}

func (f FilteringVisitor) Visit(pkg string) ([]string, error) {
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
