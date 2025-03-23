package document

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"uni-resolver-driver-did-cndid/utils"
)

const resolveIP = "139.198.21.202"
const resolvePort = "31005"
const resolutionServiceURL = "http://" + resolveIP + ":" + resolvePort + "/resolve"

// Response 统一的响应结构
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func isValidHexAddress(address string) bool {
	return utils.StringToAddress(address).EqualString(address)
}

// 写入响应
func writeResponse(w http.ResponseWriter, status int, message string, data interface{}) {
	response := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func ResolveDID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/1.0/identifiers/"):]
	if id == "" {
		logrus.Error("ID is required")
		writeResponse(w, http.StatusBadRequest, "ID is required", nil)
		return
	}

	// 构造解析服务的请求 URL
	query := url.Values{}
	query.Set("id", id)
	resolutionURL := fmt.Sprintf("%s?%s", resolutionServiceURL, query.Encode())

	// 初步判断地址是否合法
	if !isValidHexAddress(id) {
		logrus.Error("Invalid DID Address")
		writeResponse(w, http.StatusUnavailableForLegalReasons, "Invalid DID Address", nil)
		return
	}

	// 发送请求到解析服务
	resp, err := http.Get(resolutionURL)
	if err != nil {
		logrus.Errorf("Failed to resolve DID: %s", err)
		writeResponse(w, http.StatusInternalServerError, "Failed to resolve DID", nil)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			logrus.Error("Failed to resolve DID, CN-DID document not exist")
			writeResponse(w, resp.StatusCode, "Failed to resolve DID, not exist", nil)
			return
		}
		logrus.Errorf("Failed to resolve DID, response err : %s", resp.Status)
		writeResponse(w, resp.StatusCode, "Failed to resolve DID", nil)
		return
	}

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Failed to read response body: %s", err)
		writeResponse(w, http.StatusInternalServerError, "Failed to read response body", nil)
		return
	}

	// 解析响应内容
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		logrus.Errorf("Failed to parse response body: %s", err)
		writeResponse(w, http.StatusInternalServerError, "Failed to parse response body", nil)
		return
	}

	// 检查状态码和消息
	if response.Status != http.StatusOK || response.Message != "Success" {
		logrus.Errorf("Failed to resolve DID: status=%d, message=%s", response.Status, response.Message)
		writeResponse(w, response.Status, response.Message, nil)
		return
	}

	// 提取 data 字段
	var result DIDResolution
	if response.Data != nil {
		dataBytes, err := json.Marshal(response.Data)
		if err != nil {
			logrus.Errorf("Failed to marshal data field: %s", err)
			writeResponse(w, http.StatusInternalServerError, "Failed to marshal data field", nil)
			return
		}

		if err := json.Unmarshal(dataBytes, &result); err != nil {
			logrus.Errorf("Failed to parse data field: %s", err)
			writeResponse(w, http.StatusInternalServerError, "Failed to parse data field", nil)
			return
		}
	}

	// 在序列化之前，手动移除不存在的字段
	cleanedResult := cleanOptionalDIDResolution(result)

	// 返回响应内容
	writeResponse(w, http.StatusOK, "Success", cleanedResult)
}

// cleanOptionalDIDResolution 用于移除非必填的字段
func cleanOptionalDIDResolution(res DIDResolution) map[string]interface{} {
	result := make(map[string]interface{})

	// 必填字段
	result["context"] = res.Context
	result["version"] = res.Version
	result["id"] = res.ID
	result["publicKey"] = res.PublicKey
	result["authentication"] = res.Authentication
	result["extension"] = cleanExtension(res.Extension)
	result["created"] = res.Created
	result["updated"] = res.Updated

	// 可选字段
	if len(res.AlsoKnownAs) > 0 {
		result["alsoKnownAs"] = res.AlsoKnownAs
	}
	if len(res.Service) > 0 {
		var services []map[string]interface{}
		for _, serv := range res.Service {
			service := cleanService(serv)
			services = append(services, service)
		}
		result["service"] = services
	}

	if res.Proof.Creator != "" && res.Proof.SignatureValue != "" {
		result["proof"] = res.Proof
	}

	return result
}

// cleanExtension 用于清理 Extension 字段
func cleanExtension(ext Extension) map[string]interface{} {
	result := make(map[string]interface{})

	// 必填字段
	result["ttl"] = ext.TTL
	result["type"] = ext.Type

	var cleanedAttributes []map[string]interface{}
	for _, attr := range ext.Attributes {
		cleanedAttr := cleanAttribute(attr)
		if len(cleanedAttr) > 0 {
			cleanedAttributes = append(cleanedAttributes, cleanedAttr)
		}
	}
	result["attributes"] = cleanedAttributes

	// 可选字段
	if len(ext.Recovery) > 0 {
		result["recovery"] = ext.Recovery
	}
	if ext.DelegateSign.Signer != "" || ext.DelegateSign.SignatureValue != "" {
		result["delegateSign"] = ext.DelegateSign
	}
	if len(ext.VerifiableCredentials) > 0 {
		result["verifiableCredentials"] = ext.VerifiableCredentials
	}

	return result
}

// cleanAttribute 用于清理单个 Attribute 字段
func cleanAttribute(attr Attribute) map[string]interface{} {
	result := make(map[string]interface{})

	// 必填字段
	result["key"] = attr.Key

	// 可选字段
	if attr.Desc != "" {
		result["desc"] = attr.Desc
	}
	if attr.Encrypt != 0 {
		result["encrypt"] = attr.Encrypt
	}
	if attr.Format != "" {
		result["format"] = attr.Format
	}
	if attr.Value != "" {
		result["value"] = attr.Value
	}

	return result
}

// cleanService 用于清理单个 Service 字段
func cleanService(service Service) map[string]interface{} {
	result := make(map[string]interface{})

	// 必填字段
	result["id"] = service.ID
	result["type"] = service.Type
	result["serviceEndpoint"] = service.ServiceEndpoint

	// 根据 type 字段决定是否包含其他字段
	if service.Type == "DIDSubResolver" {
		result["version"] = service.Version
		result["serverType"] = service.ServerType
		result["protocol"] = service.Protocol
		result["port"] = service.Port
	}

	return result
}
