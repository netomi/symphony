/*

	MIT License

	Copyright (c) Microsoft Corporation.

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE

*/

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeploymentDeepEquals(t *testing.T) {
	deployment1 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	deployment2 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	res, err := deployment1.DeepEquals(deployment2)
	assert.Nil(t, err)
	assert.True(t, res)
}

func TestDeploymentDeepEqualsOneEmpty(t *testing.T) {
	deployment1 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	res, err := deployment1.DeepEquals(nil)
	assert.Errorf(t, err, "parameter is not a DeploymentSpec type")
	assert.False(t, res)
}

func TestDeploymentDeepEqualsSolutionNameNotMatch(t *testing.T) {
	deployment1 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	deployment2 := DeploymentSpec{
		SolutionName: "SolutionName1",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	res, err := deployment1.DeepEquals(deployment2)
	assert.Nil(t, err)
	assert.False(t, res)
}

func TestDeploymentDeepEqualsSolutionNotMatch(t *testing.T) {
	deployment1 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	deployment2 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName1",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	res, err := deployment1.DeepEquals(deployment2)
	assert.Nil(t, err)
	assert.False(t, res)
}

func TestDeploymentDeepEqualsInstanceNotMatch(t *testing.T) {
	deployment1 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	deployment2 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName1",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	res, err := deployment1.DeepEquals(deployment2)
	assert.Nil(t, err)
	assert.False(t, res)
}

func TestDeploymentDeepEqualsTargetsNotMatch(t *testing.T) {
	deployment1 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	deployment2 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName1",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo1": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	res, err := deployment1.DeepEquals(deployment2)
	assert.Nil(t, err)
	assert.False(t, res)
}

func TestDeploymentDeepEqualsDevicesNotMatch(t *testing.T) {
	deployment1 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	deployment2 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName1",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	res, err := deployment1.DeepEquals(deployment2)
	assert.Nil(t, err)
	assert.False(t, res)
}

func TestDeploymentDeepEqualsComponentStartIndexNotMatch(t *testing.T) {
	deployment1 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	deployment2 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 1,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	res, err := deployment1.DeepEquals(deployment2)
	assert.Nil(t, err)
	assert.False(t, res)
}

func TestDeploymentDeepEqualsComponentEndIndexNotMatch(t *testing.T) {
	deployment1 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	deployment2 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   1,
		ActiveTarget:        "ActiveTarget",
	}
	res, err := deployment1.DeepEquals(deployment2)
	assert.Nil(t, err)
	assert.False(t, res)
}

func TestDeploymentDeepEqualsActiveTargetNotMatch(t *testing.T) {
	deployment1 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget",
	}
	deployment2 := DeploymentSpec{
		SolutionName: "SolutionName",
		Solution: SolutionSpec{
			DisplayName: "SolutionDisplayName",
		},
		Instance: InstanceSpec{
			Name: "InstanceName",
		},
		Targets: map[string]TargetSpec{
			"foo": {
				DisplayName: "TargetName",
			},
		},
		Devices: []DeviceSpec{{
			DisplayName: "DeviceName",
		}},
		Assignments: map[string]string{
			"foo": "bar",
		},
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
		ActiveTarget:        "ActiveTarget1",
	}
	res, err := deployment1.DeepEquals(deployment2)
	assert.Nil(t, err)
	assert.False(t, res)
}

func TestGetComponentSlice(t *testing.T) {
	deployment := DeploymentSpec{
		ComponentStartIndex: 0,
		ComponentEndIndex:   0,
	}
	res := deployment.GetComponentSlice()
	assert.Equal(t, 0, len(res))
}

func TestGetComponentSliceWithValues(t *testing.T) {
	deployment := DeploymentSpec{
		ComponentStartIndex: 1,
		ComponentEndIndex:   2,
		Solution: SolutionSpec{
			Components: []ComponentSpec{
				{Name: "Component1"},
				{Name: "Component2"},
				{Name: "Component3"},
				{Name: "Component4"},
				{Name: "Component5"},
			},
		},
	}
	res := deployment.GetComponentSlice()
	assert.Equal(t, 1, len(res))
}

func TestMapsEqual(t *testing.T) {
	map1 := map[string]TargetSpec{
		"foo": {
			DisplayName: "TargetName",
		},
	}
	map2 := map[string]TargetSpec{
		"foo": {},
	}
	res := mapsEqual(map1, map2, nil)
	assert.False(t, res)
}
