package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

// todo should be read from smc instead

func getTemp() (int, error) {
	out, err := exec.Command("sysctl", "machdep.xcpm.cpu_thermal_level").Output()
	tempOutput := string(out)

	if err != nil {
		return 0, err
	}

	r := regexp.MustCompile(`[0-9].*`)
	t := r.FindString(tempOutput)

	temperature, err := strconv.Atoi(t)

	if err != nil {
		return 0, err
	}

	return temperature, nil
}

type Information struct {
	CpuTemp int
}

func main() {

	temp, err := getTemp()
	if err != nil {
		fmt.Println(err)
	}

	output := Information{
		CpuTemp: temp,
	}

	out, err := json.Marshal(output)
	fmt.Println(string(out))

}
