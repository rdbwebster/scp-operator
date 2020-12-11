package kclient

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	model "github.com/rdbwebster/scp-operator/model"
	"github.com/rdbwebster/scp-operator/stacktrace"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func ConnectToCluster(ci model.ClusterInfo) (int, error) {

	config, err := clientcmd.BuildConfigFromFlags(ci.Spec.Url, "")

	config.CAData = []byte(ci.Spec.CertAuth)
	config.CertData = []byte(ci.Spec.Cert)
	config.KeyData = []byte(ci.Spec.Token)

	// create the clientset

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		return 0, err
	}

	fmt.Printf("%+v\n", config)

	count, _ := listPods(clientset)

	return count, err

}

func getThumbprint(url string) (string, error) {
	// Parse cmdline arguments using flag package
	//	server := flag.String("server", host, "Server to ping")
	//	port := flag.Uint("port", 443, "Port that has TLS")/
	//	flag.Parse()

	//	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", *server, *port), &tls.Config{})

	// Get the ConnectionState struct as that's the one which gives us x509.Certificate struct
	cert, err := getCertificate(url)
	if err != nil {
		return "", err
	}

	fingerprint := md5.Sum(cert.Raw)

	var buf bytes.Buffer
	for i, f := range fingerprint {
		if i > 0 {
			fmt.Fprintf(&buf, ":")
		}
		fmt.Fprintf(&buf, "%02X", f)
	}
	//fmt.Printf("Fingerprint for %s: %s", *server, buf.String())

	fmt.Printf("Fingerprint for %s: %s", url, buf.String())

	return buf.String(), err
}

func createScpCluster() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	namespace := "default"

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	scpclusterRes := schema.GroupVersionResource{Group: "webapp.my.domain", Version: "v1", Resource: "scpclusters"}

	scpcluster := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "webapp.my.domain/v1",
			"kind":       "SCPcluster",
			"metadata": map[string]interface{}{
				"name": "scpcluster-sample-2",
			},
			"spec": map[string]interface{}{
				"id":        5,
				"name":      "cluster1",
				"url":       "localhost:443",
				"token":     "==== START ====",
				"cert":      "==== START ====",
				"auth":      "==== START ====",
				"connected": metav1.Now(),
			},
		},
	}

	// Create Resource
	fmt.Println("Creating deployment...")
	result, err := client.Resource(scpclusterRes).Namespace(namespace).Create(context.TODO(), scpcluster, metav1.CreateOptions{})

	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetName())
}

func getCertificate(url string) (*x509.Certificate, error) {

	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	//	conn, err := tls.Dial("tcp", url, &tls.Config{})
	conn, err := tls.Dial("tcp", url, conf)
	if err != nil {
		return &x509.Certificate{}, err
	}

	// Get the ConnectionState struct as that's the one which gives us x509.Certificate struct
	cert := conn.ConnectionState().PeerCertificates[0]

	conn.Close()

	//  crtkit.ServerCertPem = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	//	crtkit.ServerKeyPem = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})
	//	Goose.Generator.Logf(4, "PEM Certificate: %s", crtkit.ServerCertPem)

	return cert, err

}

// asn1Data starting with -----BEGIN CERTIFICATE-----
func convertPemToDer(asn1Data string) ([]byte, error) {
	var block *pem.Block
	var pemByte []byte

	block, pemByte = pem.Decode([]byte(asn1Data))

	fmt.Println(block)

	fmt.Println(pemByte)

	// just a single certificate
	cert, err := x509.ParseCertificate(block.Bytes)

	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	return cert.Raw, err

}

/*
// asn1Data starting with -----BEGIN CERTIFICATE-----
func convertPemsToDers(asn1Data string) ([][]byte, error) {

	asn1str, err := hex.DecodeString(asn1Data)

	if err != nil {
		fmt.Println(err)
		return [][]byte{}, err
	}

	certs, err := x509.ParseCertificates(asn1str)

	if err != nil {
		fmt.Println(err)
		return [][]byte{}, err
	}

	genMap := make([][]byte, len(certs))

	for idx, cert := range certs {
		fmt.Printf("%v\n", cert)
		genMap[idx] = append(genMap[idx][0], cert.Raw)
	}

}
*/

func listPods(clientset *kubernetes.Clientset) (int, error) {
	// get PodList
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		return 0, err
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// get []Pod and list pod names
	for _, pod := range pods.Items {
		//	p, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
		fmt.Fprintln(os.Stdout, pod.GetName())
	}
	return len(pods.Items), err
}

func getPodObject() *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-test-pod",
			Namespace: "default",
			Labels: map[string]string{
				"app": "demo",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:            "busybox",
					Image:           "busybox",
					ImagePullPolicy: v1.PullIfNotPresent,
					Command: []string{
						"sleep",
						"3600",
					},
				},
			},
		},
	}
}

func ReadBinaryFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		return nil, statsErr
	}

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)

	return bytes, err
}

func WriteBinaryFile(filename string, data []byte) error {

	file, err := os.Open(filename)
	if err == nil {
		file.Close()
		return errors.New("output file already exists")
	}

	err = ioutil.WriteFile(filename, data, 0644)
	return err

}
