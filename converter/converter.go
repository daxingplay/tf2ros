package converter

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Converter(dir string, outputPath string) error {
	ret := "{\"ROSTemplateFormatVersion\": \"2015-09-01\",\"Transform\": \"Aliyun::Terraform-v1.1\",\"Parameters\": \"$PARAMETERS\",\"Metadata\": \"$METADATA\",\"Workspace\": \"$WORKSPACE\"}"

	// get metadata
	metadataContent, err := GetAdditionalContentFromFile(filepath.Join(dir, "./.metadata"))
	if err != nil {
		return err
	}
	ret = strings.Replace(ret, "\"$METADATA\"", metadataContent, 1)

	// get parameters
	parametersContent, err := GetAdditionalContentFromFile(filepath.Join(dir, "./.parameters"))
	if err != nil {
		return err
	}
	ret = strings.Replace(ret, "\"$PARAMETERS\"", parametersContent, 1)

	// get workspace
	workspace, err := GetTerraformFiles(dir)
	if err != nil {
		return err
	}
	workspaceContent, err := json.Marshal(workspace)
	if err != nil {
		return err
	}
	ret = strings.Replace(ret, "\"$WORKSPACE\"", string(workspaceContent), 1)

	prettyRet, err := PrettyString(ret)
	if err != nil {
		return err
	}

	// write file
	return ioutil.WriteFile(outputPath, []byte(prettyRet), 0644)
}

func fileExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func GetAdditionalContentFromFile(filename string) (string, error) {
	content := "{}"
	if fileExists(filename) {
		metadataFile, err := ioutil.ReadFile(filename)

		if err != nil {
			return "", err
		}

		content = string(metadataFile)
	}
	return content, nil
}

func GetTerraformFiles(dir string) (map[string]string, error) {
	workspace := make(map[string]string)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".tf") {
			rel, err := filepath.Rel(dir, path)

			if err != nil {
				return err
			}

			filecontent, err := ioutil.ReadFile(path)

			if err != nil {
				return err
			}

			workspace[rel] = string(filecontent)
		}
		return nil
	})

	return workspace, err
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
