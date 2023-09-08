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

package sideload

import (
	"context"
	"encoding/json"
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/azure/symphony/api/pkg/apis/v1alpha1/model"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/contexts"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/observability"
	observ_utils "github.com/azure/symphony/coa/pkg/apis/v1alpha2/observability/utils"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/providers"
	"github.com/azure/symphony/coa/pkg/logger"
)

var sLog = logger.NewLogger("coa.runtime")

type Win10SideLoadProviderConfig struct {
	Name                string `json:"name"`
	IPAddress           string `json:"ipAddress"`
	Pin                 string `json:"pin,omitempty"`
	WinAppDeployCmdPath string `json:"winAppDeployCmdPath"`
	NetworkUser         string `json:"networkUser,omitempty"`
	NetworkPassword     string `json:"networkPassword,omitempty"`
	Silent              bool   `json:"silent,omitempty"`
}

type Win10SideLoadProvider struct {
	Config  Win10SideLoadProviderConfig
	Context *contexts.ManagerContext
}

func Win10SideLoadProviderConfigFromMap(properties map[string]string) (Win10SideLoadProviderConfig, error) {
	ret := Win10SideLoadProviderConfig{}
	if v, ok := properties["name"]; ok {
		ret.Name = v
	}
	if v, ok := properties["ipAddress"]; ok {
		ret.IPAddress = v
	} else {
		ret.IPAddress = "localhost"
	}
	if v, ok := properties["pin"]; ok {
		ret.Pin = v
	}
	if v, ok := properties["winAppDeployCmdPath"]; ok {
		ret.WinAppDeployCmdPath = v
	} else {
		ret.WinAppDeployCmdPath = "c:\\Program Files (x86)\\Windows Kits\\10\\bin\\10.0.19041.0\\x86\\WinAppDeployCmd.exe"
	}
	if v, ok := properties["networkUser"]; ok {
		ret.NetworkUser = v
	}
	if v, ok := properties["networkPassword"]; ok {
		ret.NetworkPassword = v
	}
	if v, ok := properties["silent"]; ok {
		bVal, err := strconv.ParseBool(v)
		if err != nil {
			ret.Silent = false
		} else {
			ret.Silent = bVal
		}
	}
	return ret, nil
}
func (i *Win10SideLoadProvider) InitWithMap(properties map[string]string) error {
	config, err := Win10SideLoadProviderConfigFromMap(properties)
	if err != nil {
		return err
	}
	return i.Init(config)
}

func (i *Win10SideLoadProvider) Init(config providers.IProviderConfig) error {
	_, span := observability.StartSpan("Win 10 Sideload Provider", context.Background(), &map[string]string{
		"method": "Init",
	})
	sLog.Info("~~~ Win 10 Sideload Provider ~~~ : Init()")

	updateConfig, err := toWin10SideLoadProviderConfig(config)
	if err != nil {
		return errors.New("expected Win10SideLoadProviderConfig")
	}
	i.Config = updateConfig

	observ_utils.CloseSpanWithError(span, nil)
	return nil
}
func toWin10SideLoadProviderConfig(config providers.IProviderConfig) (Win10SideLoadProviderConfig, error) {
	ret := Win10SideLoadProviderConfig{}
	data, err := json.Marshal(config)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(data, &ret)
	return ret, err
}
func (i *Win10SideLoadProvider) Get(ctx context.Context, deployment model.DeploymentSpec, references []model.ComponentStep) ([]model.ComponentSpec, error) {
	_, span := observability.StartSpan("Win 10 Sideload Provider", context.Background(), &map[string]string{
		"method": "Get",
	})
	sLog.Infof("~~~ Win 10 Sideload Provider ~~~ : getting artifacts: %s - %s", deployment.Instance.Scope, deployment.Instance.Name)

	params := make([]string, 0)
	params = append(params, "list")
	params = append(params, "-ip")
	params = append(params, i.Config.IPAddress)
	if i.Config.Pin != "" {
		params = append(params, "-pin")
		params = append(params, i.Config.Pin)
	}

	out, err := exec.Command(i.Config.WinAppDeployCmdPath, params...).Output()

	if err != nil {
		observ_utils.CloseSpanWithError(span, err)
		return nil, err
	}
	str := string(out)
	lines := strings.Split(str, "\r\n")

	desired := deployment.GetComponentSlice()

	re := regexp.MustCompile(`^(\w+\.)+\w+$`)
	ret := make([]model.ComponentSpec, 0)
	for _, line := range lines {
		if re.Match([]byte(line)) {
			mLine := line
			if strings.LastIndex(line, "__") > 0 {
				mLine = line[:strings.LastIndex(line, "__")]
			}
			for _, component := range desired {
				if component.Name == mLine {
					ret = append(ret, model.ComponentSpec{
						Name: line,
						Type: "win.uwp",
					})
				}
			}
		}
	}

	observ_utils.CloseSpanWithError(span, nil)
	return ret, nil
}
func (i *Win10SideLoadProvider) Apply(ctx context.Context, deployment model.DeploymentSpec, step model.DeploymentStep, isDryRun bool) (map[string]model.ComponentResultSpec, error) {
	_, span := observability.StartSpan("Win 10 Sideload Provider", ctx, &map[string]string{
		"method": "Apply",
	})
	sLog.Infof("~~~ Win 10 Sideload Provider ~~~ : applying artifacts: %s - %s", deployment.Instance.Scope, deployment.Instance.Name)

	components := step.GetComponents()
	err := i.GetValidationRule(ctx).Validate(components)
	if err != nil {
		observ_utils.CloseSpanWithError(span, err)
		return nil, err
	}
	if isDryRun {
		observ_utils.CloseSpanWithError(span, nil)
		return nil, nil
	}

	ret := step.PrepareResultMap()
	components = step.GetUpdatedComponents()
	if len(components) > 0 {
		for _, component := range components {
			if path, ok := component.Properties["app.package.path"].(string); ok {
				params := make([]string, 0)
				params = append(params, "install")
				params = append(params, "-ip")
				params = append(params, i.Config.IPAddress)
				if i.Config.Pin != "" {
					params = append(params, "-pin")
					params = append(params, i.Config.Pin)
				}
				params = append(params, "-file")
				params = append(params, path)

				cmd := exec.Command(i.Config.WinAppDeployCmdPath, params...)
				err := cmd.Run()
				if err != nil {
					ret[component.Name] = model.ComponentResultSpec{
						Status:  v1alpha2.UpdateFailed,
						Message: err.Error(),
					}
					observ_utils.CloseSpanWithError(span, err)
					if i.Config.Silent {
						return ret, nil
					} else {
						return ret, err
					}
				}
			}
		}
	}
	components = step.GetDeletedComponents()
	if len(components) > 0 {
		for _, component := range components {
			if component.Name != "" {
				params := make([]string, 0)
				params = append(params, "uninstall")
				params = append(params, "-ip")
				params = append(params, i.Config.IPAddress)
				if i.Config.Pin != "" {
					params = append(params, "-pin")
					params = append(params, i.Config.Pin)
				}
				params = append(params, "-package")

				name := component.Name

				// TODO: this is broken due to the refactor, the current reference is no longer available
				// for _, ref := range currentRef {
				// 	if ref.Name == name || strings.HasPrefix(ref.Name, name) {
				// 		name = ref.Name
				// 		break
				// 	}
				// }

				params = append(params, name)

				cmd := exec.Command(i.Config.WinAppDeployCmdPath, params...)
				err := cmd.Run()
				if err != nil {
					observ_utils.CloseSpanWithError(span, err)
					if i.Config.Silent {
						return ret, nil
					} else {
						return ret, err
					}
				}

			}
		}
	}
	observ_utils.CloseSpanWithError(span, nil)
	return ret, nil
}

func (i *Win10SideLoadProvider) NeedsUpdate(ctx context.Context, desired []model.ComponentSpec, current []model.ComponentSpec) bool {
	for _, d := range desired {
		found := false
		for _, c := range current {
			if c.Name == d.Name || strings.HasPrefix(c.Name, d.Name) {
				found = true
			}
		}
		if !found {
			return true
		}
	}
	return false
}
func (i *Win10SideLoadProvider) NeedsRemove(ctx context.Context, desired []model.ComponentSpec, current []model.ComponentSpec) bool {
	for _, d := range desired {
		for _, c := range current {
			if c.Name == d.Name || strings.HasPrefix(c.Name, d.Name) {
				return true
			}
		}
	}
	return false
}

func (*Win10SideLoadProvider) GetValidationRule(ctx context.Context) model.ValidationRule {
	return model.ValidationRule{
		RequiredProperties:    []string{},
		OptionalProperties:    []string{},
		RequiredComponentType: "",
		RequiredMetadata:      []string{},
		OptionalMetadata:      []string{},
		ChangeDetectionProperties: []model.PropertyDesc{
			{Name: "", IsComponentName: true, IgnoreCase: true, PrefixMatch: true},
		},
	}
}
