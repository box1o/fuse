package capabilities

import (
	"context"
	"runtime"

	"fuse/pkg/computeapi"
)

type LinuxDetector struct{}

func NewLinuxDetector() *LinuxDetector { return &LinuxDetector{} }

func (d *LinuxDetector) Detect(ctx context.Context) (computeapi.Capabilities, error) {
	if runtime.GOOS != "linux" {
		return computeapi.Capabilities{}, ErrUnsupportedPlatform
	}
	memory, err := detectMemory()
	if err != nil {
		return computeapi.Capabilities{}, err
	}
	storage, err := detectStorage()
	if err != nil {
		return computeapi.Capabilities{}, err
	}
	runtimeInfo := detectDocker(ctx)
	return computeapi.Capabilities{
		SchemaVersion: computeapi.CapabilitySchemaVersion,
		OS:            detectOS(ctx), CPU: detectCPU(), Memory: memory, Storage: storage,
		ContainerRuntime: runtimeInfo, Accelerators: detectNVIDIA(ctx, runtimeInfo.Available),
	}, nil
}
