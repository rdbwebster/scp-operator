package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	api "github.com/rdbwebster/scp-operator/api/v1"
	clientV1alpha1 "github.com/rdbwebster/scp-operator/clientset/v1alpha1"
	"github.com/rdbwebster/scp-operator/kclient"
	"github.com/rdbwebster/scp-operator/model"
	"github.com/rdbwebster/scp-operator/stacktrace"
	v1 "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var clusterCurrentId int
var serviceCurrentId int
var factoryCurrentId int
var userCurrentId int = 60

var scpclusterRes = schema.GroupVersionResource{Group: "webapp.my.domain", Version: "v1", Resource: "scpclusters"}

// datastore
var clusterInfos model.ClusterInfos
var serviceInfos model.ServiceInfos
var groupInfos model.GroupInfos
var factoryInfos model.FactoryInfos

var userInfos model.UserInfos

var clusterClient *clientV1alpha1.ExampleV1Alpha1Client

var dynamicClient dynamic.Interface
var clientset *kubernetes.Clientset

var pemData = `-----BEGIN CERTIFICATE-----
MIIDADCCAeigAwIBAgIBAjANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwptaW5p
a3ViZUNBMB4XDTIwMDQyODIwMDUwNloXDTIxMDQyOTIwMDUwNlowMTEXMBUGA1UE
ChMOc3lzdGVtOm1hc3RlcnMxFjAUBgNVBAMTDW1pbmlrdWJlLXVzZXIwggEiMA0G
CSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDy1XxsZOcvpUtFapyuk3zkRShHxhS5
PjLVBse8NRPf1H6F2Ro+HpCzJbIx98JtcYJi5Rk6R9dBAg/r8UHbWDrGtQkNLpfq
4HK5zOqBMxqOFVIO8Vlo008nqQ/5WZkrFsVAR2dN//6plCIYQY+Jb48yc0bZbZFy
mS886YT3MR6sExqYCcnjXfVCuRME7n1ZEBM1KzsbIGh0AMqFVxhpHx3Iu6woX7OA
CDO5zuzPsIkzhp3r77aPQvV91oaHU25TfETG9dy02jJAPSB+ON2lzluMECnTY8cQ
ZjIjJBC6TqVOJ7D6zjZN1J4v9iEjpsEgEGAOIiXMHyrPFEsQS5MHFOPFAgMBAAGj
PzA9MA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAQYIKwYBBQUH
AwIwDAYDVR0TAQH/BAIwADANBgkqhkiG9w0BAQsFAAOCAQEAbo8I8V8wDnKCIfN3
BN09GwCyMxx1nrauMoxoYS2/R4MlHj6+e3Uw7aJjZdMR+WSu2JHOJaGrXMdvejsp
34J8wY1d8vHjoPr5KZFb9n8u7IekYot6dDZEeRE2u3URWjZk2a2QqsBC+Fo7jELM
Tb8KLmlVE5MusAzhTlIgMchSWw3w7SuNL1G4GGolws7rpYURQWDfUXxOcpwdsllv
2IMFVMXDhIHACFWgnGL+6scwo1K+/2n2jO+IaZNm5ydiOXXaIBb1jDdOWmuamJcZ
9g4XLQqQkl1NZQdeWNTUltmq68hu4ONLy/TPdHJPktgvH+wgaLJKzmaUkSQSZuX5
sS41JA==
-----END CERTIFICATE-----`

var certAuth = `-----BEGIN CERTIFICATE-----
MIIC5zCCAc+gAwIBAgIBATANBgkqhkiG9w0BAQsFADAVMRMwEQYDVQQDEwptaW5p
a3ViZUNBMB4XDTIwMDQwOTE1MjcyN1oXDTMwMDQwNzE1MjcyN1owFTETMBEGA1UE
AxMKbWluaWt1YmVDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBANDo
ZULypJVG2/JuTNh3mHOWnBjTpxFms9I6FVTuJrulD6gDOfrMDrzYcCkHA/7Pi8s4
CVfy7tIUaiACYc9yR4DYgev9GBQCr9voBLbMwutXVK4j/g/D4Li+ZOAo/Drqf6DF
5bAMm+N9Fy87C7DMtOEt6rrgTUhBLT6bunwG40GOxrpoOpt3R7TFRURKin3QMQBe
y82gHsyacymGVLdJH0wUp6NA5GLFjlIo9CKqhxJdd2V8YqtD/0AyigJZ7ygB+xML
TWnXePuvF73oPohjv3g6znT85rWR291b09RJ0eyg0nUSlU2ruNG+Q4INHvAPozfE
YZ/vKr9zZHQouSR/HWsCAwEAAaNCMEAwDgYDVR0PAQH/BAQDAgKkMB0GA1UdJQQW
MBQGCCsGAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3
DQEBCwUAA4IBAQBjlB1+MXREuamTda/l7W3pffqKrEZAHqb+h/bv6QWdRS7q371O
RKJa/SjEICZiwASKj7ESSK00oblYknJ1xYs7tPZVHzmsjYRy1b464vHiR1vOBebh
2/9lsxB14GHXpostNchOGhaS4dUpbmvAxv19ePVd2BqlAEm/eGHVgvkCCv5MPuFr
PodsTfC/zQIExhOsTaUaTozO5JWBdDoEzhf49r6Qxo5dD6FCVSXuQuiGMP5W7Gca
h9dE0rvLTu8nqncDoyTtek8h2nKxS6Da0soyeFkik4J5kMooUskGVOunT8RTB4EH
ejqdItTXlDY7JlO4nhpOANulM3cSl8LtLM26
-----END CERTIFICATE-----`

var bearerToken = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA8tV8bGTnL6VLRWqcrpN85EUoR8YUuT4y1QbHvDUT39R+hdka
Ph6QsyWyMffCbXGCYuUZOkfXQQIP6/FB21g6xrUJDS6X6uByuczqgTMajhVSDvFZ
aNNPJ6kP+VmZKxbFQEdnTf/+qZQiGEGPiW+PMnNG2W2RcpkvPOmE9zEerBMamAnJ
4131QrkTBO59WRATNSs7GyBodADKhVcYaR8dyLusKF+zgAgzuc7sz7CJM4ad6++2
j0L1fdaGh1NuU3xExvXctNoyQD0gfjjdpc5bjBAp02PHEGYyIyQQuk6lTiew+s42
TdSeL/YhI6bBIBBgDiIlzB8qzxRLEEuTBxTjxQIDAQABAoIBACbKbZ7PG5Mj33tO
RYspqki4t2+Ht+XDhtE6zQtGm08lHbT58lQ8A7dqbSXIQznCaSatHDOQKFWNI8f3
+SI13OXDI5gEemYdxpXhoxBSfop1427ZpQO2xa07N38IjgwxDf9vqJPwMdka/btM
NcapmIFhos62zwY0bDRZxVDLeu/Xqbpktujo5bixN2No4xhCHPoMz2nc2eiuHnpQ
Qpdqr11nSrAGj7EASK3P0gcfj0DH5g28n2FX90VWbIf30QSZvXPf/Jcch8NaRkBi
lUiIxNDalyYUPyH9Wq20cLtJ2zt/XxziZiFTHVcrTJ+tN4NLeW6s8KMbbJkxhkbq
0Qn1vn0CgYEA/lPRNz3Lk+aKnQjlZPfR/Zqzqspy0a7N+WCPk6InRiFafF1v2wO4
pN7vAiOX+6424uMOtMghLkUEVJA3Gr0zf34dh89/imw4DmUCNs76QxvJdAVbcEMn
JJVbwMQgivTAQkYgnUeqiVY6jnd8HkEOVk9b/9lzTjHywtj2ZGRK3A8CgYEA9G5R
iDbmzTVgOJeAfNRvq2yAFie8rp+jRevb3bJfL84lcgdOPeEIrRnnFdueU6gBuqcA
OxKJeA9IWyKBqSNvCdxPrRJSDRdvq/EKf1gvMwdDs87H7J7Sg3cpEBJ7gkQ2fJ0E
9+76EEuYahS1KJ4JSGo6+kxEI4NH771ud8Ga/usCgYEA1vBJYbFlCsMNJLgu/oz+
uKD09QOR1Coyw25bCT9Ch9+KVI63CNb1RsluH1WrjbXnhwq0FA8LE8qaZUlYeM2r
5zTTikLQHFFncqrlGyMDmJG0SMx6Qb9PJnjgCWL3ydgdYCVaTPITa2wnot3SVNNQ
ZZs+OlUxQMWv0AKDAcdNCPsCgYEAiZVjDTIh/dYSgCg+6YTGCo67Fj1txjkTNTNK
geJ6E7WMfD/CebAmKxFOco4480u5FXAVACsx98NabfnhU+we/0TkED4kszvC3tyB
lSZ1AtsO77Hv9K99PQSgt2w/2xY8OS5E8q2wUeXLN8LKKb+y5/Drm6G8JOUrY7WT
7ZKrhNsCgYBjdIbWYdvmWeMTRG9VmIBEOz60MNPJeOF025v0HNacOnD431KQkmQA
hC5MOOdmGl/A328EhZVqRmwuSaBXNrTNSLj6KFqUZMaoawK57AD2H2NRAwojjyP2
Zjfex97rrbjORBea5+zZaWpTOsvjG/fqWLYJHQC3D3CMzEF0cNJ17g==
-----END RSA PRIVATE KEY-----`

// Give us some seed data
func init() {

	var config *rest.Config
	var err error

	if kubeconfig == "" {
		log.Printf("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		log.Printf("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		panic(err)
	}

	//	var client *clientV1alpha1.ExampleV1Alpha1Client
	clusterClient, err = clientV1alpha1.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	dynamicClient, err = dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	/*	https://pkg.go.dev/k8s.io/client-go/kubernetes */
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Static mock data
	RepoCreateCluster(model.ClusterInfo{Spec: api.SCPclusterSpec{Clustername: "Cluster One",
		Namespace: "default",
		Url:       "192.168.64.4:8443",
		Token:     bearerToken,
		Cert:      pemData,
		CertAuth:  certAuth}})
	//	RepoCreateCluster(model.ClusterInfo{Name: "Cluster Two", NAmespace: "default", Url: "192.168.0.42", Token: "", Cert: "", CertAuth: ""})

	RepoCreateMockService(model.ServiceInfo{Name: "Postgres db1", Crdname: "postgres.sql.tanzu.vmware.com", Clustername: "", Status: "Active"})
	RepoCreateMockService(model.ServiceInfo{Name: "Postgres db2", Crdname: "postgres.sql.tanzu.vmware.com", Clustername: "cluster1", Status: "Active"})

	//	RepoCreateFactory(model.FactoryInfo{Name: "etcd", Version: "1", Deploymentname: "etcd", Clustername: "LOCAL"})

	//var members1 = []string {""}
	//var members2 = []string {""}
	//var members3 = []string {""}
	//RepoCreateGroup(modelGroupInfo{Name: "platform_operators", Member: members1 })
	//RepoCreateGroup(modelGroupInfo{Name: "service_operators", Member: members2 })
	//RepoCreateGroup(modelGroupInfo{Name: "developers", Member: members3 })

	user1 := model.UserInfo{FirstName: "Patrick", LastName: "Star", Email: "developer@vmware.com", Password: "VMware1!", Id: ""}
	user1.Roles[0] = "DEVELOPER"
	RepoCreateUser(user1)

	user2 := model.UserInfo{FirstName: "Sandy", LastName: "Cheeks", Email: "admin@vmware.com", Password: "VMware1!", Id: ""}
	user2.Roles[0] = "ADMINISTRATOR"
	RepoCreateUser(user2)

	user3 := model.UserInfo{FirstName: "Eugene", LastName: "Crabs", Email: "platformop@vmware.com", Password: "VMware1!", Id: ""}
	user3.Roles[0] = "PLATFORM_OPERATOR"
	RepoCreateUser(user3)

	user4 := model.UserInfo{FirstName: "Squidward", LastName: "Tentacles", Email: "serviceop@vmware.com", Password: "VMware1!", Id: ""}
	user4.Roles[0] = "SERVICE_OPERATOR"
	RepoCreateUser(user4)

}

func RepoCreateUser(t model.UserInfo) model.UserInfo {
	userCurrentId++
	t.Id = strconv.Itoa(userCurrentId)
	userInfos = append(userInfos, t)
	return t
}

func RepoGetUserInfos() []model.UserInfo {
	return userInfos
}

//
//  Cluster Repos Methods
//

func RepoGetClusters() error {

	// *v1.SCPclusterList
	clusterList, err := clusterClient.SCPcluster("default").List(metav1.ListOptions{})
	if err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		fmt.Printf("Error retrieving clusters %+v \n", st)
		return err
	}

	// replace the cached clusterInfos
	clusterInfos = nil
	for _, t := range clusterList.Items {
		clusterInfos = append(clusterInfos,

			model.ClusterInfo{Spec: api.SCPclusterSpec{Clustername: t.Spec.Clustername, Namespace: t.Spec.Namespace, Url: t.Spec.Url,
				Token: t.Spec.Token, Cert: t.Spec.Cert, CertAuth: t.Spec.CertAuth}})
	}
	return nil

}

func RepoFindCluster(clustername string) model.ClusterInfo {
	for _, t := range clusterInfos {
		if t.Spec.Clustername == clustername {
			return t
		}
	}
	// return empty Todo if not found

	return model.ClusterInfo{Spec: api.SCPclusterSpec{Clustername: ""}}
}

func RepoCreateCluster(t model.ClusterInfo) model.ClusterInfo {

	scpcluster := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "webapp.my.domain/v1",
			"kind":       "SCPcluster",
			"metadata": map[string]interface{}{
				"name": t.Spec.Clustername,
			},
			"spec": map[string]interface{}{
				"clustername": t.Spec.Clustername,
				"namespace":   t.Spec.Namespace,
				"url":         t.Spec.Url,
				"token":       t.Spec.Token,
				"cert":        t.Spec.Cert,
				"auth":        t.Spec.CertAuth,
				"lastcontact": metav1.Now(),
			},
		},
	}

	fmt.Println("Creating cluster...")
	result, err := dynamicClient.Resource(scpclusterRes).Namespace("default").Create(context.TODO(), scpcluster, metav1.CreateOptions{})
	if err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		fmt.Printf("Error creating clusters %+v \n", st)
	} else {
		fmt.Printf("Created cluster %q.\n", result.GetName())
	}

	//clusterCurrentId++
	//t.Id = clusterCurrentId
	//clusterInfos = append(clusterInfos, t)
	return t
}

func RepoUpdateCluster(ci model.ClusterInfo) model.ClusterInfo {

	for _, t := range clusterInfos {
		if t.Spec.Clustername == ci.Spec.Clustername {
			t.Spec.Clustername = ci.Spec.Clustername
			t.Spec.Namespace = ci.Spec.Namespace
			t.Spec.Url = ci.Spec.Url
			t.Spec.Token = ci.Spec.Token
			t.Spec.Cert = ci.Spec.Cert
			t.Spec.CertAuth = ci.Spec.CertAuth
			t.Spec.Lastcontact = ci.Spec.Lastcontact
		}
	}
	return ci
}

func RepoDeleteCluster(clustername string) error {

	c := RepoFindCluster(clustername)
	if c.Spec.Clustername == "" {
		return nil
	}

	fmt.Printf("Deleting cluster... %s \n", clustername)
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	err := dynamicClient.Resource(scpclusterRes).Namespace(c.Spec.Namespace).Delete(context.TODO(), c.Spec.Clustername, deleteOptions)
	if err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		fmt.Printf("Error creating clusters %+v \n", st)
	} else {
		fmt.Printf("Deleted cluster %q.\n", clustername)
	}

	//	for i, t := range clusterInfos {
	//		if t.Id == id {
	//			clusterInfos = append(clusterInfos[:i], clusterInfos[i+1:]...)
	//			return nil
	//		}
	//	}
	//	return fmt.Errorf("Could not find Cluster with id of %d to delete", id)
	return nil
}

//
//  Service Repos Methods
//

func RepoGetServices() error {

	// reset cached service infos
	serviceInfos = nil

	// replace mock entries
	RepoCreateMockService(model.ServiceInfo{Name: "Postgres db1", Crdname: "postgres.sql.tanzu.vmware.com", Clustername: "LOCAL", Status: "Active"})
	RepoCreateMockService(model.ServiceInfo{Name: "Postgres db2", Crdname: "postgres.sql.tanzu.vmware.com", Clustername: "cluster1", Status: "Active"})

	// Get list of Managed Operators
	// *v1.ManagedOperatorList
	factoryList, err := clusterClient.ManagedOperator("default").List(metav1.ListOptions{})
	if err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		fmt.Printf("Error retrieving factories %+v \n", st)
		return err
	}

	// get instances of each managed operator
	for _, f := range factoryList.Items {

		// get the crd
		var crdresource = kclient.GetCustomResourceDefinition(dynamicClient, f.Spec.CrdName)

		// Get gvr for CRD Type
		crdgvr := kclient.GetCrdGrv(dynamicClient, crdresource)

		// get a list of those crs by crd type
		// *unstructured.UnstructuredList
		crinstances := kclient.ListCustomResourceInstances(dynamicClient, crdgvr, "default")

		// for each cr instance
		for _, cri := range crinstances.Items {

			fmt.Printf("cri %+v \n", cri)

			// Parse out CR name
			name, found, err := unstructured.NestedString(cri.Object, "metadata", "name")
			if found {
				fmt.Printf("\nSee Name: %s\n", name)
			} else if err != nil {
				st := stacktrace.New(err.Error())
				log.Printf("%s\n", st)
				fmt.Printf("Error retrieving services %+v \n", st)
			}

			// create new svc info
			newsvcinfo := model.ServiceInfo{Name: name, Crdname: "localhost",
				Clustername: "LOCAL", Status: ""}

			// copy in croutputs from managed op spec to new svcinfo
			entries := make([]api.CRentry, 0)

			fmt.Printf("length of croutputs  %d \n", len(f.Spec.CRoutputs))

			for _, oentry := range f.Spec.CRoutputs {
				newentry := api.CRentry{ControlName: oentry.ControlName, ControlType: oentry.ControlType,
					CurrentValue: "", ValueType: oentry.ValueType, CRpath: oentry.CRpath}
				entries = append(entries, newentry)
			}
			newsvcinfo.CRoutputs = entries

			fmt.Printf("ready with newsvcinfo %+v \n", newsvcinfo)

			// get interesting fields from cr, updating output values in servinceinfo
			kclient.GetResourceOutputs(dynamicClient, crdgvr, name,
				"default", newsvcinfo.CRoutputs)

			fmt.Printf("\nCRoutputs = %+v\n", newsvcinfo.CRoutputs)

			// add new entry to svcinfos cache
			serviceInfos = append(serviceInfos, newsvcinfo)
		}
	}

	return nil
}

/*
	// get service for each managed operator using app label
	for _, f := range factoryList.Items {
		services, err := GetServicesByLabel("default", f.Spec.ServiceLabel)
		if err != nil {
			st := stacktrace.New(err.Error())
			log.Printf("%s\n", st)
			fmt.Printf("Error retrieving services %+v \n", st)
			return err
			}
			for _, s := range services.Items {
				if !contains(serviceInfos, s.Name) {
				serviceInfos = append(serviceInfos,
							model.ServiceInfo{Name: s.Name, Url: "", Clustername: "LOCAL", Status: "Available"})
			}

			}
	}*/

func contains(s []model.ServiceInfo, name string) bool {
	for _, v := range s {
		if v.Name == name {
			return true
		}
	}

	return false
}

func RepoFindService(name string) model.ServiceInfo {
	for _, t := range serviceInfos {
		if t.Name == name {
			return t
		}
	}
	// return empty Todo if not found
	return model.ServiceInfo{Name: ""}
}

func RepoCreateMockService(serv model.ServiceInfo) {
	serviceInfos = append(serviceInfos, serv)
}

func RepoCreateService(serv model.ServiceInfo) model.ServiceInfo {

	// Read Custom Resource Definition by name
	// *unstructured.Unstructured
	var crdresource = kclient.GetCustomResourceDefinition(dynamicClient, serv.Crdname)

	// Get gvr for CR Type
	gvr := kclient.GetCrdGrv(dynamicClient, crdresource)

	// Parse out CRD kind
	kind, found, err := unstructured.NestedString(crdresource.Object, "spec", "names", "kind")
	if found {
		fmt.Printf("See Kind : %s\n", kind)
	} else {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		fmt.Printf("Error creating clusters %+v \n", st)
	}
	fmt.Printf("See crdresource.Object type of: %T\n", crdresource.Object)

	// Create a custom new resource with these values
	kclient.CreateCustomResourceSetParms(dynamicClient, serv.Name, gvr, kind, serv)

	serviceInfos = append(serviceInfos, serv)
	return serv
}

func RepoUpdateService(ci model.ServiceInfo) model.ServiceInfo {

	for _, t := range serviceInfos {
		if t.Name == ci.Name {
			t.Name = ci.Name
			t.Crdname = ci.Crdname
			t.Status = ci.Status
		}
	}
	return ci
}

func RepoDeleteService(name string) error {
	for i, t := range serviceInfos {
		if t.Name == name {
			serviceInfos = append(serviceInfos[:i], serviceInfos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Service with name of %s to delete", name)
}

//
//  Factory Repos Methods
//

func RepoGetFactories() error {

	// *v1.ManagedOperatorList
	factoryList, err := clusterClient.ManagedOperator("default").List(metav1.ListOptions{})
	if err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		fmt.Printf("Error retrieving factories %+v \n", st)
		return err
	}

	// replace the cached
	factoryInfos = nil

	// add factory to available list only if it has a deployment
	for _, f := range factoryList.Items {
		//	deployments, err := GetDeploymentsByField("default", "metadata.name=" + f.Name)
		//	if err != nil {
		//		st := stacktrace.New(err.Error())
		//		log.Printf("%s\n", st)
		//		fmt.Printf("Error retrieving factories %+v \n", st)
		//		return err
		//	 }
		//	 if len(deployments.Items) > 0 {
		//	factoryInfos = append(factoryInfos, f.Spec)
		factoryInfos = append(factoryInfos,
			model.FactoryInfo{Spec: f.Spec, Clustername: "LOCAL"})
		//	}

	}
	return nil

}

func GetServicesByLabel(namespace string, labelSelector string) (*core.ServiceList, error) {

	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return nil, err
	}
	for _, d := range services.Items {
		fmt.Printf("Service  %s\n", d.Name)
	}

	// https://pkg.go.dev/k8s.io/api/apps/v1
	return services, nil
}

func GetDeploymentsByLabel(namespace string, labelSelector string) (*v1.DeploymentList, error) {
	// "app=scp-spa"
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	fmt.Printf("Type %T \n", deployments)
	if err != nil {
		return nil, err
	}
	for _, d := range deployments.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	// https://pkg.go.dev/k8s.io/api/apps/v1
	return deployments, nil
}

func GetDeploymentsByField(namespace string, fieldSelector string) (*v1.DeploymentList, error) {
	// "app=scp-spa"
	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{FieldSelector: fieldSelector})
	fmt.Printf("Type %T \n", deployments)
	if err != nil {
		return nil, err
	}
	for _, d := range deployments.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}

	// https://pkg.go.dev/k8s.io/api/apps/v1
	return deployments, nil
}

func RepoFindFactory(name string) model.FactoryInfo {
	for _, t := range factoryInfos {
		fmt.Printf("Compare '%s' and '%s' \n", t.Spec.Name, name)
		if t.Spec.Name == name {
			fmt.Printf("Found")
			return t
		}
	}
	// return empty factory if not found
	return model.FactoryInfo{Spec: api.ManagedOperatorSpec{}, Clustername: ""}
}

//this is bad, I don't think it passes race condtions
func RepoCreateFactory(t model.FactoryInfo) model.FactoryInfo {

	factoryInfos = append(factoryInfos, t)
	return t
}

func RepoUpdateFactory(ci model.FactoryInfo) model.FactoryInfo {

	for _, t := range factoryInfos {
		if t.Spec.Name == ci.Spec.Name {
			t.Spec.Name = ci.Spec.Name
			// TODO update more than name
		}
	}
	return ci
}

func RepoDeleteFactory(name string) error {
	for i, t := range factoryInfos {
		if t.Spec.Name == name {
			factoryInfos = append(factoryInfos[:i], factoryInfos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find factory with name of %s to delete", name)
}
