package main

import (
	"errors"
	"flag"
	"fmt"
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

func (dataset Dataset) DestroySnapshot(name string) {
	fmt.Printf("Destroying %s \n", name)
	err := exec.Command("zfs", "destroy", dataset.Name+"@"+name).Run()
	crash(err)
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

func (dataset Dataset) listSnapshots() {

	output, err := exec.Command("zfs", "list", "-t", "snap", "-r", "-d", "1", "-H", dataset.Name).Output()
	crash(err)

	lines := strings.Split(string(output), "\n")
	for i := 0; i < len(lines)-1; i++ {
		lineSplit := strings.Split(lines[i], "\t")
		name := strings.Replace(lineSplit[0], dataset.Name+"@", "", 1)
		fmt.Printf("%s\t%s\t%s\n", name, lineSplit[1], lineSplit[3])
	}
}

func (dataset Dataset) DestroyAllSnapshots() {

	output, err := exec.Command("zfs", "list", "-t", "snap", "-r", "-d", "1", "-H", dataset.Name).Output()
	crash(err)

	lines := strings.Split(string(output), "\n")
	for i := 0; i < len(lines)-1; i++ {
		lineSplit := strings.Split(lines[i], "\t")
		name := strings.Replace(lineSplit[0], dataset.Name+"@", "", 1)
		dataset.DestroySnapshot(name)
	}
}

func currentDataset() (Dataset, error) {
	path, err := os.Getwd()
	crash(err)

	//TODO dont do this on linux ?
	if strings.HasPrefix(path, "/home") {
		path = strings.Replace(path, "/home", "/usr/home", 1)
	}

	fmt.Printf("Current path %s\n", path)

	datasets := list()
	var possibleDataset Dataset
	for _, dataset := range datasets {
		if strings.HasPrefix(path, dataset.Mountpoint) && len(dataset.Mountpoint) > len(possibleDataset.Mountpoint) {
			possibleDataset = dataset
		}
	}
	if possibleDataset.Name == "" {
		return Dataset{}, errors.New("Cant find current dataset")
	}

	return possibleDataset, nil
}

func main() {

	all := flag.Bool("a", false, "All the things")

	flag.Parse()

	dataset, err := currentDataset()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current dataset %s\n", dataset.Name)

	switch flag.Arg(0) {
	case "list":
		dataset.listSnapshots()
	case "destroy", "rm":
		name := flag.Arg(1)
		if name != "" {
			dataset.DestroySnapshot(name)
		} else if *all == true {
			dataset.DestroyAllSnapshots()
		} else {
			fmt.Println("Must provide name")
		}
	default:
		flag.PrintDefaults()
	}

}
