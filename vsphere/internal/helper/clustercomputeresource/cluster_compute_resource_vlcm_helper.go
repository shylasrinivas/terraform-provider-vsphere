package clustercomputeresource

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"

	"gitlab.eng.vmware.com/vmehta/vsphere-client-bindings-go/esx/settings/clusters/software"
	"gitlab.eng.vmware.com/vmehta/vsphere-client-bindings-go/esx/settings/clusters/software/drafts"
)

// const basePath = ""
//
// func getSessionID() (string, error) {
// 	cfg := session.NewConfiguration()
// 	cfg.BasePath = basePath
// 	tr := &http.Transport{
// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
// 	}
// 	cfg.HTTPClient = &http.Client{Transport: tr}
// 	ctx := context.WithValue(context.Background(), session.ContextBasicAuth, session.BasicAuth{
// 		UserName: "",
// 		Password: "",
// 	})
//
// 	client := session.NewAPIClient(cfg)
// 	sessionId, _, err := client.CisSessionApi.Create(ctx)
//
// 	if err != nil {
// 		fmt.Println(err)
// 		return "", err
// 	} else {
// 		fmt.Println("CIS session ID is ", sessionId)
// 		return sessionId, nil
// 	}
// }

func GetSoftwareCfg(sessionId string, basePath string, insecureFlag bool) *software.Configuration {
	cfg := software.NewConfiguration()
	cfg.BasePath = basePath
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureFlag},
	}
	cfg.HTTPClient = &http.Client{Transport: tr}
	cfg.AddDefaultHeader("vmware-api-session-id", sessionId)
	return cfg
}

func PreCheck(sessionId string, cfg *software.Configuration, clusterComputeResource string, hosts []string, commit string) error {
	log.Printf("[DEBUG] Precheck for compute cluster %s started", clusterComputeResource)
	ctx := context.WithValue(context.Background(), software.ContextBasicAuth, software.APIKey{
		Key:    sessionId,
		Prefix: "Bearer",
	})

	client := software.NewAPIClient(cfg)
	settingsClustersSoftwareCheckSpec := software.SettingsClustersSoftwareCheckSpec{
		Commit: commit,
		Hosts:  hosts,
	}
	vmwTask := "true"
	log.Printf("[DEBUG] Begin call to check API")
	apiResponse, httpResponse, err := client.EsxSettingsClustersSoftwareApi.Checktask(ctx, vmwTask, clusterComputeResource, settingsClustersSoftwareCheckSpec)
	if err != nil {
		log.Printf("[DEBUG] Error during check %s ", err)
		return err
	} else {
		log.Printf("[DEBUG] apiResponse from check %s ", apiResponse)
		log.Println("[DEBUG] httpResponse HTTP apiResponse from check ", httpResponse)
	}
	log.Printf("[DEBUG] Precheck for compute cluster %s complete", clusterComputeResource)
	return nil
}

func Remediate(sessionId string, cfg *software.Configuration, clusterComputeResource string, hosts []string, commit string, acceptEula bool) error {
	log.Printf("[DEBUG] Remediation for compute cluster %s started", clusterComputeResource)
	ctx := context.WithValue(context.Background(), software.ContextBasicAuth, software.APIKey{
		Key:    sessionId,
		Prefix: "Bearer",
	})

	client := software.NewAPIClient(cfg)
	settingsClustersSoftwareApplySpec := software.SettingsClustersSoftwareApplySpec{
		Commit:     commit,
		Hosts:      hosts,
		AcceptEula: acceptEula,
	}
	vmwTask := "true"
	log.Printf("[DEBUG] Begin call to apply API")
	apiResponse, httpResponse, err := client.EsxSettingsClustersSoftwareApi.Applytask(ctx, vmwTask, clusterComputeResource, settingsClustersSoftwareApplySpec)
	if err != nil {
		log.Printf("[DEBUG] Error during apply %s ", err)
		return err
	} else {
		log.Printf("[DEBUG] apiResponse from apply %s ", apiResponse)
		log.Println("[DEBUG] httpResponse HTTP apiResponse from apply ", httpResponse)
	}
	log.Printf("[DEBUG] Remediation for compute cluster %s complete", clusterComputeResource)
	return nil
}

func Export(sessionId string, cfg *software.Configuration, clusterComputeResource string, exportSoftwareSpec bool, exportIsoImage bool, exportOfflineBundle bool) (string, error) {
	log.Printf("[DEBUG] Export compute cluster %s image", clusterComputeResource)
	ctx := context.WithValue(context.Background(), software.ContextBasicAuth, software.APIKey{
		Key:    sessionId,
		Prefix: "Bearer",
	})

	client := software.NewAPIClient(cfg)
	settingsClustersSoftwareExportSpec := software.SettingsClustersSoftwareExportSpec{
		ExportSoftwareSpec:  exportSoftwareSpec,
		ExportIsoImage:      exportIsoImage,
		ExportOfflineBundle: exportOfflineBundle,
	}

	log.Printf("[DEBUG] Begin call to export API")
	apiResponse, httpResponse, err := client.EsxSettingsClustersSoftwareApi.Export(ctx, clusterComputeResource, settingsClustersSoftwareExportSpec)
	if err != nil {
		log.Printf("[DEBUG] Error during export %s ", err)
		return "", err
	} else {
		log.Printf("[DEBUG] apiResponse from export %s ", apiResponse)
		log.Println("[DEBUG] httpResponse HTTP apiResponse from export ", httpResponse)
		// Eg apiResponse {"SOFTWARE_SPEC":"http://sc1-10-78-165-8.eng.vmware.com:9084/vum-filedownload/download?file=SOFTWARE_SPEC_1310523462.json"}
	}
	log.Printf("[DEBUG] Export compute cluster %s image complete", clusterComputeResource)
	exportedImage := apiResponse["SOFTWARE_SPEC"]
	return exportedImage, nil
}

func GetDraftsCfg(sessionId string, basePath string, insecureFlag bool) *drafts.Configuration {
	cfg := drafts.NewConfiguration()
	cfg.BasePath = basePath

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureFlag},
	}
	cfg.HTTPClient = &http.Client{Transport: tr}
	cfg.AddDefaultHeader("vmware-api-session-id", sessionId)
	return cfg
}

func Import(sessionId string, cfg *drafts.Configuration, cluster string, location string, fileId string, softwareSpec string) error {
	log.Printf("[DEBUG] Import compute cluster %s image", cluster)
	ctx := context.WithValue(context.Background(), drafts.ContextBasicAuth, drafts.APIKey{
		Key:    sessionId,
		Prefix: "Bearer",
	})

	client := drafts.NewAPIClient(cfg)
	// sourceType SettingsClustersSoftwareDraftsSourceType
	settingsClustersSoftwareDraftsImportSpec := drafts.SettingsClustersSoftwareDraftsImportSpec{
		SourceType: "PULL",
		// Location of the software specification file to be imported. This field is optional and it is only relevant when the value of Drafts.ImportSpec.source-type is PULL.
		Location: location,
		// File identifier returned by the file upload endpoint after file is uploaded. This field is optional and it is only relevant when the value of Drafts.ImportSpec.source-type is PUSH.
		FileId: fileId,
		// The JSON string representing the desired software specification. This field is optional and it is only relevant when the value of Drafts.ImportSpec.source-type is JSON_STRING.
		SoftwareSpec: softwareSpec,
	}

	log.Printf("[DEBUG] Begin call to import API")
	apiResponse, httpResponse, err := client.EsxSettingsClustersSoftwareDraftsApi.ImportSoftwareSpec(ctx, cluster, settingsClustersSoftwareDraftsImportSpec)
	if err != nil {
		log.Printf("[DEBUG] Error during import %s ", err)
		return err
	} else {
		log.Printf("[DEBUG] apiResponse from import %s ", apiResponse)
		log.Println("[DEBUG] httpResponse HTTP apiResponse from import ", httpResponse)
	}
	log.Printf("[DEBUG] Import compute cluster %s image complete", cluster)
	return nil
}

func Commit(sessionId string, cfg *drafts.Configuration, cluster string, draft string, commitMessage string) error {
	log.Printf("[DEBUG] Commit compute cluster %s image started", cluster)
	ctx := context.WithValue(context.Background(), drafts.ContextBasicAuth, drafts.APIKey{
		Key:    sessionId,
		Prefix: "Bearer",
	})

	client := drafts.NewAPIClient(cfg)
	settingsClustersSoftwareDraftsCommitSpec := drafts.SettingsClustersSoftwareDraftsCommitSpec{
		// Message to include with the commit. If unset, message is set to empty string.
		Message: commitMessage,
	}
	vmwTask := "true"
	log.Printf("[DEBUG] Begin call to commit API")
	apiResponse, httpResponse, err := client.EsxSettingsClustersSoftwareDraftsApi.Committask(ctx, vmwTask, cluster, draft, settingsClustersSoftwareDraftsCommitSpec)
	if err != nil {
		log.Printf("[DEBUG] Error during commit %s ", err)
		return err
	} else {
		log.Printf("[DEBUG] apiResponse from commit %s ", apiResponse)
		log.Println("[DEBUG] httpResponse HTTP apiResponse from commit ", httpResponse)
	}
	log.Printf("[DEBUG] Commit compute cluster %s image complete", cluster)
	return nil
}
