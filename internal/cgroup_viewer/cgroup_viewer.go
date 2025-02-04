package cgroup_viewer

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type cgroupDir struct {
	controllers      string //cgroup.controllers cpuset cpu io memory hugetlb pids rdma misc
	cpusetEffective  string // cat cpuset.cpus.effective 0-127
	cpusetExclussive string //cpuset.cpus.exclusive
	cpusetPartition  string //cpuset.cpus.partition  member
	cpuWeight        string //cpu.weight
	cpuWeightNice    string //cpu.weight.nice
	memoryCurrent    string //memory.current

}

func readDir(dir string) (error, cgroupDir) {
	result := cgroupDir{}
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
	_, result := readDir(cgroupK8sRoot)
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
				_, result := readDir(cgroupK8sRoot + "/" + file.Name())
				fmt.Printf("\t%d %s\n", count, file.Name())
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
							_, result := readDir(cgroupK8sRoot + "/" + file.Name() + "/" + container.Name())
							fmt.Printf("\t\t%s\n", container.Name())
							fmt.Printf("\t\ttcpusetPartition  %s", result.cpusetPartition)
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
	_, result = readDir(cgroupK8sRoot)
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
				_, result := readDir(cgroupK8sRoot + "/" + file.Name())
				fmt.Printf("\t%d %s\n", count, file.Name())
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
							_, result := readDir(cgroupK8sRoot + "/" + file.Name() + "/" + container.Name())
							fmt.Printf("\t\t%s\n", container.Name())
							fmt.Printf("\t\ttcpusetPartition  %s", result.cpusetPartition)
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
