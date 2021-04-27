package common

import (
	"context"
	"fmt"
	"github.com/citihub/probr-sdk/audit"
	"github.com/citihub/probr-sdk/probeengine"
	"github.com/citihub/probr-sdk/probeengine/opa"
	"github.com/citihub/probr-sdk/providers/azure/connection"
	"github.com/citihub/probr-sdk/utils"
	"log"
	"strings"
)

// ScenarioState is the base struct for handling state across steps in a scenario
type ScenarioState struct {
	Name         string
	CurrentStep  string
	Audit        *audit.ScenarioAudit
	Probe        *audit.Probe
	Ctx          context.Context
	Tags         map[string]*string
	AZConnection connection.Azure
	//additional variables to hold state goes here
}

const opaPackageName = "aks"

// GetScenarioState returns the ScenarioState struct - useful if used as an embedded field (see comment on OPAProbe)
func (s ScenarioState) GetScenarioState() ScenarioState {
	return s
}

/*
   OPAProbe is a common function for any features using OPA as part of the scenario.
   To use this function the scenarioState that is part of the standard Probr should use
   common.ScenarioState as an embedded field to base the struct from, as follows:

   type scenarioState struct {
     common.ScenarioState
     more                  string
     fields                int
   }

   Then pass a pointer to common.ScenarioState into OPAProbe as follows:

   baseState := scenario.GetScenarioState()
   common.OPAProbe("opa_function", json, &baseState)
*/

func OPAProbe(opaFuncName string, aksJson []byte, scenario *ScenarioState) (err error) {

	var stepTrace strings.Builder

	stepTrace, payload, err := utils.AuditPlaceholders()
	if err != nil {
		return err
	}

	defer func() {
		scenario.Audit.AuditScenarioStep(scenario.CurrentStep, stepTrace.String(), payload, err)
	}()

	payload = struct {
		Placeholder string
	}{
		"placeholder",
	}

	stepTrace.WriteString(fmt.Sprintf("Use OPA function %s to evaluate this cluster; ", opaFuncName))

	err = eval(opaFuncName, aksJson)
	return
}

func eval(functionName string, aksJson []byte) (err error) {
	// true = return nil
	// false = return new err
	// err = return err

	regoFilePath := probeengine.GetFilePath("internal", "common", "aks.rego")
	log.Printf("[DEBUG] common.go: eval(): regoFilePath = %s", regoFilePath)
	var r bool

	r, err = opa.Eval(regoFilePath, opaPackageName, functionName, &aksJson)

	if err != nil {
		log.Printf("[DEBUG] opa.Eval returned an error")
		return
	}

	if r == false {
		log.Printf("[DEBUG] Rego function %s returned result of 'non-compliant'", functionName)
		err = fmt.Errorf("Rego function %s returned result of 'non-compliant'", functionName)
	} else {
		log.Printf("[DEBUG] Rego function %s returned result of 'compliant'", functionName)
		err = nil
	}
	return
}

func AnAzureKubernetesClusterWeCanReadTheConfigurationOf(scenario *ScenarioState) (json []byte, err error) {
	var stepTrace strings.Builder

	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		scenario.Audit.AuditScenarioStep(scenario.CurrentStep, stepTrace.String(), payload, err)
	}()

	payload = struct {
		Placeholder string
	}{
		"placeholder",
	}

	stepTrace.WriteString("Get the configuration of the AKS cluster; ")

	json, err = getClusterConfigJSON(scenario)

	if err != nil {
		log.Printf("Error loading JSON: %v", err)
		return
	}

	if len(json) == 0 {
		err = fmt.Errorf("aksJson empty")
	}
	return
}

func getClusterConfigJSON(scenario *ScenarioState) (json []byte, err error) {
	json, err = scenario.AZConnection.GetManagedClusterJSON("probr-demo-rg", "probr-demo-cluster")
	if err != nil {
		log.Printf("Error loading JSON: %v", err)
		return
	}

	if len(json) == 0 {
		err = fmt.Errorf("aksJson empty")
	}
	return
}
