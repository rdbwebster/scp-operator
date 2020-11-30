package main

import (
	"fmt"
	"strconv"

	"github.com/rdbwebster/scp-rest-svr/model"
)

var clusterCurrentId int
var serviceCurrentId int
var factoryCurrentId int
var userCurrentId int = 60

// datastore
var clusterInfos model.ClusterInfos
var serviceInfos model.ServiceInfos
var factoryInfos model.FactoryInfos
var userInfos model.UserInfos

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
	RepoCreateCluster(model.ClusterInfo{Name: "Cluster One", Url: "192.168.64.4:8443", Token: bearerToken, Cert: pemData, CertAuth: certAuth})
	RepoCreateCluster(model.ClusterInfo{Name: "Cluster Two", Url: "192.168.0.42", Token: "", Cert: "", CertAuth: ""})
	RepoCreateService(model.ServiceInfo{Name: "Postgres db1", Url: "192.168.0.42", Status: "Active"})
	RepoCreateService(model.ServiceInfo{Name: "Postgres db2", Url: "192.168.64.4:8443", Status: "Active"})
	RepoCreateFactory(model.FactoryInfo{Name: "Tanzu Postgres Operator", Url: "http://192.168.0.42", Status: "Active"})

	user1 := model.UserInfo{FirstName: "Patrick", LastName: "Star", Email: "developer@vmware.com", Password: "VMware1!", Id: ""}
	user1.Roles[0] = "TENANT_USER"
	RepoCreateUser(user1)

	user2 := model.UserInfo{FirstName: "Sandy", LastName: "Cheeks", Email: "admin@vmware.com", Password: "VMware1!", Id: ""}
	user2.Roles[0] = "TENANT_ADMIN"
	RepoCreateUser(user2)

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

func RepoFindCluster(id int) model.ClusterInfo {
	for _, t := range clusterInfos {
		if t.Id == id {
			return t
		}
	}
	// return empty Todo if not found

	return model.ClusterInfo{Id: 0}
}

//this is bad, I don't think it passes race condtions
func RepoCreateCluster(t model.ClusterInfo) model.ClusterInfo {
	clusterCurrentId++
	t.Id = clusterCurrentId
	clusterInfos = append(clusterInfos, t)
	return t
}

func RepoUpdateCluster(ci model.ClusterInfo) model.ClusterInfo {

	for _, t := range clusterInfos {
		if t.Id == ci.Id {
			t.Name = ci.Name
			t.Url = ci.Url
			t.Token = ci.Token
			t.Cert = ci.Cert
			t.CertAuth = ci.CertAuth
			t.Connected = ci.Connected
		}
	}
	return ci
}

func RepoDeleteCluster(id int) error {
	for i, t := range clusterInfos {
		if t.Id == id {
			clusterInfos = append(clusterInfos[:i], clusterInfos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Cluster with id of %d to delete", id)
}

//
//  Service Repos Methods
//

func RepoFindService(id int) model.ServiceInfo {
	for _, t := range serviceInfos {
		if t.Id == id {
			return t
		}
	}
	// return empty Todo if not found
	return model.ServiceInfo{Id: 0}
}

//this is bad, I don't think it passes race condtions
func RepoCreateService(t model.ServiceInfo) model.ServiceInfo {
	serviceCurrentId++
	t.Id = serviceCurrentId
	serviceInfos = append(serviceInfos, t)
	return t
}

func RepoUpdateService(ci model.ServiceInfo) model.ServiceInfo {

	for _, t := range serviceInfos {
		if t.Id == ci.Id {
			t.Name = ci.Name
			t.Url = ci.Url
			t.Status = ci.Status
		}
	}
	return ci
}

func RepoDeleteService(id int) error {
	for i, t := range serviceInfos {
		if t.Id == id {
			serviceInfos = append(serviceInfos[:i], serviceInfos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Service with id of %d to delete", id)
}

//
//  Factory Repos Methods
//

func RepoFindFactory(id int) model.FactoryInfo {
	for _, t := range factoryInfos {
		if t.Id == id {
			return t
		}
	}
	// return empty Todo if not found
	return model.FactoryInfo{Id: 0}
}

//this is bad, I don't think it passes race condtions
func RepoCreateFactory(t model.FactoryInfo) model.FactoryInfo {
	factoryCurrentId++
	t.Id = factoryCurrentId
	factoryInfos = append(factoryInfos, t)
	return t
}

func RepoUpdateFactory(ci model.FactoryInfo) model.FactoryInfo {

	for _, t := range factoryInfos {
		if t.Id == ci.Id {
			t.Name = ci.Name
			t.Url = ci.Url
			t.Status = ci.Status
		}
	}
	return ci
}

func RepoDeleteFactory(id int) error {
	for i, t := range factoryInfos {
		if t.Id == id {
			factoryInfos = append(factoryInfos[:i], factoryInfos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find factory with id of %d to delete", id)
}
