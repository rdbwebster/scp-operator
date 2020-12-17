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
	"strconv"
	"strings"

	api "github.com/rdbwebster/scp-operator/api/v1"
	"github.com/rdbwebster/scp-operator/model"
	"github.com/rdbwebster/scp-operator/stacktrace"
	"gopkg.in/yaml.v2"
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

func GetResource(client dynamic.Interface, gvr schema.GroupVersionResource,
	name string, namespace string) *unstructured.Unstructured {

	var err error
	var result *unstructured.Unstructured

	if namespace == "" {
		result, err = client.Resource(gvr).Get(context.TODO(), name, metav1.GetOptions{})
	} else {
		//	result, getErr := client.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
		result, err = client.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	}

	if err != nil {
		panic(fmt.Errorf("failed to get latest version of Custom Resource:name= %s, namespace=%s, gvr= %+v, err= %v",
			name, namespace, gvr, err))
	}

	fmt.Printf("Resource= \n %#v", result)

	return result
}

func GetResourceOutputs(client dynamic.Interface, gvr schema.GroupVersionResource, name string,
	namespace string, croutputs []api.CRentry) {

	//	unstructuredresource := GetCustomResourceDefinition(client, name)

	unstructuredresource := GetResource(client, gvr, name, namespace)
	fmt.Printf("\nRetrieved Resource: %+v\n", unstructuredresource)

	for index, entry := range croutputs {

		if strings.HasPrefix(entry.CRpath, "status.") {

			if entry.ValueType == "text" {
				// Parse out CRD text field
				resultvalue, found, err := unstructured.NestedString(unstructuredresource.Object, "status",
					strings.TrimPrefix(entry.CRpath, "status."))
				if found {
					fmt.Printf("\nSee Output Value start: %s\n", resultvalue)
					croutputs[index].CurrentValue = resultvalue
					fmt.Printf("\nSee Output Value end: %s\n", croutputs[index].CurrentValue)
				} else {
					fmt.Printf("\nCannot find Output Value %s, error = %+v", resultvalue, err)
				}

			} else if entry.ValueType == "number" {
				// Parse out CRD number field

				resultvalue, found, err := unstructured.NestedInt64(unstructuredresource.Object, "status",
					strings.TrimPrefix(entry.CRpath, "status."))
				if found {
					fmt.Printf("\nSee Output Value: %d\n", resultvalue)
					croutputs[index].CurrentValue = strconv.FormatInt(resultvalue, 10)
				} else {
					fmt.Printf("\nCannot find Output Value %d, error = %+v", resultvalue, err)
				}

			}

		}
	}

}

func ListCustomResourceInstances(client dynamic.Interface, gvr schema.GroupVersionResource, namespace string) *unstructured.UnstructuredList {

	// other ways to specifiy fields https://github.com/kubernetes/apimachinery/issues/47
	// deployments, err := clientset.AppsV1().Deployments("default").List(context.TODO(), metav1.ListOptions{FieldSelector: "metadata.name=scp-rest-svr"})

	results, err := client.Resource(gvr).Namespace(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d CRs of type %s\n", len(results.Items), gvr.Resource)
	fmt.Printf("See type of: %T\n", results)

	return results
}

func Between(value string, a string, b string) string {
	// Get substring between two strings.
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func Before(value string, a string) string {
	// Get substring before a string.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	return value[0:pos]
}

func After(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func GetGrv(client dynamic.Interface, crdresource *unstructured.Unstructured) schema.GroupVersionResource {

	// Parse out group
	group, found, err := unstructured.NestedString(crdresource.Object, "apiVersion")
	group = Before(group, "/")
	if found {
		fmt.Printf("See Group: %s\n", group)
	} else {
		panic(err)
	}

	// Parse out CRD kind plural
	plural, found, err := unstructured.NestedString(crdresource.Object, "spec", "names", "plural")
	if found {
		fmt.Printf("See Kind Plural: %s\n", plural)
	} else {
		panic(err)
	}

	// Parse out api/version
	version, found, err := unstructured.NestedString(crdresource.Object, "apiVersion")
	version = After(version, "/")
	if found {
		fmt.Printf("See Version: %s\n", version)
	} else {
		panic(err)
	}

	return schema.GroupVersionResource{Group: group, Version: version, Resource: plural}

}

func GetCrdGrv(client dynamic.Interface, crdresource *unstructured.Unstructured) schema.GroupVersionResource {

	// Parse out to confirm CRD name
	name, found, err := unstructured.NestedString(crdresource.Object, "metadata", "name")
	if found {
		fmt.Printf("\nSee cr name in crd: %s\n", name)
	}

	// Parse out CRD group
	group, found, err := unstructured.NestedString(crdresource.Object, "spec", "group")
	if found {
		fmt.Printf("See cr Group: %s\n", group)
	}

	// Parse out CRD kind plural
	plural, found, err := unstructured.NestedString(crdresource.Object, "spec", "names", "plural")
	if found {
		fmt.Printf("See Kind Plural: %s\n", plural)
	}

	// Parse out CRD version
	// https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1/unstructured
	versions, found, err := unstructured.NestedSlice(crdresource.Object, "spec", "versions")
	if found {
		fmt.Printf("See Versions: %s of type %T\n", versions, versions)
	} else {
		fmt.Printf("Cannot find version %+v", err)
	}

	// Must cast
	// https://github.com/kubernetes/client-go/blob/master/examples/dynamic-create-update-delete-deployment/main.go
	foo := versions[0].(map[string]interface{})
	fmt.Printf("See foo type of: %T\n", foo)

	version, found, err := unstructured.NestedString(foo, "name")
	if found {
		fmt.Printf("See Version: %s of type %T\n", version, version)
	} else {
		fmt.Printf("Cannot find name %+v", err)
	}

	return schema.GroupVersionResource{Group: group, Version: version, Resource: plural}

}

func CreateCustomResourceSetParms(client dynamic.Interface, name string, gvr schema.GroupVersionResource,
	kind string, values model.ServiceInfo) {

	//	crGvr := schema.GroupVersionResource{Group: group, Version: version, Resource: plural}

	gv := gvr.Group + "/" + gvr.Version

	cr := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": gv,
			"kind":       kind,
			"metadata": map[string]interface{}{
				"name": name,
			},
			"spec": map[string]interface{}{
				//		"size": 2,
				//		"version": "3.1.10",
			},
		},
	}

	// add entries
	for _, entry := range values.CRinputs {

		if strings.HasPrefix(entry.CRpath, "spec.") {
			if entry.ValueType == "text" {

				err := unstructured.SetNestedField(cr.Object, entry.CurrentValue, "spec",
					strings.TrimPrefix(entry.CRpath, "spec."))
				if err != nil {
					fmt.Printf("\n Cannot set Value: %s, error = %+v \n", entry.CRpath, err)
				}
			} else if entry.ValueType == "number" {
				// Parse out CRD number field

				value, err := strconv.ParseInt(entry.CurrentValue, 10, 64)
				if err != nil {
					panic(err)
				}
				err = unstructured.SetNestedField(cr.Object, value, "spec",
					strings.TrimPrefix(entry.CRpath, "spec."))
				if err != nil {
					fmt.Printf("\n Cannot set Value %s, error = %+v \n", entry.CRpath, err)
				}
			}
		}
	}

	// set the spec updated spec field
	//	if err := unstructured.SetNestedField(result.Object, containers, "spec", "template", "spec", "containers"); err != nil {
	//		panic(err)
	//	}

	// Create Deployment
	fmt.Printf("Creating %s resource...", kind)
	result, err := client.Resource(gvr).Namespace("default").Create(context.TODO(), cr, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetName())
	fmt.Printf("Unstructured \n %+v \n", cr.Object)
	y, err := yaml.Marshal(cr.Object)
	if err != nil {
		fmt.Printf("Marshal: %v", err)
	}

	fmt.Println(string(y))

}

func GetCustomResourceDefinition(client dynamic.Interface, name string) *unstructured.Unstructured {

	gvr := schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}

	//	result, getErr := client.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	result, getErr := client.Resource(gvr).Get(context.TODO(), name, metav1.GetOptions{})

	if getErr != nil {
		panic(fmt.Errorf("failed to get latest version of Custom Resource Definition: %v", getErr))
	}

	return result
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
