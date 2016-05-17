package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var token string

// init takes a flag (-token) or an environment variable (LIFXTOKEN) to set the
// LIFX OAuth Access Token needed for HTTP API access.
func init() {

	tokenstr := flag.String("token", "", "LIFX OAuth Access Token")
	flag.Parse()

	// Check for LIFX OAuth Access Token
	token = string(*tokenstr)
	if token == "" {
		token = os.Getenv("LIFXTOKEN")
	}
	if token == "" {
		log.Println("Please use either a flag or env variable (LIFXTOKEN) to specify LIFX OAuth Access Token.")
		flag.PrintDefaults()
		os.Exit(1)
	}

}

func main() {
	// commands (implemented):
	//  list [bulbid] - list status for all or individual bulb
	//  toggle bulbid [bulbid...] - toggle individual or multiple bulbs
	//  bri bulbid value - adjust bulb brightness to value

	log.Println("HTTP LIFX control")
	//log.Println("LIFX OAuth Access Token", token)

	if len(flag.Args()) > 0 {
		log.Println("Args called:", flag.Args())
		// list: lists status of bulb or bulbs
		// list - list all status
		// list id - list status of one bulb
		if flag.Args()[0] == "list" {
			if len(flag.Args()) == 1 {
				err := listBulbs()
				if err != nil {
					log.Println("Unable to list bulbs.")
					log.Println(err)
				}
				return
			}
			// no else needed here, return above precludes
			err := bulbStatus(flag.Args()[1])
			if err != nil {
				log.Println("Unable to list status for", flag.Args()[1])
				log.Println(err)
			}
			return
		}

		// toggle: bulb on|off
		// toggle id - toggle single bulb
		// toggle id id id - toggle list of bulbs
		if flag.Args()[0] == "toggle" {
			if len(flag.Args()) == 1 {
				log.Println("Must provide a bulb ID or bulb IDs to toggle.")
				return
			} else if len(flag.Args()) == 2 {
				err := toggleBulb(flag.Args()[1])
				if err != nil {
					log.Println(err)
				}
				return
			} else {
				err := toggleBulbs(flag.Args()[1:])
				if err != nil {
					log.Println(err)
				}
				return
			}
		}

		// Adjust state: brightness, hue, kelvin, saturation
		// bri bulbid value - single bulb brightness
		// hue bulbid value - single bulb hue
		// kel bulbid value - single bulb kelvin
		// sat bulbid value - single bulb saturation
		//attributes := []string{"bri", "hue", "kel", "sat"}
		if strings.HasPrefix(flag.Args()[0], "bri") {
			// needs arg length check (0 bri, 1 bulbid, 2 value) - single bulb
			// could also be (0 bri, 1..n bulbid, n+1 value) - array of bulbs
			log.Println("Adjusting brightness")
			err := setBrightness(flag.Args()[1], flag.Args()[2])
			if err != nil {
				log.Println("Can't even", err)
			}

			return
		}

		if strings.HasPrefix(flag.Args()[0], "hue") {
			log.Println("Adjusting hue")
			log.Pritnln("Unimplemented")
			return
		}

		if strings.HasPrefix(flag.Args()[0], "kel") {
			log.Println("Adjusting kelvin")
			log.Pritnln("Unimplemented")
			return
		}

		if strings.HasPrefix(flag.Args()[0], "sat") {
			log.Println("Adjusting saturation")
			log.Pritnln("Unimplemented")
			return
		}
	}

}

func setBrightness(bulbID string, brightness string) error {

	brival, err := strconv.ParseFloat(brightness, 64)
	if err != nil {
		log.Printf("Can't parse %s to float.", brightness)
		return err
	}
	state := State{Brightness: brival, Selector: bulbID}
	states := States{[]State{state}}
	err = setStatus(bulbID, states)
	if err != nil {
		return err
	}
	return nil
}

func setStatus(bulbID string, states States) error {
	client := &http.Client{}
	statebytes, err := json.Marshal(states)
	log.Printf("%s", statebytes)
	if err != nil {
		log.Println("Can't turn bulb states into bytes, aborting call to LIFX.")
		return err
	}
	buf := bytes.NewBuffer(statebytes)
	req, err := http.NewRequest("PUT", "https://api.lifx.com/v1/lights/states", buf)
	req.Header.Set("Authorization", "Bearer "+token)
	if err != nil {
		log.Println("Error creating request for LIFX")
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error executing request to LIFX")
		return err
	}
	log.Println(resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't read body")
		return err
	}
	log.Printf("%s", body)

	return nil
}

func toggleBulbs(bulbs []string) error {
	badness := false
	for _, v := range bulbs {
		err := toggleBulb(v)
		if err != nil {
			badness = true
			log.Println(err)
		}
	}
	if badness {
		return errors.New("Multiple errors detected.")
	}
	return nil
}

func toggleBulb(bulbID string) error {

	client := &http.Client{}
	durationstr := `{"duration":"2"}`
	buf := strings.NewReader(durationstr)
	req, err := http.NewRequest("POST", "https://api.lifx.com/v1/lights/id:"+bulbID+"/toggle", buf)
	if err != nil {
		log.Println("Unable to construct request to toggle bulb", bulbID)
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Unable to call LIFX")
		return err
	}

	log.Println(resp.Status)
	responsebytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't read bytes")
	}
	//if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
	fmt.Printf("%s\n", responsebytes)
	//}

	return nil
}

func listBulbs() error {
	return bulbStatus("all")
}

func bulbStatus(bulbID string) error {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.lifx.com/v1/lights/"+bulbID, nil)
	if err != nil {
		log.Println("Can't construct HTTP request to list lightbulbs.")
		return err
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Can't make an HTTP request to list lightbulbs.")
		return err
	}

	msgbytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Can't convert HTTP response to bytes.")
		return err
	}

	if resp.StatusCode != 200 {
		log.Println("Status", resp.Status)
		fmt.Printf("%s\n", msgbytes)
		return err
	}

	var lights []Lifx
	json.Unmarshal(msgbytes, &lights)
	for k, v := range lights {
		fmt.Printf("%d %s %11s %3s %2.2f (%.2f, %v, %.2f)\n", k, v.ID, v.Label, v.Power, v.Brightness, v.Color.Hue, v.Color.Kelvin, v.Color.Saturation)
	}

	return nil
}
