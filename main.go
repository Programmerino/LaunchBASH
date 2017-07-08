package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"strings"

	"github.com/rakyll/launchpad"
)

var hitX int
var hitY int
var bindings map[string]string
var ch <-chan launchpad.Hit

func start() {
	pad, err := launchpad.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer pad.Close()

	pad.Clear()

	ch = pad.Listen()
	for {
		select {
		case hit := <-ch:
			fmt.Printf("%v, %v\n", hit.X, hit.Y)
			if _, ok := bindings[string(hit.X)+string(hit.Y)]; ok {
				output, err := exec.Command("bash", "-c", bindings[string(hit.X)+string(hit.Y)]).CombinedOutput()
				if err != nil {
					os.Stderr.WriteString(err.Error())
				}
				fmt.Println(string(output))
			}
			pad.Clear()
			pad.Light(hit.X, hit.Y, 3, 3)
		}
	}
}
func main() {
	keypressed := false
	fmt.Println("Learning mode!")
	fmt.Println("Please press the button that you would like to bind")
	pad, err := launchpad.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer pad.Close()

	pad.Clear()

	ch = pad.Listen()
	for !keypressed {
		select {
		case hit := <-ch:
			fmt.Printf("%v, %v\n", hit.X, hit.Y)
			hitX = hit.X
			hitY = hit.Y
			pad.Clear()
			pad.Light(hit.X, hit.Y, 3, 3)
			keypressed = true
		}
	}
	fmt.Println("What command in your default shell would you like to bind this to?")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Printf("You want to bind %s to %v, %v?\n", text, hitX, hitY)
	yn, _ := reader.ReadString('\n')
	yn = strings.ToUpper(yn)
	if strings.TrimRight(yn, "\n") == "N" {
		pad.Close()
		main()
	}
	bindings[string(hitX)+string(hitY)] = text
	fmt.Printf("Do you want to bind more keys?\n")
	yn1, _ := reader.ReadString('\n')
	yn1 = strings.ToUpper(yn1)
	if strings.TrimRight(yn1, "\n") == "Y" {
		pad.Close()
		main()
	} else {
		fmt.Println(yn1)
	}
	fmt.Println("Complete!")
	pad.Close()
	start()
}

func init() {
	bindings = make(map[string]string)
}
