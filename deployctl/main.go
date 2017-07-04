package main

func main() {

	dir, err := createArtifactsDir("/tmp")
	if err != nil {
		panic(err.Error())
	}

	setLogFile(dir)

}
