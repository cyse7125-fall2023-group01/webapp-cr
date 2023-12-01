package main

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

func createCustomResource(check HTTPCheck) error {
	// Use in-cluster configuration
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	// Create a dynamic client
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	// Your custom resource definition
	u := &unstructured.Unstructured{}
	u.Object = map[string]interface{}{
		"apiVersion": "monitoring.httpcheck.io/v1alpha1",
		"kind":       "CronJob",
		"metadata": map[string]interface{}{
			"name": check.Name,
		},
		"spec": map[string]interface{}{
			"ID":                        check.ID,
			"name":                      check.Name,
			"uri":                       check.URI,
			"is_paused":                 check.IsPaused,
			"num_retries":               check.NumRetries,
			"uptime_sla":                check.UptimeSLA,
			"response_time_sla":         check.ResponseTimeSLA,
			"use_ssl":                   check.UseSSL,
			"response_status_code":      check.ResponseStatusCode,
			"check_interval_in_seconds": check.CheckIntervalSeconds,
		},
	}

	// Create the Custom Resource
	result, err := client.Resource(schema.GroupVersionResource{
		Group:    "monitoring.httpcheck.io",
		Version:  "v1alpha1",
		Resource: "cronjobs",
	}).Namespace("webapp").Create(context.Background(), u, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Println("Created Custom Resource:", result)
	return nil
}

func deleteCustomResource(Name string) error {
	// Use in-cluster configuration
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	// Create a dynamic client
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	// Delete the Custom Resource
	err = client.Resource(schema.GroupVersionResource{
		Group:    "monitoring.httpcheck.io",
		Version:  "v1alpha1",
		Resource: "cronjobs",
	}).Namespace("webapp").Delete(context.Background(), Name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Println("Deleted Custom Resource:", Name)
	return nil
}

func updateCustomResource(check HTTPCheck) error {
	// Use in-cluster configuration
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	// Create a dynamic client
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	// Retrieve the current resource version
	existingResource, err := client.Resource(schema.GroupVersionResource{
		Group:    "monitoring.httpcheck.io",
		Version:  "v1alpha1",
		Resource: "cronjobs",
	}).Namespace("webapp").Get(context.Background(), check.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// Extract the resource version
	resourceVersion, found, err := unstructured.NestedString(existingResource.Object, "metadata", "resourceVersion")
	if err != nil || !found {
		return fmt.Errorf("failed to retrieve resource version: %v", err)
	}

	// Create the updated Custom Resource
	u := &unstructured.Unstructured{}
	u.Object = map[string]interface{}{
		"apiVersion": "monitoring.httpcheck.io/v1alpha1",
		"kind":       "CronJob",
		"metadata": map[string]interface{}{
			"name":            check.Name,
			"resourceVersion": resourceVersion, // Include the resource version for the update
		},
		"spec": map[string]interface{}{
			"ID":                        check.ID,
			"name":                      check.Name,
			"uri":                       check.URI,
			"is_paused":                 check.IsPaused,
			"num_retries":               check.NumRetries,
			"uptime_sla":                check.UptimeSLA,
			"response_time_sla":         check.ResponseTimeSLA,
			"use_ssl":                   check.UseSSL,
			"response_status_code":      check.ResponseStatusCode,
			"check_interval_in_seconds": check.CheckIntervalSeconds,
		},
	}

	// Update the Custom Resource
	result, err := client.Resource(schema.GroupVersionResource{
		Group:    "monitoring.httpcheck.io",
		Version:  "v1alpha1",
		Resource: "cronjobs",
	}).Namespace("webapp").Update(context.Background(), u, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	log.Println("Updated Custom Resource:", result)
	return nil
}
