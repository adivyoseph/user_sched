package exporter

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
)

type ContainerState struct {
	Id       string `yaml:"id"`
	Limit    int    `yaml:"limit"`
	Requests int    `yaml:"requests"`
}

type PodState struct {
	Name       string                    `yaml:"name"`
	Reserved   bool                      `yam;:"reserved"`
	Containers map[string]ContainerState `yaml:"containers"`
}

type ChipletsState struct {
	Id         int                 `yaml: "id"`
	PowerState string              `yaml:"powerstate"`
	Capacity   int                 `yaml:"capacity"`
	Limit      int                 `yaml:"limit"`
	Requests   int                 `yaml:"requests"`
	Cpuset     string              `yaml:"cpuset"`
	Reserved   string              `yaml:"reserved"`
	Pods       map[string]PodState `yaml:"pods"`
}

type ChipletsYAML struct {
	ApiVersion string          `yaml:"apiVersion"`
	Chiplets   []ChipletsState `yaml:"chiplets"`
}

type exporterState struct {
	chipletCnt        int
	chipletPrefString string
	typeCnt           int
	typeString        []string
}

var chiplets ChipletsYAML

func (s *exporterState) addType(typeName string) {
	s.typeString = append(s.typeString, typeName)
	s.typeCnt++
}

func Exporter_main() {
	state := exporterState{
		chipletPrefString: "chiplet_",
	}
	state.addType("power")
	state.addType("requests")
	state.addType("podCnt")
	state.addType("freq")

	//read state file once to get chiplet count
	readState()
	state.chipletCnt = len(chiplets.Chiplets)

	chipletStats := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ChipletStatus",
			Help: "labels control mapping",
		},
		[]string{"chiplet", "type"},
	)

	prometheus.MustRegister(chipletStats)

	go func() {
		counter := 0
		for {
			readState()
			for c := 0; c < state.chipletCnt; c++ {
				chipletName := fmt.Sprintf("%s%03d", state.chipletPrefString, c)
				t := 0
				powerLevel := 0
				if chiplets.Chiplets[c].PowerState == "low " {
					powerLevel = 10
				}
				if chiplets.Chiplets[c].PowerState == "norm" {
					powerLevel = 50
				}
				if chiplets.Chiplets[c].PowerState == "high" {
					powerLevel = 90
				}
				chipletStats.WithLabelValues(chipletName, state.typeString[t]).Set(float64(powerLevel))
				//fmt.Printf("%s %s %f %s\n", chipletName, state.typeString[t], float64(powerLevel), chiplets.Chiplets[c].PowerState)
				//freq
				//request
				requests := chiplets.Chiplets[c].Requests / 160
				chipletStats.WithLabelValues(chipletName, state.typeString[1]).Set(float64(requests))
				//pod count
				podCnt := len(chiplets.Chiplets[c].Pods) * 5
				chipletStats.WithLabelValues(chipletName, state.typeString[2]).Set(float64(podCnt))
				fmt.Printf("%s %s %03.1f %3d \trequests %05d\t power %s\n",
					chipletName,
					state.typeString[2],
					float64(podCnt),
					len(chiplets.Chiplets[c].Pods),

					chiplets.Chiplets[c].Requests,
					chiplets.Chiplets[c].PowerState)

			}
			fmt.Printf("--- %d\n", counter)
			counter++
			time.Sleep(time.Second * 3)
		}

	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":50001", nil))

}

func readState() {
	f, err := os.ReadFile("/tmp/state.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(f, &chiplets); err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%+v\n", chiplets)

}
