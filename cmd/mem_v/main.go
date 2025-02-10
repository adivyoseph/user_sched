package main

import (
	"fmt"

	"os"
	"strings"
)

func myCpuid(op uint32) (eax, ebx, ecx, edx uint32)
func myCpuidex(op, op2 uint32) (eax, ebx, ecx, edx uint32)

func main() {
	//eax, _, _, _ := myCpuid(0)
	//fmt.Printf("Max Std Func %d\n", eax)
	if !isAmd() {
		fmt.Printf("Not AMD device\n")
		os.Exit(1)
	}
	fmt.Printf("AMD device\n")
	fmt.Printf("model: %s\n", brandName())
	_, result := hasMSR()
	fmt.Printf("MSR %s\n", result)
	_, result = hasPQM()
	fmt.Printf("PQM %s\n", result)
	_, result = hasEffFreq()
	fmt.Printf("EffFreq %s\n\n", result)

	fmt.Printf("ecx:  Extended Topology Enumeration\n")
	a, b, c, d := myCpuidex(0xB, 0)
	fmt.Printf("0   ThreadMaskWidth %d\n", int(a&0xf))
	fmt.Printf("0   NumLogProc      %d\n", int(b&0xff))
	fmt.Printf("0   HierarchyLevel  %d\n", int((c>>8)&0xf))
	fmt.Printf("0   x2APIC_ID       %d\n", int(d&0xffff))

	a, b, c, d = myCpuidex(0xB, 1)
	fmt.Printf("1   CoreMaskWidth   %d\n", int(a&0xf))
	fmt.Printf("1   NumLogCores     %d\n", int(b&0xff))
	fmt.Printf("1   HierarchyLevel  %d\n", int((c>>8)&0xf))
	fmt.Printf("1   x2APIC_ID       %d\n", int(d&0xffff))
	fmt.Printf("\nPQM Capabilities \n")
	_, b, _, d = myCpuidex(0xf, 0)
	fmt.Printf("0   Max_RMID        %d\n", int(b))
	fmt.Printf("0   L3CacheMon      %d\n", int(d>>1))
	fmt.Printf("\nPQM  L3 Cache Monitoring Capabilities \n")
	a, b, c, d = myCpuidex(0xf, 1)
	fmt.Printf("1   CounterSize      %d\n", int(a&0xf))
	fmt.Printf("1   ScaleFactor      %d\n", int(b))
	fmt.Printf("1   Max_RMIDr        %d\n", int(c))
	fmt.Printf("1  L3CacheBWMonEvt1  %d\n", int((d>>2)&0x1))
	fmt.Printf("1  L3CacheBWMonEvt0  %d\n", int((d>>1)&0x1))
	fmt.Printf("1  L3CacheOccMon     %d\n", int(d&0x1))
	fmt.Printf("\nPQOS Enforcement (PQE) \n")
	_, _, _, d = myCpuidex(0x10, 0)
	fmt.Printf("0   L3Alloc         %d\n", int((d>>1)&0x1))
	fmt.Printf("\nPQOS L3 Cache Allocation Enforcement Capabilities \n")
	a, b, c, d = myCpuidex(0x10, 1)
	fmt.Printf("1   CBM_LEN          %d\n", int(a&0xf))
	fmt.Printf("1   L3ShareAllocMask %0X\n", int(b))
	fmt.Printf("0   CDP              %d\n", int((c>>2)&0x1))
	fmt.Printf("0   COS_MAX          %d\n", int(d&0xff))
	fmt.Printf("\nCache Topology Information\n")
	var i uint32
	i = 0
	exit := 1
	for exit > 0 {
		a, b, c, d = myCpuidex(0x8000001d, i)

		if (a & 0xf) == 0 {
			exit = 0
		} else {
			fmt.Printf("%d level %d type %d\n", i, int((a>>5)&0x7), int(a&0x3))
			fmt.Printf("\t CacheNumWays %d CacheNumSets %d CacheInclusive %d\n", int(b>>22), int(c), int((d>>1)&0x1))
			i++
		}
	}
	/*	for i = 1; i <= 16; i++ {
			_, b, c, _ = myCpuidex(0x8000001e, i)
			fmt.Printf("%02d ComputeUnitId %d threads %d\n", i, int(b&0xff), int((b>>8)&0xff)+1)
			fmt.Printf("\t NodeId %d NodesPerProcessor %d\n", int(c&0xff), int((c>>8)&0xff)+1)
		}
	*/

}

func isAmd() bool {
	_, b, c, d := myCpuid(0)
	v := string(valAsString(b, d, c))
	//fmt.Printf("Vendor %s\n", v)
	return "AuthenticAMD" == v
}

func hasMSR() (bool, string) {
	_, _, _, edx := myCpuid(1)
	enabled := int((edx >> 5) & 0x1)
	if enabled == 1 {
		return true, "enabled"
	}
	return false, "disabled"
}

func hasPQM() (bool, string) {
	_, ebx, _, _ := myCpuid(7)
	enabled := int((ebx >> 12) & 0x1)
	if enabled == 1 {
		return true, "enabled"
	}
	return false, "disabled"
}

func hasEffFreq() (bool, string) {
	_, _, _, edx := myCpuid(6)
	enabled := int(edx & 0x1)
	if enabled == 1 {
		return true, "enabled"
	}
	return false, "disabled"

}

func maxExtendedFunction() uint32 {
	eax, _, _, _ := myCpuid(0x80000000)
	return eax
}

func valAsString(values ...uint32) []byte {
	r := make([]byte, 4*len(values))
	for i, v := range values {
		dst := r[i*4:]
		dst[0] = byte(v & 0xff)
		dst[1] = byte((v >> 8) & 0xff)
		dst[2] = byte((v >> 16) & 0xff)
		dst[3] = byte((v >> 24) & 0xff)
		switch {
		case dst[0] == 0:
			return r[:i*4]
		case dst[1] == 0:
			return r[:i*4+1]
		case dst[2] == 0:
			return r[:i*4+2]
		case dst[3] == 0:
			return r[:i*4+3]
		}
	}
	return r
}

func brandName() string {
	if maxExtendedFunction() >= 0x80000004 {
		v := make([]uint32, 0, 48)
		for i := uint32(0); i < 3; i++ {
			a, b, c, d := myCpuid(0x80000002 + i)
			v = append(v, a, b, c, d)
		}
		return strings.Trim(string(valAsString(v...)), " ")
	}
	return "unknown"
}
