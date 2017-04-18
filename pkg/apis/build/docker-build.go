// Copyright © 2017 huang jia <449264675@qq.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package build

import (
	"encoding/json"
	"net/http"

	"apiserver/pkg/api/build"
	r "apiserver/pkg/router"
	"apiserver/pkg/util/log"

	"github.com/gorilla/mux"
	"context"
	"github.com/docker/docker/api/types"
	"apiserver/pkg/storage/mysqld"
	"os"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/container"
	//"github.com/google/go-github/github"
    	"golang.org/x/oauth2"
    	githuboauth "golang.org/x/oauth2/github"
	"fmt"
)

var (
	engine = mysqld.GetEngine()
	oauthConf = &oauth2.Config{
		ClientID:     "17f03800ee9e0b7caf95",
		ClientSecret: "542dd6ffef2da3ac66ed493e82fba8fd0de089ba",
		//ClientID:     "",
		//ClientSecret: "",
		// select level of access you want https://developer.github.com/v3/oauth/#scopes
		Scopes:       []string{"user:email", "repo"},
		Endpoint:     githuboauth.Endpoint,
    }
    // random string for oauth2 API calls to protect against CSRF
    oauthStateString = "thisshouldberandom"
)

const(
	OK int = iota
	ERROR
)

func Register(rout *mux.Router) {
	r.RegisterHttpHandler(rout, "/build", "POST", OnlineBuild)
	r.RegisterHttpHandler(rout, "/build", "PUT", OfflineBuild)
}

//OnlineBuild build application online
func OnlineBuild(request *http.Request) (string, interface{}) {
	decoder := json.NewDecoder(request.Body)
	builder := &build.Build{}
	err := decoder.Decode(builder)
	if err != nil {
		log.Errorf("decode the request body err:%v", err)
		return r.StatusBadRequest, "json format error"
	}


	result,err := BuildImage(builder)
	if err != nil && result != OK {
		log.Errorf("build image failed,error:%v",err)
		return r.StatusBadRequest,"build image error"
	}

	return r.StatusCreated, nil
}

//OfflineBuild build application offline
func OfflineBuild(request *http.Request) (string, interface{}) {
	decoder := json.NewDecoder(request.Body)
	builder := &build.Build{}
	err := decoder.Decode(builder)
	if err != nil {
		log.Errorf("decode the request body err:%v", err)
		return r.StatusBadRequest, "json format error"
	}

	result,err := BuildImage(builder)
	if err != nil && result != OK {
		log.Errorf("build image failed,error:%v",err)
		return r.StatusBadRequest,"build image error"
	}

	return r.StatusCreated, nil
}

func BuildImage(builder *build.Build) (int,error) {
	hosts := fmt.Sprintf("%s:%s",builder.Host,builder.Port)
	cli := &http.Client{
		Transport: new(http.Transport),
	}
	DockerClient, err := client.NewClient(hosts,builder.Apiversion , cli, nil)
	if err != nil{
		fmt.Println("error in ImageBuild:",err)
		panic(err)
	}
	ctx := context.Background()
	tags := fmt.Sprintf("%s:%s",builder.Image,builder.Version)

	options := types.ImageBuildOptions{
		Tags: []string{tags},
		SuppressOutput: true,
		Isolation: container.Isolation("default"),
		NoCache: true,
		Remove: false,
		Labels: map[string]string{builder.Image:builder.Version},
		ForceRemove: true,
		PullParent: false,
		Dockerfile: "Dockerfile",
	}
	buildcontext,err := os.Open(builder.Path)
	if err != nil{
		log.Error("error in Open buildcontext:",err)
		return ERROR,err
	}
	defer buildcontext.Close()
	resp,err := DockerClient.ImageBuild(ctx,buildcontext,options)

	if err != nil && resp.OSType != ""{
		log.Error("error in ImageBuild:",err)
		return ERROR,err
	}
	return OK,nil
}

