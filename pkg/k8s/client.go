package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
	"runtimectl/dao"
	"runtimectl/pkg/util"
	"time"
)

type K8sClient struct {
	DynamicClient *dynamic.DynamicClient
	Client        *kubernetes.Clientset
}

func Init(path string) *K8sClient {
	dynamic, err := dynamicConfigInit(path)
	if err != nil {
		log.Fatalf("Error building k8s client: %s", err.Error())
	}
	c, err := configInit(path)
	if err != nil {
		log.Fatalf("Error building k8s client: %s", err.Error())
	}
	client := &K8sClient{
		DynamicClient: dynamic,
		Client:        c,
	}
	return client
}

func (sdk *K8sClient) Patch() error {
	crdList, err := sdk.getAllDevbox()
	for _, crd := range crdList.Items {
		log.Println("Patching devbox ", crd.GetName())
		name, version, image, err := sdk.getRuntimeNameAndVersionAndImage(crd.GetName(), crd.GetNamespace())
		if err != nil {
			return err
		}

		templateUID, _ := dao.GetTemplateID(name, version, image)
		patchData := map[string]interface{}{
			"spec": map[string]interface{}{
				"templateID": templateUID,
			},
			"metadata": map[string]interface{}{
				"labels": map[string]string{
					"devbox.sealos.io/patched": "true",
				},
			},
		}
		patchBytes, err := json.Marshal(patchData)
		if err != nil {
			return err
		}

		_, err = sdk.DynamicClient.Resource(getDevboxSchema()).Namespace(crd.GetNamespace()).Patch(context.Background(), crd.GetName(), types.MergePatchType, patchBytes, metav1.PatchOptions{})
		if err != nil {
			log.Println("Error patching devbox ", crd.GetName())
			return err
		}
		log.Println("Patched devbox ", crd.GetName())
	}
	return err
}

func (sdk *K8sClient) SyncToDB() error {
	runtimeList, err := sdk.GetAllRuntime()
	if err != nil {
		return err
	}
	for _, r := range runtimeList.Items {
		active, _, _ := unstructured.NestedString(r.Object, "spec", "state")
		class, _, _ := unstructured.NestedString(r.Object, "spec", "classRef")
		runtimeClass, err := sdk.GetRuntimeClass(class)
		if err != nil {
			log.Println("Error getting runtime class:", err)
			log.Println("Error runtimeclass:", runtimeClass)
			log.Println("Error runtime:", class)
		}
		kind, _, err := unstructured.NestedString(runtimeClass.Object, "spec", "kind")
		kind = parseKind(kind)

		version, _, _ := unstructured.NestedString(r.Object, "spec", "version")
		image, _, _ := unstructured.NestedString(r.Object, "spec", "config", "image")
		deleteTime, found, _ := unstructured.NestedString(r.Object, "spec", "runtimeVersion")

		if !found || deleteTime == "" {
			deleteTime = util.RandomDateString()
		}
		parsedTime, err := time.Parse("2006-01-02-1504", deleteTime)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return err
		}

		config, err := sdk.GetRuntimeConfig(r)

		if err := dao.CreateOrUpdateTemplateRepository(class, kind); err != nil {
			fmt.Println("Error creating or updating template repository:", err)
			return err
		}
		t := dao.GetTemplateRepository(class)
		if err := dao.CreateOrUpdateTemplate(version, t.UID, image, config, active, parsedTime); err != nil {
			fmt.Println("Error creating or updating template:", err)
			return err
		}
	}
	return nil
}

func parseKind(kind string) string {
	switch kind {
	case "Framework":
		return "FRAMEWORK"
	case "Language":
		return "LANGUAGE"
	case "OS":
		return "OS"
	case "Custom":
		return "CUSTOM"
	}
	return ""
}

func (sdk *K8sClient) GetRuntimeConfig(r unstructured.Unstructured) (string, error) {
	configData, found, err := unstructured.NestedMap(r.Object, "spec", "config")
	if err != nil || !found {
		fmt.Println("spec.config not found or error occurred")
		return "", err
	}
	delete(configData, "image")
	config, err := json.MarshalIndent(configData, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	return string(config), nil
}

func (sdk *K8sClient) getAllDevbox() (*unstructured.UnstructuredList, error) {
	// Retrieve all devbox resources
	allDevboxes, err := sdk.DynamicClient.Resource(getDevboxSchema()).Namespace("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Filter out devboxes that have the label "devbox.sealos.io/patched" set to "true"
	filteredDevboxes := &unstructured.UnstructuredList{}
	for _, devbox := range allDevboxes.Items {
		labels := devbox.GetLabels()
		if labels == nil || labels["devbox.sealos.io/patched"] != "true" {
			filteredDevboxes.Items = append(filteredDevboxes.Items, devbox)
		}
	}
	return filteredDevboxes, nil
}

func (sdk *K8sClient) GetAllRuntime() (*unstructured.UnstructuredList, error) {
	return sdk.DynamicClient.Resource(getRuntimeSchema()).Namespace("").List(context.Background(), metav1.ListOptions{})
}

func (sdk *K8sClient) GetRuntimeClass(name string) (*unstructured.Unstructured, error) {
	return sdk.DynamicClient.Resource(getRuntimeClassSchema()).Namespace("devbox-system").Get(context.Background(), name, metav1.GetOptions{})
}

func (sdk *K8sClient) getRuntimeNameAndVersionAndImage(name, namespace string) (string, string, string, error) {
	runtimeRef, err := sdk.getDevboxRuntimeRef(name, namespace)
	if err != nil {
		log.Printf("Error getting CRD instance: %s\n", err.Error())
		return "", "", "", err
	}
	unstructuredObj, err := sdk.DynamicClient.Resource(getRuntimeSchema()).Namespace("devbox-system").Get(context.Background(), runtimeRef, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting runtime CRD instance: %s\n", err.Error())
		return "", "", "", err
	}
	n, found, err := unstructured.NestedString(unstructuredObj.Object, "spec", "classRef")
	if err != nil || !found {
		log.Println("spec field not found or error occurred")
		return "", "", "", err
	}
	version, found, err := unstructured.NestedString(unstructuredObj.Object, "spec", "version")
	if err != nil || !found {
		log.Println("spec field not found or error occurred")
		return "", "", "", err
	}
	image, found, err := unstructured.NestedString(unstructuredObj.Object, "spec", "config", "image")
	if err != nil || !found {
		log.Println("spec field not found or error occurred")
		return "", "", "", err
	}
	return n, version, image, nil
}

func (sdk *K8sClient) getDevboxRuntimeRef(name, namespace string) (string, error) {
	unstructuredObj, err := sdk.DynamicClient.Resource(getDevboxSchema()).Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting devbox CRD instance: %s\n", err.Error())
		return "", err
	}
	runtimeRefName, found, err := unstructured.NestedString(unstructuredObj.Object, "spec", "runtimeRef", "name")
	if err != nil || !found {
		log.Println("spec field not found or error occurred")
		return "", err
	}
	return runtimeRefName, nil
}

func getDevboxSchema() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "devbox.sealos.io",
		Version:  "v1alpha1",
		Resource: "devboxes",
	}
}

func getRuntimeSchema() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "devbox.sealos.io",
		Version:  "v1alpha1",
		Resource: "runtimes",
	}
}

func getRuntimeClassSchema() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "devbox.sealos.io",
		Version:  "v1alpha1",
		Resource: "runtimeclasses",
	}
}
