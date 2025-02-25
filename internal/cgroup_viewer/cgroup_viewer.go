package cgroup_viewer

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

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
	Id         int                 `yaml:"id"`
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

var chiplets ChipletsYAML

func initState() {
	f, err := os.ReadFile("/tmp/state.yaml")
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(f, &chiplets); err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%+v\n", chiplets)

}

type cgroupDir struct {
	stateName        string
	reserved         bool
	stateCpuSet      string
	stateRequest     int
	controllers      string //cgroup.controllers cpuset cpu io memory hugetlb pids rdma misc
	cpusetEffective  string // cat cpuset.cpus.effective 0-127
	cpusetExclussive string //cpuset.cpus.exclusive
	cpusetPartition  string //cpuset.cpus.partition  member
	cpuWeight        string //cpu.weight
	cpuWeightNice    string //cpu.weight.nice
	memoryCurrent    string //memory.current

}

func myreadDir(dir string, prefix string, t int) (error, cgroupDir) {
	result := cgroupDir{}
	if t != 1 && t != 4 {
		initState()
		if t == 2 {
			//burstable pod
			temp := strings.TrimPrefix(dir, prefix)
			temp = strings.TrimSuffix(temp, ".slice")
			temp = strings.TrimPrefix(temp, "kubepods-burstable-pod")
			temp = strings.Replace(temp, "_", "", -1)
			//fmt.Printf("myReadir 2 temp %s\n", temp)
			for x := 0; x < len(chiplets.Chiplets); x++ {
				for pod, podstate := range chiplets.Chiplets[x].Pods {
					tempPod := strings.Replace(pod, "-", "", -1)
					if tempPod == temp {
						result.stateName = podstate.Name
						result.reserved = podstate.Reserved
						if result.reserved {
							result.stateCpuSet = chiplets.Chiplets[x].Reserved
						} else {
							result.stateCpuSet = chiplets.Chiplets[x].Cpuset
						}
						break
					}
				}

			}
		}
		if t == 3 {
			//burstable container
			temp := strings.TrimPrefix(dir, prefix)
			//fmt.Printf("myReadir 3 temp 1 %s\n", temp)
			temp = strings.TrimPrefix(temp, "cri-containerd-")
			//fmt.Printf("myReadir 3 temp 2 %s\n", temp)
			temp = strings.TrimSuffix(temp, ".scope")
			temp = strings.Replace(temp, "_", "", -1)
			//fmt.Printf("myReadir 3 temp %s\n", temp)
			for x := 0; x < len(chiplets.Chiplets); x++ {
				for pod, podstate := range chiplets.Chiplets[x].Pods {
					for con, containerState := range chiplets.Chiplets[x].Pods[pod].Containers {
						//fmt.Printf("myReadir 3 containerState.Id %s\n", containerState.Id)
						if temp == containerState.Id {
							result.stateName = con
							if podstate.Reserved {
								result.stateCpuSet = chiplets.Chiplets[x].Reserved
							} else {
								result.stateCpuSet = chiplets.Chiplets[x].Cpuset
							}
							result.stateRequest = containerState.Requests
							break
						}
					}
				}

			}
		}
		if t == 5 {
			//burstable pod
			temp := strings.TrimPrefix(dir, prefix)
			temp = strings.TrimPrefix(temp, "kubepods-besteffort-pod")
			temp = strings.TrimSuffix(temp, ".slice")
			temp = strings.Replace(temp, "_", "", -1)
			for x := 0; x < len(chiplets.Chiplets); x++ {
				for pod, podstate := range chiplets.Chiplets[x].Pods {
					tempPod := strings.Replace(pod, "-", "", -1)
					if tempPod == temp {
						result.stateName = podstate.Name
						result.reserved = podstate.Reserved
						if result.reserved {
							result.stateCpuSet = chiplets.Chiplets[x].Reserved

						} else {
							result.stateCpuSet = chiplets.Chiplets[x].Cpuset
						}
						break
					}
				}

			}
		}
		if t == 6 {
			//burstable container
			temp := strings.TrimPrefix(dir, prefix)
			//fmt.Printf("myReadir 3 temp 1 %s\n", temp)
			temp = strings.TrimPrefix(temp, "cri-containerd-")
			//fmt.Printf("myReadir 3 temp 2 %s\n", temp)
			temp = strings.TrimSuffix(temp, ".scope")
			temp = strings.Replace(temp, "_", "", -1)
			//fmt.Printf("myReadir 3 temp %s\n", temp)
			for x := 0; x < len(chiplets.Chiplets); x++ {
				for pod, podstate := range chiplets.Chiplets[x].Pods {
					for con, containerState := range chiplets.Chiplets[x].Pods[pod].Containers {
						//fmt.Printf("myReadir 3 containerState.Id %s\n", containerState.Id)
						if temp == containerState.Id {
							result.stateName = con
							if podstate.Reserved {
								result.stateCpuSet = chiplets.Chiplets[x].Reserved
							} else {
								result.stateCpuSet = chiplets.Chiplets[x].Cpuset
							}
							result.stateRequest = containerState.Requests
							break
						}
					}
				}

			}
		}
	}

	file := fmt.Sprintf("%s/cpuset.cpus.partition", dir)
	content, err := os.ReadFile(file)
	result.cpusetPartition = string(content)

	file = fmt.Sprintf("%s/cpuset.cpus.effective", dir)
	content, err = os.ReadFile(file)
	result.cpusetEffective = string(content)

	file = fmt.Sprintf("%s/cpuset.cpus.exclusive", dir)
	content, err = os.ReadFile(file)
	result.cpusetExclussive = string(content)

	file = fmt.Sprintf("%s/cpu.weight", dir)
	content, err = os.ReadFile(file)
	result.cpuWeight = string(content)

	file = fmt.Sprintf("%s/cpu.weight.nice", dir)
	content, err = os.ReadFile(file)
	result.cpuWeightNice = string(content)

	return err, result
}

func Viewer_main() {
	cgroupK8sRoot := fmt.Sprintf("/sys/fs/cgroup/kubepods.slice/kubepods-burstable.slice")

	if _, err := os.Stat(cgroupK8sRoot); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does *not* exist
			fmt.Printf("dir %s does not exit\n", cgroupK8sRoot)
			os.Exit(1)
		} else {
			// Schrodinger: file may or may not exist. See err for details.
			// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
			fmt.Printf("dir %s othe error %v\n", cgroupK8sRoot, err)
			os.Exit(1)
		}

	}
	// path/to/whatever exists
	fmt.Printf("dir %s found\n", cgroupK8sRoot)
	_, result := myreadDir(cgroupK8sRoot, "", 1)
	fmt.Printf("result:\n")
	fmt.Printf("   cpusetPartition  %s", result.cpusetPartition)
	fmt.Printf("   cpusetEffective  %s", result.cpusetEffective)
	fmt.Printf("   cpusetExclussive %s", result.cpusetExclussive)
	fmt.Printf("   cpuWeight        %s", result.cpuWeight)

	rootFiles, err := os.ReadDir(cgroupK8sRoot)
	if err != nil {
		os.Exit(1)
	}
	count := 0
	for _, file := range rootFiles {
		if file.IsDir() {
			_, _, found := strings.Cut(file.Name(), "kubepods-burstable")
			if found {
				_, result := myreadDir(cgroupK8sRoot+"/"+file.Name(), cgroupK8sRoot+"/", 2)
				fmt.Printf("\t%d %s\n", count, file.Name())
				fmt.Printf("\tPodName          %s\n", result.stateName)
				fmt.Printf("\tReserved         %t\n", result.reserved)
				fmt.Printf("\tcpuSet           %s\n", result.stateCpuSet)
				fmt.Printf("\tcpusetPartition  %s", result.cpusetPartition)
				fmt.Printf("\tcpusetEffective  %s", result.cpusetEffective)
				fmt.Printf("\tcpusetExclussive %s", result.cpusetExclussive)
				fmt.Printf("\tcpuWeight        %s", result.cpuWeight)
				containers, errs := os.ReadDir(cgroupK8sRoot + "/" + file.Name())
				if errs != nil {
					os.Exit(1)
				}
				for _, container := range containers {
					if container.IsDir() {
						_, _, foundCon := strings.Cut(container.Name(), "cri-container")
						if foundCon {
							_, result := myreadDir(cgroupK8sRoot+"/"+file.Name()+"/"+container.Name(), cgroupK8sRoot+"/"+file.Name()+"/", 3)
							fmt.Printf("\t\t%s\n", container.Name())
							fmt.Printf("\t\tContainerName    %s\n", result.stateName)
							fmt.Printf("\t\tcpuSet           %s\n", result.stateCpuSet)
							fmt.Printf("\t\trequests         %d\n", result.stateRequest)
							fmt.Printf("\t\ttcpusetPartition %s", result.cpusetPartition)
							fmt.Printf("\t\tcpusetEffective  %s", result.cpusetEffective)
							fmt.Printf("\t\tcpusetExclussive %s", result.cpusetExclussive)
							fmt.Printf("\t\tcpuWeight        %s", result.cpuWeight)

						}
					}

				}

				count++

			}
		}
	}
	cgroupK8sRoot = fmt.Sprintf("/sys/fs/cgroup/kubepods.slice/kubepods-besteffort.slice")
	if _, err := os.Stat(cgroupK8sRoot); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// path/to/whatever does *not* exist
			fmt.Printf("dir %s does not exit\n", cgroupK8sRoot)
			os.Exit(1)
		} else {
			// Schrodinger: file may or may not exist. See err for details.
			// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
			fmt.Printf("dir %s othe error %v\n", cgroupK8sRoot, err)
			os.Exit(1)
		}

	}
	// path/to/whatever exists
	fmt.Printf("dir %s found\n", cgroupK8sRoot)
	_, result = myreadDir(cgroupK8sRoot, "", 4)
	fmt.Printf("result:\n")
	fmt.Printf("   cpusetPartition  %s", result.cpusetPartition)
	fmt.Printf("   cpusetEffective  %s", result.cpusetEffective)
	fmt.Printf("   cpusetExclussive %s", result.cpusetExclussive)
	fmt.Printf("   cpuWeight        %s", result.cpuWeight)
	rootFiles, err = os.ReadDir(cgroupK8sRoot)
	if err != nil {
		os.Exit(1)
	}

	for _, file := range rootFiles {
		if file.IsDir() {
			_, _, found := strings.Cut(file.Name(), "kubepods-best")
			if found {
				_, result := myreadDir(cgroupK8sRoot+"/"+file.Name(), cgroupK8sRoot+"/", 5)
				fmt.Printf("\t%d %s\n", count, file.Name())
				fmt.Printf("\tPodName          %s\n", result.stateName)
				fmt.Printf("\tReserved         %t\n", result.reserved)
				fmt.Printf("\tcpusetPartition  %s", result.cpusetPartition)
				fmt.Printf("\tcpusetEffective  %s", result.cpusetEffective)
				fmt.Printf("\tcpusetExclussive %s", result.cpusetExclussive)
				fmt.Printf("\tcpuWeight        %s", result.cpuWeight)
				containers, errs := os.ReadDir(cgroupK8sRoot + "/" + file.Name())
				if errs != nil {
					os.Exit(1)
				}
				for _, container := range containers {
					if container.IsDir() {
						_, _, foundCon := strings.Cut(container.Name(), "cri-container")
						if foundCon {
							_, result := myreadDir(cgroupK8sRoot+"/"+file.Name()+"/"+container.Name(), cgroupK8sRoot+"/"+file.Name()+"/", 6)
							fmt.Printf("\t\t%s\n", container.Name())
							fmt.Printf("\t\tContainerName    %s\n", result.stateName)
							fmt.Printf("\t\tcpuSet           %s\n", result.stateCpuSet)
							fmt.Printf("\t\trequests         %d\n", result.stateRequest)
							fmt.Printf("\t\ttcpusetPartition %s", result.cpusetPartition)
							fmt.Printf("\t\tcpusetEffective  %s", result.cpusetEffective)
							fmt.Printf("\t\tcpusetExclussive %s", result.cpusetExclussive)
							fmt.Printf("\t\tcpuWeight        %s", result.cpuWeight)

						}
					}

				}

				count++

			}
		}
	}

}
