package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/rdbwebster/scp-rest-svr/kclient"
	"github.com/rdbwebster/scp-rest-svr/model"
	"github.com/rdbwebster/scp-rest-svr/stacktrace"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

//
// currentuser

func GetSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, " {\"forcePasswordChange\": false}")
}

//
// Login
//

func Login(w http.ResponseWriter, r *http.Request) {
	var req SessionRequest
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))

	if err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(err)
		}
	}

	if err := r.Body.Close(); err != nil {
		log.Print(err)
	}
	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(err)
		}
	}

	var found bool

	for _, v := range RepoGetUserInfos() {

		if v.Email == req.Email && v.Password == req.Password {
			found = true
			v.Password = "####"
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Header().Set("Authorization", v.Id)
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(v); err != nil {
				log.Print(err)
			}
			fmt.Println("\033[33m", "Login Successful")
		}
	}

	if !found {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusForbidden, Text: "Invalid Login"}); err != nil {
			log.Print(err)
		}
		fmt.Println("\033[33m", "Login Unsuccessful")
	}
}

//
// Cluster Handlers
//

func GetClusters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(clusterInfos); err != nil {
		log.Print(err)
	}
}

func GetCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id int
	var err error
	if id, err = strconv.Atoi(vars["id"]); err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	clusterInfo := RepoFindCluster(id)
	if clusterInfo.Id == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(clusterInfo); err != nil {
			log.Print(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		log.Print(err)
	}

}

func ConnectCluster(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id int
	var err error
	var podCount int
	if id, err = strconv.Atoi(vars["id"]); err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	clusterInfo := RepoFindCluster(id)
	if clusterInfo.Id == 0 {
		// If we didn't find it, 404
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
			log.Print(err)
		}
		return
	}

	if podCount, err = kclient.ConnectToCluster(clusterInfo); err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Cluster has %d pods", podCount)

	w.WriteHeader(http.StatusOK)

	if _, err := fmt.Fprintf(w, "{Pods: %d}", podCount); err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func CreateCluster(w http.ResponseWriter, r *http.Request) {
	var clusterInfo model.ClusterInfo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Print(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Print(err)
	}
	if err := json.Unmarshal(body, &clusterInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(err)
		}
	}

	t := RepoCreateCluster(clusterInfo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Print(err)
	}
}

func DeleteCluster(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		RepoDeleteCluster(id)
		w.WriteHeader(http.StatusNoContent)
	}
	return
}

func UpdateCluster(w http.ResponseWriter, r *http.Request) {
	var clusterInfo model.ClusterInfo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Print(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Print(err)
	}
	if err := json.Unmarshal(body, &clusterInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(err)
		}
	}
	t := RepoUpdateCluster(clusterInfo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Print(err)
	}
}

//
// Service Handlers
//

func GetServices(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(serviceInfos); err != nil {
		log.Print(err)
	}
}

func GetService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id int
	var err error
	if id, err = strconv.Atoi(vars["id"]); err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	serviceInfo := RepoFindService(id)
	if serviceInfo.Id == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(serviceInfo); err != nil {
			log.Print(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		log.Print(err)
	}

}

func CreateService(w http.ResponseWriter, r *http.Request) {
	var serviceInfo model.ServiceInfo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Print(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Print(err)
	}
	if err := json.Unmarshal(body, &serviceInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(err)
		}
	}

	t := RepoCreateService(serviceInfo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Print(err)
	}
}

func DeleteService(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		RepoDeleteService(id)
		w.WriteHeader(http.StatusNoContent)
	}
	return
}

func UpdateService(w http.ResponseWriter, r *http.Request) {
	var serviceInfo model.ServiceInfo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Print(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Print(err)
	}
	if err := json.Unmarshal(body, &serviceInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(err)
		}
	}
	t := RepoUpdateService(serviceInfo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Print(err)
	}
}

//
// Factory Handlers
//

func GetFactories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(factoryInfos); err != nil {
		log.Print(err)
	}
}

func GetFactory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id int
	var err error
	if id, err = strconv.Atoi(vars["id"]); err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	factoryInfo := RepoFindFactory(id)
	if factoryInfo.Id == 0 {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(factoryInfo); err != nil {
			log.Print(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(jsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		log.Print(err)
	}

}

func CreateFactory(w http.ResponseWriter, r *http.Request) {
	var factoryInfo model.FactoryInfo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Print(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Print(err)
	}
	if err := json.Unmarshal(body, &factoryInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusUnprocessableEntity) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(err)
		}
	}

	t := RepoCreateFactory(factoryInfo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Print(err)
	}
}

func DeleteFactory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		st := stacktrace.New(err.Error())
		log.Printf("%s\n", st)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		RepoDeleteFactory(id)
		w.WriteHeader(http.StatusNoContent)
	}
	return
}

func UpdateFactory(w http.ResponseWriter, r *http.Request) {
	var factoryInfo model.FactoryInfo
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Print(err)
	}
	if err := r.Body.Close(); err != nil {
		log.Print(err)
	}
	if err := json.Unmarshal(body, &factoryInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Print(err)
		}
	}
	t := RepoUpdateFactory(factoryInfo)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		log.Print(err)
	}
}
