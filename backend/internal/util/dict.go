package util

import (
	"bufio"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

func LoadSubDomainDict(path, dictType string) (dict []string, err error) {
	var filepath string
	switch dictType {
	case "1":
		filepath = path + "/subdomains-5000.txt"
	case "2":
		filepath = path + "/subdomains-20000.txt"
	case "3":
		filepath = path + "/subdomains-110000.txt"
	default:
		filepath = path + "/subdomains-5000.txt"
	}

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := bufio.NewScanner(file)
	for buf.Scan() {
		if buf.Text() != "" { // 去除空行
			dict = append(dict, buf.Text())
		}
	}
	return dict, nil
}

func LoadCDNList(path string) (data map[string][]string, err error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	yamlData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	data = make(map[string][]string)
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
