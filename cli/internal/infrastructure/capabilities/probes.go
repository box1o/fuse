package capabilities

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"fuse/pkg/computeapi"
)

var ErrUnsupportedPlatform = errors.New("compute capability detection currently supports Linux only")

func detectOS(ctx context.Context) computeapi.OperatingSystem {
	values := readKeyValueFile("/etc/os-release", "=")
	kernel, _ := exec.CommandContext(ctx, "uname", "-r").Output()
	name := strings.Trim(values["ID"], `"`)
	if name == "" {
		name = runtime.GOOS
	}
	return computeapi.OperatingSystem{Name: name, Version: strings.Trim(values["VERSION_ID"], `"`), Kernel: strings.TrimSpace(string(kernel)), Architecture: runtime.GOARCH}
}

func detectCPU() computeapi.CPU {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return computeapi.CPU{Model: "unknown", PhysicalCores: runtime.NumCPU(), LogicalCores: runtime.NumCPU()}
	}
	defer file.Close()
	var vendor, model, physicalID, coreID string
	cores := map[string]struct{}{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if physicalID != "" || coreID != "" {
				cores[physicalID+":"+coreID] = struct{}{}
			}
			physicalID, coreID = "", ""
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "vendor_id":
			if vendor == "" {
				vendor = value
			}
		case "model name":
			if model == "" {
				model = value
			}
		case "physical id":
			physicalID = value
		case "core id":
			coreID = value
		}
	}
	physical := len(cores)
	if physical == 0 {
		physical = runtime.NumCPU()
	}
	if model == "" {
		model = "unknown"
	}
	return computeapi.CPU{Vendor: vendor, Model: model, PhysicalCores: physical, LogicalCores: runtime.NumCPU()}
}

func detectMemory() (computeapi.Memory, error) {
	fields := strings.Fields(readKeyValueFile("/proc/meminfo", ":")["MemTotal"])
	if len(fields) == 0 {
		return computeapi.Memory{}, fmt.Errorf("cannot detect total memory")
	}
	kib, err := strconv.ParseUint(fields[0], 10, 64)
	if err != nil {
		return computeapi.Memory{}, err
	}
	return computeapi.Memory{TotalBytes: kib * 1024}, nil
}

func detectStorage() (computeapi.Storage, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs("/", &stat); err != nil {
		return computeapi.Storage{}, err
	}
	return computeapi.Storage{TotalBytes: stat.Blocks * uint64(stat.Bsize)}, nil
}

func detectDocker(ctx context.Context) computeapi.ContainerRuntime {
	output, err := exec.CommandContext(ctx, "docker", "--version").Output()
	if err != nil {
		return computeapi.ContainerRuntime{Name: "docker", Available: false}
	}
	version := strings.TrimSpace(strings.TrimPrefix(string(output), "Docker version "))
	if index := strings.Index(version, ","); index >= 0 {
		version = version[:index]
	}
	return computeapi.ContainerRuntime{Name: "docker", Version: version, Available: true}
}

func detectNVIDIA(ctx context.Context, dockerAvailable bool) []computeapi.Accelerator {
	output, err := exec.CommandContext(ctx, "nvidia-smi", "--query-gpu=uuid,name,memory.total,driver_version", "--format=csv,noheader,nounits").Output()
	if err != nil {
		return []computeapi.Accelerator{}
	}
	runtimeAvailable := false
	if dockerAvailable {
		info, infoErr := exec.CommandContext(ctx, "docker", "info", "--format", "{{json .Runtimes}}").Output()
		runtimeAvailable = infoErr == nil && strings.Contains(strings.ToLower(string(info)), "nvidia")
	}
	result := make([]computeapi.Accelerator, 0)
	for _, line := range strings.Split(strings.TrimSpace(string(output)), "\n") {
		parts := strings.Split(line, ",")
		if len(parts) < 4 {
			continue
		}
		memoryMiB, _ := strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 64)
		result = append(result, computeapi.Accelerator{Kind: "gpu", Vendor: "nvidia", DeviceID: strings.TrimSpace(parts[0]), Model: strings.TrimSpace(parts[1]), MemoryBytes: memoryMiB * 1024 * 1024, DriverVersion: strings.TrimSpace(parts[3]), RuntimeAvailable: runtimeAvailable})
	}
	return result
}

func readKeyValueFile(path, separator string) map[string]string {
	values := map[string]string{}
	file, err := os.Open(path)
	if err != nil {
		return values
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), separator, 2)
		if len(parts) == 2 {
			values[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return values
}
