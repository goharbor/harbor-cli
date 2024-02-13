package utils
import (
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/ghodss/yaml"
	"github.com/goharbor/go-client/pkg/harbor"
	v2client "github.com/goharbor/go-client/pkg/sdk/v2.0/client"
)
// Returns Harbor v2 client for given clientConfig
func GetClientByConfig(clientConfig *harbor.ClientSetConfig) *v2client.HarborAPI {
	cs, err := harbor.NewClientSet(clientConfig)
	if err != nil {
		panic(err)
	}
	return cs.V2()
}
// Returns Harbor v2 client after resolving the credential name
func GetClientByCredentialName(credentialName string) *v2client.HarborAPI {
	credential, err := resolveCredential(credentialName)
	if err != nil {
		panic(err)
	}
	clientConfig := &harbor.ClientSetConfig{
		URL:      credential.ServerAddress,
		Username: credential.Username,
		Password: credential.Password,
	}
	return GetClientByConfig(clientConfig)
}

func PrintPayloadFormat(otype string, payload interface{}) string {
	var res []byte
	output, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	if otype == "yaml" {
		res, _ = yaml.JSONToYAML(output)
		return string(res)
	}
	if otype == "json" {

		return string(output)
	}
	return string(output)
}

func ConvertSize(sizeinByte int64) string {
	size := sizeinByte/1024
	if (sizeinByte/1024) > 1024 {
		return strconv.Itoa(int(size)) + "GiB"
	} else {
		return strconv.Itoa(int(size)) + "KiB"
	}
}

func ShortText(s string, i int) string {
    if len(s) < i {
        return s
    }
    if utf8.ValidString(s[:i]) {
        return s[:i]
    }
    return s[:i+1]
}

func AccessCheck(check string) string {
	if check == "true" {
		return "public"
	}
	return "private"
}