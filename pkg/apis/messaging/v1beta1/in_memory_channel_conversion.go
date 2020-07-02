/*
Copyright 2020 The Knative Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"context"
	"fmt"

	eventingduckv1 "knative.dev/eventing/pkg/apis/duck/v1"
	eventingduckv1beta1 "knative.dev/eventing/pkg/apis/duck/v1beta1"
	"knative.dev/pkg/apis"

	"knative.dev/eventing/pkg/apis/messaging"
	v1 "knative.dev/eventing/pkg/apis/messaging/v1"
)

// ConvertTo implements apis.Convertible
// Converts source (from v1beta1.InMemoryChannel) into v1.InMemoryChannel
func (source *InMemoryChannel) ConvertTo(ctx context.Context, obj apis.Convertible) error {
	switch sink := obj.(type) {
	case *v1.InMemoryChannel:
		sink.ObjectMeta = source.ObjectMeta
		if sink.Annotations == nil {
			sink.Annotations = make(map[string]string)
		}
		sink.Annotations[messaging.SubscribableDuckVersionAnnotation] = "v1"
		if err := source.Status.ConvertTo(ctx, &sink.Status); err != nil {
			return err
		}
		return source.Spec.ConvertTo(ctx, &sink.Spec)
	default:
		return fmt.Errorf("unknown version, got: %T", sink)
	}
}

// ConvertTo helps implement apis.Convertible
func (source *InMemoryChannelSpec) ConvertTo(ctx context.Context, sink *v1.InMemoryChannelSpec) error {
	sink.SubscribableSpec = eventingduckv1.SubscribableSpec{}
	if err := source.SubscribableSpec.ConvertTo(ctx, &sink.SubscribableSpec); err != nil {
		return err
	}
	if source.Delivery != nil {
		sink.Delivery = &eventingduckv1.DeliverySpec{}
		return source.Delivery.ConvertTo(ctx, sink.Delivery)
	}
	return nil
}

// ConvertTo helps implement apis.Convertible
func (source *InMemoryChannelStatus) ConvertTo(ctx context.Context, sink *v1.InMemoryChannelStatus) error {
	source.Status.ConvertTo(ctx, &sink.Status)
	sink.AddressStatus = source.AddressStatus
	sink.SubscribableStatus = eventingduckv1.SubscribableStatus{}
	if err := source.SubscribableStatus.ConvertTo(ctx, &sink.SubscribableStatus); err != nil {
		return err
	}
	return nil
}

// ConvertFrom implements apis.Convertible.
// Converts obj v1.InMemoryChannel into v1beta1.InMemoryChannel
func (sink *InMemoryChannel) ConvertFrom(ctx context.Context, obj apis.Convertible) error {
	switch source := obj.(type) {
	case *v1.InMemoryChannel:
		sink.ObjectMeta = source.ObjectMeta
		if err := sink.Status.ConvertFrom(ctx, source.Status); err != nil {
			return err
		}
		if err := sink.Spec.ConvertFrom(ctx, source.Spec); err != nil {
			return err
		}
		if sink.Annotations == nil {
			sink.Annotations = make(map[string]string)
		}
		sink.Annotations[messaging.SubscribableDuckVersionAnnotation] = "v1beta1"
		return nil
	default:
		return fmt.Errorf("unknown version, got: %T", source)
	}
}

// ConvertFrom helps implement apis.Convertible
func (sink *InMemoryChannelSpec) ConvertFrom(ctx context.Context, source v1.InMemoryChannelSpec) error {
	if source.Delivery != nil {
		sink.Delivery = &eventingduckv1beta1.DeliverySpec{}
		if err := sink.Delivery.ConvertFrom(ctx, source.Delivery); err != nil {
			return err
		}
	}
	sink.SubscribableSpec = eventingduckv1beta1.SubscribableSpec{}
	if err := sink.SubscribableSpec.ConvertFrom(ctx, source.SubscribableSpec); err != nil {
		return err
	}
	return nil
}

// ConvertFrom helps implement apis.Convertible
func (sink *InMemoryChannelStatus) ConvertFrom(ctx context.Context, source v1.InMemoryChannelStatus) error {
	source.Status.ConvertTo(ctx, &sink.Status)
	sink.AddressStatus = source.AddressStatus
	sink.SubscribableStatus = eventingduckv1beta1.SubscribableStatus{}
	if err := sink.SubscribableStatus.ConvertFrom(ctx, source.SubscribableStatus); err != nil {
		return err
	}
	return nil
}
