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

func (sdk *K8sClient) getAllDevbox() (*unstructured.UnstructuredList, error) {
	return sdk.DynamicClient.Resource(getDevboxSchema()).Namespace("").List(context.Background(), metav1.ListOptions{})
}

func (sdk *K8sClient) Patch() error {
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"templateID": "newValue",
		},
	}
	patchBytes, err := json.Marshal(patchData)
	if err != nil {
		return err
	}

	crdList, err := sdk.getAllDevbox()

	for _, crd := range crdList.Items {
		sdk.DynamicClient.Resource(getDevboxSchema()).Namespace(crd.GetNamespace()).Patch(context.Background(), crd.GetName(), types.MergePatchType, patchBytes, metav1.PatchOptions{})

		fmt.Printf("CRD Instance Name: %s\n", crd.GetName())
	}

	return err
}

func (sdk *K8sClient) getRuntimeRef(name, namespace string) (string, error) {
	unstructuredObj, err := sdk.DynamicClient.Resource(getDevboxSchema()).Namespace(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		log.Printf("Error getting CRD instance: %s\n", err.Error())
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
