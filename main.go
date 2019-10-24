package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

func crash(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type Dataset struct {
	Name       string
	Mountpoint string
}

// Only returns mounted datasets
func list() []Dataset {

	output, err := exec.Command("zfs", "list", "-H").Output()
	crash(err)
	lines := strings.Split(string(output), "\n")
	var datasets []Dataset
	for i := 0; i < len(lines)-1; i++ {
		lineSplit := strings.Split(lines[i], "\t")
		mountPoint := lineSplit[4]
		if mountPoint != "/" && mountPoint != "none" && mountPoint != "-" && mountPoint != "/" {
			datasets = append(datasets, Dataset{Name: lineSplit[0], Mountpoint: lineSplit[4]})
		}
	}

	return datasets
}

func currentDataset() (Dataset, error) {
	path, err := os.Getwd()
	crash(err)

	datasets := list()
	for _, dataset := range datasets {
		if strings.HasPrefix(path, dataset.Mountpoint) {
			return dataset, nil
		}
	}
	return Dataset{}, errors.New("Cant find current dataset")
}

func main() {
	dataset, err := currentDataset()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Current dataset %s\t", dataset.Name)

}
