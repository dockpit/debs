package debs

type Manager struct {
	Dir string
}

// Manages microservice dependencies, uses go tool approach
// meaning that it retrieves dependencies by using go-like import paths
func NewManager(dir string) *Manager {
	return &Manager{
		Dir: dir,
	}
}

// install the dependency but if it already exists, do nothing
func (m *Manager) Install(ipath string) error {

	//create package
	p := NewPackage(ipath)

	//download it
	return DownloadPackage(p, m.Dir, false)
}

// return the absolate path of a service by its import path
func (m *Manager) Locate(ipath string) (string, error) {

	//create package
	p := NewPackage(ipath)

	//@todo return more info
	_, _, _, root, err := ExpandPackage(p, m.Dir)
	if err != nil {
		return "", err
	}

	return root, nil
}

// update the dependency, if its not installed install
func (m *Manager) Upsert(ipath string) error {

	//create package
	p := NewPackage(ipath)

	//download it with update flag enabled
	return DownloadPackage(p, m.Dir, true)
}
