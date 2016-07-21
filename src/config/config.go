package config

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	conf Config
)

const (
	Version = "1.0.0"
)

type Config struct {
	Volume        int
	Cache         int
	CacheSeek     int
	Radios        []RadioCfg
	SelectedRadio int
}

type RadioCfg struct {
	Name string
	Url  string
}

func Init() {

	conf.Volume = 100
	conf.Cache = 4096
	conf.CacheSeek = 16
	conf.Radios = make([]RadioCfg, 0)
	conf.SelectedRadio = 0

	Load()

	if len(conf.Radios) == 0 {
		AddRadio(RadioCfg{Name: "Ultra HD", Url: "http://nashe2.hostingradio.ru/ultra-192.mp3"})
	}
}

func SetCache(val int) {
	conf.Cache = val
	Save()
}

func SetCacheSeek(val int) {
	conf.CacheSeek = val
	Save()
}

func SetSelectedRadio(val int) {
	conf.SelectedRadio = val
	Save()
}

func GetCache() int {
	return conf.Cache
}

func GetCacheSeek() int {
	return conf.CacheSeek
}

func GetSelectedRadio() int {
	return conf.SelectedRadio
}

func GetRadios() []RadioCfg {
	return conf.Radios
}

func AddRadio(r RadioCfg) {
	if conf.Radios == nil {
		conf.Radios = make([]RadioCfg, 1)
	}
	conf.Radios = append(conf.Radios, r)
	Save()
}

func SetRadio(r RadioCfg, pos int) {
	if pos >= 0 && pos < len(conf.Radios) {
		conf.Radios[pos] = r
		Save()
	}
}

func DelRadio(pos int) {
	if pos >= 0 && pos < len(conf.Radios) {
		conf.Radios = append(conf.Radios[:pos], conf.Radios[pos+1:]...)
		Save()
	}
}

func Save() {
	buf, err := xml.MarshalIndent(conf, "", " ")
	if err == nil {
		filenamecfg := filepath.Join(filepath.Dir(os.Args[0]), "sradio.cfg")
		ff, err := os.Create(filenamecfg)
		if err == nil {
			_, err := ff.Write(buf)
			if err != nil {
				log.Println("Error write file", err)
			}
		} else {
			log.Println("Error create cfg file", err)
		}
	} else {
		log.Println("Error make xml", err)
	}
}

func Load() {
	filenamecfg := filepath.Join(filepath.Dir(os.Args[0]), "sradio.cfg")
	buf, err := ioutil.ReadFile(filenamecfg)
	if err == nil {
		err = xml.Unmarshal(buf, &conf)
		if err != nil {
			log.Println("Error unmarshal xml", err)
		}
	} else {
		log.Println("Error read cfg file", err)
	}
}

/*
<Config>
 <Volume>100</Volume>
 <Cache>4096</Cache>
 <CacheSeek>16</CacheSeek>
 <Radios>
  <Name>Ultra HD</Name>
  <Url>http://nashe2.hostingradio.ru/ultra-192.mp3</Url>
 </Radios>
 <Radios>
  <Name>Ultra</Name>
  <Url>http://nashe2.hostingradio.ru/ultra-128.mp3</Url>
 </Radios>
 <Radios>
  <Name>Наше Радио</Name>
  <Url>http://nashe1.hostingradio.ru/nashe-128.mp3</Url>
 </Radios>
 <Radios>
  <Name>Наше 2.0</Name>
  <Url>http://nashe1.hostingradio.ru/nashe20-128.mp3</Url>
 </Radios>
 <Radios>
  <Name>Радио К</Name>
  <Url>http://188.120.253.194:8000/play</Url>
 </Radios>
 <Radios>
  <Name>Nomercy Radio</Name>
  <Url>http://stream6.radiostyle.ru:8006/nomercy</Url>
 </Radios>
 <Radios>
  <Name>Maximum</Name>
  <Url>http://icecast.radiomaximum.cdnvideo.ru/maximum.mp3</Url>
 </Radios>
 <Radios>
  <Name>AvSIM</Name>
  <Url>http://radio.avsim.su:8000/stream</Url>
 </Radios>
 <Radios>
  <Name>RockFM</Name>
  <Url>http://nashe.streamr.ru/rock-128.mp3</Url>
 </Radios>
 <SelectedRadio>5</SelectedRadio>
</Config>
*/
