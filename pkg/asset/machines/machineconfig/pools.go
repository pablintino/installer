package machineconfig

import (
	"fmt"

	"github.com/openshift/api/features"
	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	machineConfigPoolRequiredForUpgradeLabel = "operator.machineconfiguration.openshift.io/required-for-upgrade"
	labelMCOBuiltIn                          = "machineconfiguration.openshift.io/mco-built-in"
	labelPoolPrefix                          = "pools.operator.machineconfiguration.openshift.io"
	machineConfigRoleLabel                   = "machineconfiguration.openshift.io/role"
	PoolArbiter                              = "arbiter"
	PoolMaster                               = "master"
	PoolWorker                               = "worker"
)

func GenerateMachineConfigPool(installConfig *types.InstallConfig, pool string) (*mcfgv1.MachineConfigPool, error) {
	osImageStream := mcfgv1.OSImageStreamReference{}
	if installConfig.EnabledFeatureGates().Enabled(features.FeatureGateOSStreams) && installConfig.OSImageStream != "" {
		if installConfig.OSImageStream != "" {
			osImageStream.Name = string(installConfig.OSImageStream)
		}
	}

	var generatedPool *mcfgv1.MachineConfigPool
	switch pool {
	case PoolMaster, PoolArbiter:
		generatedPool = createMachineConfigPool(pool, true, osImageStream)
	case PoolWorker:
		generatedPool = createMachineConfigPool(pool, false, osImageStream)
	default:
		return nil, fmt.Errorf("unsupported MachineConfigPool %q", pool)

	}
	return generatedPool, nil
}

// createMachineConfigPool creates a MachineConfigPool with the given name and configuration.
// requiredForUpgrade indicates whether the pool is required for cluster upgrades (master and arbiter pools).
func createMachineConfigPool(poolName string, requiredForUpgrade bool, osImageStream mcfgv1.OSImageStreamReference) *mcfgv1.MachineConfigPool {
	labels := map[string]string{
		labelMCOBuiltIn: "",
		fmt.Sprintf("%s/%s", labelPoolPrefix, poolName): "",
	}
	if requiredForUpgrade {
		labels[machineConfigPoolRequiredForUpgradeLabel] = ""
	}

	return &mcfgv1.MachineConfigPool{
		TypeMeta: metav1.TypeMeta{
			APIVersion: mcfgv1.SchemeGroupVersion.String(),
			Kind:       "MachineConfigPool",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   poolName,
			Labels: labels,
		},
		Spec: mcfgv1.MachineConfigPoolSpec{
			MachineConfigSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					machineConfigRoleLabel: poolName,
				},
			},
			NodeSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					fmt.Sprintf("node-role.kubernetes.io/%s", poolName): "",
				},
			},
			OSImageStream: osImageStream,
		},
	}
}
