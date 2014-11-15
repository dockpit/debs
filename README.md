# dockpit/debs
A library that wraps the go tool version control download logic. It allows for the downloading of version controlled packages using go-like import path syntax:

```Go

	dir, err := ioutil.TempDir("", "deps")
	if err != nil {
		t.Fatal(err)
	}

	//create manager; provide directory in which to
	//install packages, it will be created if it 
	//it doesn't exist 
	m := tool.NewManager(dir)

	//install using go-like imporpath
	err = m.Install("github.com/golang/example")
	if err != nil {
		t.Fatal(err)
	}

```