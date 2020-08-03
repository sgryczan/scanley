package scanner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/sgryczan/scanley/models"
)

func ReadScanFromFile(filename string) (*models.ScanResults, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	hs := models.ScanResults{}
	err = json.Unmarshal([]byte(file), &hs)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return &hs, nil
}

//ListScanFiles lists all hosts in inventory
func ListScanFiles() ([]string, error) {
	dir := "scans/"
	files := []string{}

	children, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range children {
		files = append(files, file.Name())
	}
	return files, nil
}

func ExpandPortRanges(portRanges []string) ([]int, error) {
	err := ValidatePortRanges(portRanges)
	if err != nil {
		return nil, err
	}
	result := []int{}
	for _, r := range portRanges {
		s := strings.Split(r, "-")
		start, _ := strconv.Atoi(s[0])
		end, _ := strconv.Atoi(s[1])
		pr := makeRange(start, end)
		result = append(result, pr...)
	}
	return result, nil
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func ValidatePorts(ports []int) error {
	err := checkPortRange(ports)
	if err != nil {
		return err
	}
	return nil
}

func ValidatePortRanges(portRanges []string) error {
	for _, r := range portRanges {
		match, _ := regexp.MatchString(`^\d*-\d*$`, r)
		if !(match) {
			return errors.New(fmt.Sprintf("Invalid port range: %s", r))
		}
	}
	return nil
}

func checkPortRange(ports []int) error {
	for _, p := range ports {
		if p < 1 || p > 65535 {
			return errors.New("Port must be between 1-65535")
		}
	}
	return nil
}

func UniquePorts(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
