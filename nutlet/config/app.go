package config

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)


func LoadAppDefinitionsCfg(file string) (*NutshellCfg, error) {
	if data, err := ioutil.ReadFile(file); err != nil {
		return nil, err
	} else {
		c, err := unmarshalYaml(data)
		return c, err
	}
}

func LoadProcfile(file string)([]string, error) {
	if data, err := ioutil.ReadFile(file); err != nil {
		return make([]string, 0), err
	} else {
		apps := make([]string, 0)
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			fields := strings.Split(line, ":")
			if len(fields) > 0  {
				apps = append(apps, fields[0])
				//fields[0], nil
			} else {
				return nil, errors.New(fmt.Sprintf("valid procfile line:%s", line))
			}
		}
		return apps, err
	}
}

func unmarshalYaml(data []byte) (*NutshellCfg, error) {
	c := &NutshellCfg{}
	err := yaml.Unmarshal(data, &c)
	return c, err
}



func WatchFile(fileName string, action func() error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					err := action()
					if err != nil {
						log.Printf("do action for file:%s change error:%v", fileName, err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()


	err = action()
	if err != nil {
		log.Printf("do action for file:%s change error:%v\n", fileName, err)
	}

	err = watcher.Add(fileName)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<- done
		log.Printf("done watch file:%s\n", fileName)
	}()
}
