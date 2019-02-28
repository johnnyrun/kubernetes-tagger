/*
 * Author: Alexandre Havrileck (Oxyno-zeta)
 * Date: 28/02/2019
 * Licence: See Readme
 */
package business

import (
	"fmt"

	"github.com/oxyno-zeta/kubernetes-tagger/pkg/kubernetes-tagger/config"
	"github.com/oxyno-zeta/kubernetes-tagger/pkg/kubernetes-tagger/resources"
	"github.com/oxyno-zeta/kubernetes-tagger/pkg/kubernetes-tagger/rules"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type BusinessContext struct {
	KubernetesClient  *kubernetes.Clientset
	MainConfiguration *config.MainConfiguration
	Rules             []*rules.Rule
}

func (context *BusinessContext) handlePersistentVolumeAdd(obj interface{}) {
	pv := obj.(*v1.PersistentVolume)
	context.runForPV(pv)
}
func (context *BusinessContext) handlePersistentVolumeDelete(obj interface{}) {
	// Nothing to do
}

func (context *BusinessContext) handlePersistentVolumeUpdate(old, current interface{}) {
	currentPersistentVolume := current.(*v1.PersistentVolume)
	context.runForPV(currentPersistentVolume)
}

func (context *BusinessContext) runForPV(pv *v1.PersistentVolume) {
	resource, err := resources.New(context.KubernetesClient, pv, context.MainConfiguration.Config)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Check if resource can be processed
	if !resource.CanBeProcessed() {
		return
	}

	// Get actual tags
	actualTags, err := resource.GetActualTags()
	if err != nil {
		fmt.Println(err)
		return
	}

	availableTagValues, err := resource.GetAvailableTagValues()
	if err != nil {
		fmt.Println(err)
		return
	}

	delta, err := rules.CalculateTags(actualTags, availableTagValues, context.Rules)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = resource.ManageTags(delta)
	if err != nil {
		fmt.Println(err)
		return
	}
}