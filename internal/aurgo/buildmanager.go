package aurgo

type BuildEnv interface {
	Create(location string) error
	Exists(location string) (bool, error)
}

func NewBuildManager(buildEnv BuildEnv, config Config) BuildManager {
	return BuildManager{
		config:   config,
		buildEnv: buildEnv,
	}
}

type BuildManager struct {
	config   Config
	buildEnv BuildEnv
}

func (b BuildManager) Provision() error {
	chrootPath := b.config.ChrootPath()

	exists, err := b.buildEnv.Exists(chrootPath)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	err = b.buildEnv.Create(chrootPath)
	if err != nil {
		return err
	}

	return nil
}
