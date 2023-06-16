package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

func main() {
	err := createPod("my-pod")
	if err != nil {
		panic(err)
	}

	time.Sleep(10 * time.Second)

	err = deletePod("my-pod")
	if err != nil {
		panic(err)
	}
}

func deletePod(name string) error {
	reqDelete, err := http.NewRequest(
		"DELETE",
		"http://127.0.0.1:8001/api/v1/namespaces/default/pods/"+name,
		nil,
	)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(reqDelete) // ➍
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("Delete Failed")
		return err
	} else {
		fmt.Println("Delete Successfully ") // ➑
		return nil
	}

}

func createPod(name string) error {
	pod := createPodObject(name) // ➊

	serializer := getJSONSerializer()
	postBody, err := serializePodObject(serializer, pod) // ➋
	if err != nil {
		return err
	}

	reqCreate, err := buildPostRequest(postBody) // ➌
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(reqCreate) // ➍
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 300 { // ➏
		fmt.Println("Create Successfully ") // ➑
	} else {
		fmt.Println("Create Failed")
	}
	return nil
}

func createPodObject(name string) *corev1.Pod { // ➊
	pod := corev1.Pod{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "runtime",
					Image: "nginx",
				},
			},
		},
	}

	pod.SetName(name)
	pod.SetLabels(map[string]string{
		"app.kubernetes.io/component": "my-component",
		"app.kubernetes.io/name":      "a-name",
	})
	return &pod
}

func serializePodObject( // ➋
	serializer runtime.Serializer,
	pod *corev1.Pod,
) (
	io.Reader,
	error,
) {
	var buf bytes.Buffer
	err := serializer.Encode(pod, &buf)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

func buildPostRequest( // ➌
	body io.Reader,
) (
	*http.Request,
	error,
) {
	reqCreate, err := http.NewRequest(
		"POST",
		"http://127.0.0.1:8001/api/v1/namespaces/default/pods",
		body,
	)
	if err != nil {
		return nil, err
	}
	reqCreate.Header.Add(
		"Accept",
		"application/json",
	)
	reqCreate.Header.Add(
		"Content-Type",
		"application/json",
	)
	return reqCreate, nil
}

func getJSONSerializer() runtime.Serializer {
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(
		schema.GroupVersion{
			Group:   "",
			Version: "v1",
		},
		&corev1.Pod{},
		&metav1.Status{},
	)
	return kjson.NewSerializerWithOptions(
		kjson.SimpleMetaFactory{},
		nil,
		scheme,
		kjson.SerializerOptions{},
	)
}
