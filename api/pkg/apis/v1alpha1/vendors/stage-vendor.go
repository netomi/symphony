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

package vendors

import (
	"context"

	"github.com/azure/symphony/api/pkg/apis/v1alpha1/managers/activations"
	"github.com/azure/symphony/api/pkg/apis/v1alpha1/managers/campaigns"
	"github.com/azure/symphony/api/pkg/apis/v1alpha1/managers/stage"
	"github.com/azure/symphony/api/pkg/apis/v1alpha1/model"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/managers"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/providers"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/providers/pubsub"
	"github.com/azure/symphony/coa/pkg/apis/v1alpha2/vendors"
	"github.com/azure/symphony/coa/pkg/logger"
)

var sLog = logger.NewLogger("coa.runtime")

type StageVendor struct {
	vendors.Vendor
	StageManager       *stage.StageManager
	CampaignsManager   *campaigns.CampaignsManager
	ActivationsManager *activations.ActivationsManager
}

func (s *StageVendor) GetInfo() vendors.VendorInfo {
	return vendors.VendorInfo{
		Version:  s.Vendor.Version,
		Name:     "Stage",
		Producer: "Microsoft",
	}
}

func (o *StageVendor) GetEndpoints() []v1alpha2.Endpoint {
	return []v1alpha2.Endpoint{}
}

func (s *StageVendor) Init(config vendors.VendorConfig, factories []managers.IManagerFactroy, providers map[string]map[string]providers.IProvider, pubsubProvider pubsub.IPubSubProvider) error {
	err := s.Vendor.Init(config, factories, providers, pubsubProvider)
	if err != nil {
		return err
	}
	for _, m := range s.Managers {
		if c, ok := m.(*stage.StageManager); ok {
			s.StageManager = c
		}
		if c, ok := m.(*campaigns.CampaignsManager); ok {
			s.CampaignsManager = c
		}
		if c, ok := m.(*activations.ActivationsManager); ok {
			s.ActivationsManager = c
		}
	}
	if s.StageManager == nil {
		return v1alpha2.NewCOAError(nil, "stage manager is not supplied", v1alpha2.MissingConfig)
	}
	if s.CampaignsManager == nil {
		return v1alpha2.NewCOAError(nil, "campaigns manager is not supplied", v1alpha2.MissingConfig)
	}
	if s.ActivationsManager == nil {
		return v1alpha2.NewCOAError(nil, "activations manager is not supplied", v1alpha2.MissingConfig)
	}
	s.Vendor.Context.Subscribe("activation", func(topic string, event v1alpha2.Event) error {
		var actData v1alpha2.ActivationData
		var aok bool
		if actData, aok = event.Body.(v1alpha2.ActivationData); !aok {
			return v1alpha2.NewCOAError(nil, "event body is not an activation job", v1alpha2.BadRequest)
		}
		campaign, err := s.CampaignsManager.GetSpec(context.Background(), actData.Campaign)
		if err != nil {
			return err
		}
		activation, err := s.ActivationsManager.GetSpec(context.Background(), actData.Activation)
		if err != nil {
			return err
		}

		evt, err := s.StageManager.HandleActivationEvent(context.Background(), actData, *campaign.Spec, activation)
		if err != nil {
			return err
		}
		if evt != nil {
			s.Vendor.Context.Publish("trigger", v1alpha2.Event{
				Body: *evt,
			})
		}
		return nil
	})
	s.Vendor.Context.Subscribe("trigger", func(topic string, event v1alpha2.Event) error {
		status := model.ActivationStatus{
			Stage:        "",
			NextStage:    "",
			Outputs:      nil,
			Status:       v1alpha2.Untouched,
			ErrorMessage: "",
			IsActive:     true,
		}
		triggerData := v1alpha2.ActivationData{}
		var aok bool
		if triggerData, aok = event.Body.(v1alpha2.ActivationData); !aok {
			err = v1alpha2.NewCOAError(nil, "event body is not an activation job", v1alpha2.BadRequest)
			status.Status = v1alpha2.BadRequest
			status.ErrorMessage = err.Error()
			status.IsActive = false
			return s.ActivationsManager.ReportStatus(context.Background(), triggerData.Activation, status)
		}
		campaign, err := s.CampaignsManager.GetSpec(context.Background(), triggerData.Campaign)
		if err != nil {
			status.Status = v1alpha2.BadRequest
			status.ErrorMessage = err.Error()
			status.IsActive = false
			return s.ActivationsManager.ReportStatus(context.Background(), triggerData.Activation, status)
		}
		status.Stage = triggerData.Stage
		status.ActivationGeneration = triggerData.ActivationGeneration
		status.Status = v1alpha2.Accepted
		err = s.ActivationsManager.ReportStatus(context.Background(), triggerData.Activation, status)
		if err != nil {
			return err
		}
		status, activation := s.StageManager.HandleTriggerEvent(context.Background(), *campaign.Spec, triggerData)
		err = s.ActivationsManager.ReportStatus(context.Background(), triggerData.Activation, status)
		if err != nil {
			return err
		}
		if activation != nil && status.Status != v1alpha2.Done {
			s.Vendor.Context.Publish("trigger", v1alpha2.Event{
				Body: *activation,
			})
		}
		return nil
	})
	return nil
}
