package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

/*
 * @author leig
 * @date 2022/01/10/10:22 PM
 */

var (
	cmd           string
	rootPath      string
	interfacePath = "Interface"
	addOnsPath    = "AddOns"
	accountPath   = "Account"
	wtfPath       = "WTF"
	fontsPath     = "Fonts"
	schema        = 0
	wpmConfig     *WpmConfig
	err           error
)

type Schema int

const (
	All Schema = iota
	Express
)

type WpmConfig struct {
	RootPath      string `json:"root-path"`
	AddOnsPath    string `json:"add-ons-path"`
	InterfacePath string `json:"interface-path"`
	AccountPath   string `json:"account-path"`
	WTFPath       string `json:"wtf-path"`
	FontsPath     string `json:"fonts-path"`
	Schema        Schema `json:"schema"`
}

// backups "recover"

func init() {
	cmd = os.Args[len(os.Args)-1]
	var cfg *ini.File
	if cfg, err = ini.Load("wpm.ini"); err != nil {
		log.Println(err)
	} else {
		rootPath = cfg.Section("").Key("RootPath").String()
		if schema, err = cfg.Section("").Key("Schema").Int(); err != nil {
			log.Println(err)
		}
	}

	flag.StringVar(&rootPath, "r", rootPath, "WOW 安装根路径")
	flag.IntVar(&schema, "s", schema, "备份模式 默认备份全部文件，包括游戏设置；简单模式 -s 1 ")
	flag.Parse()

	if strings.EqualFold(rootPath, "") {
		log.Println("rootPath is nil!")
		os.Exit(1)
	}

	wpmConfig = &WpmConfig{
		RootPath:      rootPath,
		InterfacePath: interfacePath,
		AddOnsPath:    addOnsPath,
		AccountPath:   accountPath,
		WTFPath:       wtfPath,
		FontsPath:     fontsPath,
		Schema:        Schema(schema),
	}
	jsonStr, _ := json.Marshal(wpmConfig)
	log.Println("wpmConfig: " + string(jsonStr))
}

func main() {
	paths := make(map[string]string)
	switch wpmConfig.Schema {
	case All:
		paths[wpmConfig.InterfacePath] = filepath.Join(wpmConfig.RootPath, wpmConfig.InterfacePath)
		paths[wpmConfig.WTFPath] = filepath.Join(wpmConfig.RootPath, wpmConfig.WTFPath)
		paths[wpmConfig.FontsPath] = filepath.Join(wpmConfig.RootPath, wpmConfig.FontsPath)
		break
	case Express:
		paths[wpmConfig.AddOnsPath] = filepath.Join(filepath.Join(wpmConfig.RootPath, wpmConfig.InterfacePath), wpmConfig.AddOnsPath)
		paths[wpmConfig.AccountPath] = filepath.Join(filepath.Join(wpmConfig.RootPath, wpmConfig.WTFPath), wpmConfig.AccountPath)
		paths[wpmConfig.FontsPath] = filepath.Join(wpmConfig.RootPath, wpmConfig.FontsPath)
		break
	default:
		fmt.Println("Schema is wrong!")
	}
	switch cmd {
	case "backups":
		backups(paths)
	case "recover":
		recover(paths, wpmConfig.RootPath)
	default:
		log.Println("请输入操作命令")
	}
	os.Exit(0)
}

func backups(paths map[string]string) {
	var wg sync.WaitGroup
	wg.Add(len(paths))
	for k, v := range paths {
		go Zip(k+".zip", v, wg.Done)
	}
	wg.Wait()
}

func recover(paths map[string]string, rootPath string) {
	s := string(os.PathSeparator)
	rootPath = rootPath[:strings.Index(rootPath, s)+1]
	var wg sync.WaitGroup
	wg.Add(len(paths))
	var f *os.File
	for k := range paths {
		if f, err = os.Open(k + ".zip"); err != nil {
			log.Println(err)
		} else {
			go UnZip(rootPath, f.Name(), wg.Done)
		}
	}
	wg.Wait()
}