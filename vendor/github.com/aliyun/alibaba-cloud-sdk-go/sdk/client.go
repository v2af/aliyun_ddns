/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sdk

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"net"
	"net/http"
	"fmt"
	"strconv"
)

// this value will be replaced while build: -ldflags="-X sdk.version=x.x.x"
var Version = "0.0.1"

type Client struct {
	regionId       string
	config         *Config
	signer         auth.Signer
	httpClient     *http.Client
	asyncTaskQueue chan func()

	debug     bool
	isRunning bool
}

func (client *Client) Init() (err error) {
	panic("not support yet")
}

func (client *Client) InitWithOptions(regionId string, config *Config, credential auth.Credential) (err error) {
	client.isRunning = true
	client.regionId = regionId
	client.config = config
	if err != nil {
		return
	}
	client.httpClient = &http.Client{}

	if config.HttpTransport != nil {
		client.httpClient.Transport = config.HttpTransport
	}

	if config.Timeout > 0 {
		client.httpClient.Timeout = config.Timeout
	}

	if config.EnableAsync {
		client.EnableAsync(config.GoRoutinePoolSize, config.MaxTaskQueueSize)
	}

	client.signer, err = auth.NewSignerWithCredential(credential, client.ProcessCommonRequestWithSigner)

	return
}

func (client *Client) EnableAsync(routinePoolSize, maxTaskQueueSize int) {
	client.asyncTaskQueue = make(chan func(), maxTaskQueueSize)
	for i := 0; i < routinePoolSize; i++ {
		go func() {
			for client.isRunning {
				select {
				case task := <-client.asyncTaskQueue:
					task()
				}
			}
		}()
	}
}

func (client *Client) InitWithAccessKey(regionId, accessKeyId, accessKeySecret string) (err error) {
	config := client.InitClientConfig()
	credential := &credentials.BaseCredential{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	return client.InitWithOptions(regionId, config, credential)
}

func (client *Client) InitWithSecurityToken(regionId, accessKeyId, accessKeySecret, securityToken string) (err error) {
	config := client.InitClientConfig()
	credential := &credentials.StsTokenCredential{
		AccessKeyId:       accessKeyId,
		AccessKeySecret:   accessKeySecret,
		AccessKeyStsToken: securityToken,
	}
	return client.InitWithOptions(regionId, config, credential)
}

func (client *Client) InitWithStsRoleArn(regionId, accessKeyId, accessKeySecret, roleArn, roleSessionName string) (err error) {
	config := client.InitClientConfig()
	credential := &credentials.StsRoleArnCredential{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		RoleArn:         roleArn,
		RoleSessionName: roleSessionName,
	}
	return client.InitWithOptions(regionId, config, credential)
}

func (client *Client) InitWithRsaKeyPair(regionId, publicKeyId, privateKey string, sessionExpiration int) (err error) {
	config := client.InitClientConfig()
	credential := &credentials.RsaKeyPairCredential{
		PrivateKey:        privateKey,
		PublicKeyId:       publicKeyId,
		SessionExpiration: sessionExpiration,
	}
	return client.InitWithOptions(regionId, config, credential)
}

func (client *Client) InitWithStsRoleNameOnEcs(regionId, roleName string) (err error) {
	config := client.InitClientConfig()
	credential := &credentials.StsRoleNameOnEcsCredential{
		RoleName: roleName,
	}
	return client.InitWithOptions(regionId, config, credential)
}

func (client *Client) InitClientConfig() (config *Config) {
	if client.config != nil {
		return client.config
	}else{
		return NewConfig()
	}
}

func (client *Client) DoAction(request requests.AcsRequest, response responses.AcsResponse) (err error) {
	return client.DoActionWithSigner(request, response, nil)
}

func (client *Client) DoActionWithSigner(request requests.AcsRequest, response responses.AcsResponse, signer auth.Signer) (err error) {

	// add clientVersion
	request.GetHeaders()["x-sdk-core-version"] = Version

	regionId := client.regionId
	if len(request.GetRegionId()) > 0 {
		regionId = request.GetRegionId()
	}

	// resolve endpoint
	resolveParam := &endpoints.ResolveParam{
		Domain:               request.GetDomain(),
		Product:              request.GetProduct(),
		RegionId:             regionId,
		LocationProduct:      request.GetLocationServiceCode(),
		LocationEndpointType: request.GetLocationEndpointType(),
		CommonApi:            client.ProcessCommonRequest,
	}
	endpoint, err := endpoints.Resolve(resolveParam)
	if err != nil {
		return
	}
	request.SetDomain(endpoint)

	// init request params
	err = requests.InitParams(request)
	if err != nil {
		return
	}

	// signature
	if signer != nil {
		err = auth.Sign(request, signer, regionId)
	} else {
		err = auth.Sign(request, client.signer, regionId)
	}

	if err != nil {
		return
	}

	requestMethod := request.GetMethod()
	requestUrl := request.GetUrl()
	body := request.GetBodyReader()
	httpRequest, err := http.NewRequest(requestMethod, requestUrl, body)
	if err != nil {
		return
	}
	for key, value := range request.GetHeaders() {
		httpRequest.Header[key] = []string{value}
	}
	var httpResponse *http.Response
	for retryTimes := 0; retryTimes <= client.config.MaxRetryTime; retryTimes++ {
		httpResponse, err = client.httpClient.Do(httpRequest)

		var timeout bool
		// receive error
		if err != nil {
			if timeout = isTimeout(err); !timeout {
				// if not timeout error, return
				return
			} else if retryTimes >= client.config.MaxRetryTime {
				// timeout but reached the max retry times, return
				timeoutErrorMsg := fmt.Sprintf(errors.TimeoutErrorMessage, strconv.Itoa(retryTimes + 1), strconv.Itoa(retryTimes + 1))
				err = errors.NewClientError(errors.TimeoutErrorCode, timeoutErrorMsg, err)
				return
			}
		}
		//  if status code >= 500 or timeout, will trigger retry
		if client.config.AutoRetry && (timeout || isServerError(httpResponse)) {
			continue
		}
		break
	}
	err = responses.Unmarshal(response, httpResponse, request.GetAcceptFormat())
	return
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	netErr, isNetError := err.(net.Error)
	return isNetError && netErr.Timeout()
}

func isServerError(httpResponse *http.Response) bool {
	return httpResponse.StatusCode >= http.StatusInternalServerError
}

func (client *Client) AddAsyncTask(task func()) (err error) {
	if client.asyncTaskQueue != nil {
		client.asyncTaskQueue <- task
	} else {
		err = errors.NewClientError(errors.AsyncFunctionNotEnabledCode, errors.AsyncFunctionNotEnabledMessage, nil)
	}
	return
}

func NewClient() (client *Client, err error) {
	client = &Client{}
	err = client.Init()
	return
}

func NewClientWithOptions(regionId string, config *Config, credential auth.Credential) (client *Client, err error) {
	client = &Client{}
	err = client.InitWithOptions(regionId, config, credential)
	return
}

func NewClientWithAccessKey(regionId, accessKeyId, accessKeySecret string) (client *Client, err error) {
	client = &Client{}
	err = client.InitWithAccessKey(regionId, accessKeyId, accessKeySecret)
	return
}

func NewClientWithStsToken(regionId, accessKeyId, accessKeySecret, accessKeyStsToken string) (client *Client, err error) {
	client = &Client{}
	err = client.InitWithSecurityToken(regionId, accessKeyId, accessKeySecret, accessKeyStsToken)
	return
}

func NewClientWithRsaKeyPair(regionId string, publicKeyId, privateKey string, sessionExpiration int) (client *Client, err error) {
	client = &Client{}
	err = client.InitWithRsaKeyPair(regionId, publicKeyId, privateKey, sessionExpiration)
	return
}

func NewClientWithStsRoleNameOnEcs(regionId string, roleName string) (client *Client, err error) {
	client = &Client{}
	err = client.InitWithStsRoleNameOnEcs(regionId, roleName)
	return
}

func NewClientWithStsRoleArn(regionId string, accessKeyId, accessKeySecret, roleArn, roleSessionName string) (client *Client, err error) {
	client = &Client{}
	err = client.InitWithStsRoleArn(regionId, accessKeyId, accessKeySecret, roleArn, roleSessionName)
	return
}

func (client *Client) ProcessCommonRequest(request *requests.CommonRequest) (response *responses.CommonResponse, err error) {
	request.TransToAcsRequest()
	response = responses.NewCommonResponse()
	err = client.DoAction(request, response)
	return
}

func (client *Client) ProcessCommonRequestWithSigner(request *requests.CommonRequest, signerInterface interface{}) (response *responses.CommonResponse, err error) {
	if signer, isSigner := signerInterface.(auth.Signer); isSigner {
		request.TransToAcsRequest()
		response = responses.NewCommonResponse()
		err = client.DoActionWithSigner(request, response, signer)
		return
	} else {
		panic("should not be here")
	}
}

func (client *Client) Shutdown() {
	client.signer.Shutdown()
	close(client.asyncTaskQueue)
	client.isRunning = false
}
