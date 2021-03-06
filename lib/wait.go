package lib

import (
	"github.com/cri-o/cri-o/oci"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/wait"
)

func isStopped(c *ContainerServer, ctr *oci.Container) bool {
	c.runtime.UpdateContainerStatus(ctr)
	cStatus := ctr.State()
	return cStatus.Status == oci.ContainerStateStopped
}

// ContainerWait stops a running container with a grace period (i.e., timeout).
func (c *ContainerServer) ContainerWait(container string) (int32, error) {
	ctr, err := c.LookupContainer(container)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to find container %s", container)
	}

	err = wait.PollImmediateInfinite(1,
		func() (bool, error) {
			return isStopped(c, ctr), nil
		},
	)

	if err != nil {
		return 0, err
	}
	exitCode := ctr.State().ExitCode
	c.ContainerStateToDisk(ctr)
	return exitCode, nil
}
