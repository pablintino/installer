package machineconfig

import (
	"fmt"
	"path/filepath"

	"sigs.k8s.io/yaml"

	mcfgv1 "github.com/openshift/api/machineconfiguration/v1"
	"github.com/openshift/installer/pkg/asset"
)

const (
	machineConfigFileName     = "99_openshift-machineconfig_%s.yaml"
	machineConfigPoolFileName = "99_openshift-machineconfigpool_%s.yaml"
)

var (
	machineConfigFileNamePattern = fmt.Sprintf(machineConfigFileName, "*")
)

// GenerateMachineConfigFiles creates manifest files containing the MachineConfigs.
func GenerateMachineConfigFiles(configs []*mcfgv1.MachineConfig, role, directory string) ([]*asset.File, error) {
	var ret []*asset.File
	for _, c := range configs {
		if c == nil {
			continue
		}
		configData, err := yaml.Marshal(c)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &asset.File{
			// Note that we should always be generating the role name in our MCs,
			// but just to ensure uniqueness we add the array index and the role too.
			Filename: filepath.Join(directory, fmt.Sprintf(machineConfigFileName, c.ObjectMeta.Name)),
			Data:     configData,
		})
	}
	if len(ret) == 0 {
		return nil, nil
	}
	return ret, nil
}

// GenerateMachineConfigPoolFile creates manifest files containing the MachineConfigPool.
func GenerateMachineConfigPoolFile(pool *mcfgv1.MachineConfigPool, role, directory string) (*asset.File, error) {
	configData, err := yaml.Marshal(pool)
	if err != nil {
		return nil, err
	}
	return &asset.File{
		Filename: filepath.Join(directory, fmt.Sprintf(machineConfigPoolFileName, pool.ObjectMeta.Name)),
		Data:     configData,
	}, nil
}

// IsMachineConfigManifest tests whether the specified filename is a MachineConfig manifest.
func IsMachineConfigManifest(filename string) (bool, error) {
	matched, err := filepath.Match(machineConfigFileNamePattern, filename)
	if err != nil {
		return false, err
	}
	return matched, nil
}

// IsMachineConfigPoolManifest tests whether the specified filename is a MachineConfigPool manifest.
func IsMachineConfigPoolManifest(filename string) bool {
	for _, role := range []string{PoolWorker, PoolMaster, PoolArbiter} {
		if filename == fmt.Sprintf(machineConfigFileName, role) {
			return true
		}
	}
	return false
}

// LoadMachineConfigs loads the MachineConfig manifests.
func LoadMachineConfigs(f asset.FileFetcher, role, directory string) ([]*asset.File, error) {
	return f.FetchByPattern(filepath.Join(directory, machineConfigFileNamePattern))
}

// LoadMachineConfigPool loads the MachineConfig manifest.
func LoadMachineConfigPool(f asset.FileFetcher, role, directory string) (*asset.File, error) {
	return f.FetchByName(filepath.Join(directory, fmt.Sprintf(machineConfigFileName, role)))
}
