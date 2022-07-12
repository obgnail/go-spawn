package config

import (
	"bufio"
	"encoding/json"
	"github.com/juju/errors"
	"io/ioutil"
	"os"
)

var Config *Conf

type Conf struct {
	Addr             string `json:"addr"`
	User             string `json:"user"`
	SSHDirPath       string `json:"ssh_dir_path"`
	Password         string `json:"password"`
	PrivateKey       string `json:"private_key"`
	Timeout          int    `json:"timeout"`
	CommandChainPath string `json:"command_chain_path"`
	Commands         []string
}

func tryGetPath(path string) string {
	// 向上找5层，满足在一些单元测试中加载不了配置文件的问题
	for i := 0; i < 5; i++ {
		if _, err := os.Stat(path); err == nil {
			break
		} else {
			path = ".." + string(os.PathSeparator) + path
		}
	}
	return path
}

func ReadConfig(path string) (*Conf, error) {
	config := new(Conf)
	path = tryGetPath(path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Trace(err)
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, errors.Trace(err)
	}

	cmdPath := tryGetPath(config.CommandChainPath)
	file, err := os.Open(cmdPath)
	if err != nil {
		return nil, errors.Trace(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var cmds []string
	for scanner.Scan() {
		cmds = append(cmds, scanner.Text())
	}
	config.Commands = cmds
	return config, nil
}

func init() {
	configPath := "config/conf.json"
	var err error
	Config, err = ReadConfig(configPath)
	if err != nil {
		panic(err)
	}
}
