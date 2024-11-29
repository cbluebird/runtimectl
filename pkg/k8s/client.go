package sdk

import (
	"context"
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"log"
	"runtimectl/dao"
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
		name, version, image, err := sdk.getRuntimeNameAndVersionAndImage(crd.GetName(), crd.GetNamespace())
		if err != nil {
			return err
		}

		templateUID, _ := dao.GetTemplateID(name, version, image)
		patchData := map[string]interface{}{
			"spec": map[string]interface{}{
				"templateID": templateUID,
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

func (sdk *K8sClient) getAllDevbox() (*unstructured.UnstructuredList, error) {
	return sdk.DynamicClient.Resource(getDevboxSchema()).Namespace("").List(context.Background(), metav1.ListOptions{})
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
