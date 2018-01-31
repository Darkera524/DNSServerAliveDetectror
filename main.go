package main

import "flag"

func main(){
	cfg := flag.String("c", "cfg.json", "configuration file")

	ParseConfig(*cfg)

	go CronConfig(60, *cfg)

	go CronDetect()

	select{}
}
