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
	"sort"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

var token string

// init takes a flag (-token) or an environment variable (LIFXTOKEN) to set the
// LIFX OAuth Access Token needed for HTTP API access.
func init() {
	flag.StringVar(&token, "token", "", "LIFX OAuth Access Token")
	flag.Parse()

	// Check for LIFX OAuth Access Token
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
			log.Println("Unimplemented")
			return
		}

		if strings.HasPrefix(flag.Args()[0], "kel") {
			log.Println("Adjusting kelvin")
			log.Println("Unimplemented")
			return
		}

		if strings.HasPrefix(flag.Args()[0], "sat") {
			log.Println("Adjusting saturation")
			log.Println("Unimplemented")
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
		return errors.New("multiple errors detected")
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

// listBulbs list all lights
func listBulbs() error {
	// https://api.developer.lifx.com/docs/list-lights
	return bulbStatus("all")
}

// bulbStatus individual bulb status (unless bulbID == "all")
func bulbStatus(bulbID string) error {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.lifx.com/v1/lights/%s", bulbID)
	req, err := http.NewRequest("GET", url, nil)
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
	log.Printf("lights: %d", len(lights))
	// sort by group name
	sort.Slice(lights, func(i, j int) bool {
		return lights[i].Group.Name < lights[j].Group.Name
	})
	data := [][]string{}
	for k, v := range lights {
		data = append(data, []string{
			fmt.Sprintf("%d", k),
			v.ID,
			v.Label,
			v.Group.Name,
			v.Power,
			fmt.Sprintf("%.2f", v.Brightness),
			fmt.Sprintf("%.2f", v.Color.Hue),
			fmt.Sprintf("%.0f", v.Color.Kelvin),
			fmt.Sprintf("%.1f", v.Color.Saturation),
			//v.LastSeen.Format(time.RFC1123),
		})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"idx", "ID", "Label", "Group", "Power", "Brightness", "Hue", "Kelvin", "Sat"})
	table.SetBorder(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(data)
	table.Render()

	return nil
}
