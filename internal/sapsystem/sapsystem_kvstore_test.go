package sapsystem

import (
	"fmt"
	"os"
	"path"
	"testing"

	consulApi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/trento-project/trento/internal/consul"
	consulMocks "github.com/trento-project/trento/internal/consul/mocks"
	"github.com/trento-project/trento/internal/sapsystem/sapcontrol"
	sapcontrolMocks "github.com/trento-project/trento/internal/sapsystem/sapcontrol/mocks"
)

func TestStore(t *testing.T) {
	host, _ := os.Hostname()
	consulInst := new(consulMocks.Client)
	kv := new(consulMocks.KV)

	consulInst.On("KV").Return(kv)
	kvPath := path.Join(fmt.Sprintf(consul.KvHostsSAPSystemPath, host), "PRD")

	mockWebService := new(sapcontrolMocks.WebService)

	expectedPutMap := map[string]interface{}{
		"sid":  "PRD",
		"type": Database,
		"instances": map[string]*SAPInstance{
			"HDB00": &SAPInstance{
				Name: "HDB00",
				Type: Database,
				Host: "host1",
				SAPControl: &SAPControl{
					webService: mockWebService,
					Processes: map[string]*sapcontrol.OSProcess{
						"enserver": {
							Name:        "enserver",
							Description: "foobar",
							Dispstatus:  sapcontrol.STATECOLOR_GREEN,
							Textstatus:  "Running",
							Starttime:   "1",
							Elapsedtime: "2",
							Pid:         30787,
						},
						"msg_server": {
							Name:        "msg_server",
							Description: "foobar2",
							Dispstatus:  sapcontrol.STATECOLOR_YELLOW,
							Textstatus:  "Stopping",
							Starttime:   "43",
							Elapsedtime: "",
							Pid:         30786,
						},
					},
					Properties: map[string]*sapcontrol.InstanceProperty{
						"INSTANCE_NAME": {
							Property:     "INSTANCE_NAME",
							Propertytype: "string",
							Value:        "HDB00",
						},
						"SAPSYSTEMNAME": {
							Property:     "SAPSYSTEMNAME",
							Propertytype: "string",
							Value:        "PRD",
						},
						"HANA Roles": {
							Property:     "HANA Roles",
							Propertytype: "type3",
							Value:        "some hana value",
						},
					},
					Instances: map[string]*sapcontrol.SAPInstance{
						"host1": {
							Hostname:      "host1",
							InstanceNr:    0,
							HttpPort:      50013,
							HttpsPort:     50014,
							StartPriority: "0.3",
							Features:      "some features",
							Dispstatus:    sapcontrol.STATECOLOR_GREEN,
						},
						"host2": {
							Hostname:      "host2",
							InstanceNr:    1,
							HttpPort:      50113,
							HttpsPort:     50114,
							StartPriority: "0.3",
							Features:      "some other features",
							Dispstatus:    sapcontrol.STATECOLOR_YELLOW,
						},
					},
				},
			},
		},
	}

	kv.On("DeleteTree", kvPath, (*consulApi.WriteOptions)(nil)).Return(nil, nil)
	kv.On("PutMap", kvPath, expectedPutMap).Return(nil, nil)
	testLock := consulApi.Lock{}
	consulInst.On("AcquireLockKey", fmt.Sprintf(consul.KvHostsSAPSystemPath, host)).Return(&testLock, nil)

	s := SAPSystem{
		SID:  "PRD",
		Type: Database,
		Instances: map[string]*SAPInstance{
			"HDB00": &SAPInstance{
				Name: "HDB00",
				Type: Database,
				Host: "host1",
				SAPControl: &SAPControl{
					webService: mockWebService,
					Processes: map[string]*sapcontrol.OSProcess{
						"enserver": {
							Name:        "enserver",
							Description: "foobar",
							Dispstatus:  sapcontrol.STATECOLOR_GREEN,
							Textstatus:  "Running",
							Starttime:   "1",
							Elapsedtime: "2",
							Pid:         30787,
						},
						"msg_server": {
							Name:        "msg_server",
							Description: "foobar2",
							Dispstatus:  sapcontrol.STATECOLOR_YELLOW,
							Textstatus:  "Stopping",
							Starttime:   "43",
							Elapsedtime: "",
							Pid:         30786,
						},
					},
					Properties: map[string]*sapcontrol.InstanceProperty{
						"INSTANCE_NAME": {
							Property:     "INSTANCE_NAME",
							Propertytype: "string",
							Value:        "HDB00",
						},
						"SAPSYSTEMNAME": {
							Property:     "SAPSYSTEMNAME",
							Propertytype: "string",
							Value:        "PRD",
						},
						"HANA Roles": {
							Property:     "HANA Roles",
							Propertytype: "type3",
							Value:        "some hana value",
						},
					},
					Instances: map[string]*sapcontrol.SAPInstance{
						"host1": {
							Hostname:      "host1",
							InstanceNr:    0,
							HttpPort:      50013,
							HttpsPort:     50014,
							StartPriority: "0.3",
							Features:      "some features",
							Dispstatus:    sapcontrol.STATECOLOR_GREEN,
						},
						"host2": {
							Hostname:      "host2",
							InstanceNr:    1,
							HttpPort:      50113,
							HttpsPort:     50114,
							StartPriority: "0.3",
							Features:      "some other features",
							Dispstatus:    sapcontrol.STATECOLOR_YELLOW,
						},
					},
				},
			},
		},
	}

	err := s.Store(consulInst)

	assert.NoError(t, err)

	kv.AssertExpectations(t)
}

func TestLoad(t *testing.T) {
	host, _ := os.Hostname()
	kvPath := fmt.Sprintf(consul.KvHostsSAPSystemPath, host)
	consulInst := new(consulMocks.Client)
	kv := new(consulMocks.KV)

	listMap := map[string]interface{}{
		"PRD": map[string]interface{}{
			"sid":  "PRD",
			"type": Database,
			"instances": map[string]interface{}{
				"HDB00": map[string]interface{}{
					"name": "HDB00",
					"type": Database,
					"host": "host1",
					"sapcontrol": map[string]interface{}{
						"processes": map[string]interface{}{
							"enserver": map[string]interface{}{
								"Name":        "enserver",
								"Description": "foobar",
								"Dispstatus":  sapcontrol.STATECOLOR_GREEN,
								"Textstatus":  "Running",
								"Starttime":   "1",
								"Elapsedtime": "2",
								"Pid":         30787,
							},
							"msg_server": map[string]interface{}{
								"Name":        "msg_server",
								"Description": "foobar2",
								"Dispstatus":  sapcontrol.STATECOLOR_YELLOW,
								"Textstatus":  "Stopping",
								"Starttime":   "43",
								"Elapsedtime": "",
								"Pid":         30786,
							},
						},
						"properties": map[string]interface{}{
							"INSTANCE_NAME": map[string]interface{}{
								"Property":     "INSTANCE_NAME",
								"Propertytype": "string",
								"Value":        "HDB00",
							},
							"SAPSYSTEMNAME": map[string]interface{}{
								"Property":     "SAPSYSTEMNAME",
								"Propertytype": "string",
								"Value":        "PRD",
							},
							"HANA Roles": map[string]interface{}{
								"Property":     "HANA Roles",
								"Propertytype": "type3",
								"Value":        "some hana value",
							},
						},
						"instances": map[string]interface{}{
							"host1": map[string]interface{}{
								"Hostname":      "host1",
								"InstanceNr":    0,
								"HttpPort":      50013,
								"HttpsPort":     50014,
								"StartPriority": "0.3",
								"Features":      "some features",
								"Dispstatus":    sapcontrol.STATECOLOR_GREEN,
							},
							"host2": map[string]interface{}{
								"Hostname":      "host2",
								"InstanceNr":    1,
								"HttpPort":      50113,
								"HttpsPort":     50114,
								"StartPriority": "0.3",
								"Features":      "some other features",
								"Dispstatus":    sapcontrol.STATECOLOR_YELLOW,
							},
						},
					},
				},
			},
		},
	}

	kv.On("ListMap", kvPath, kvPath).Return(listMap, nil)
	consulInst.On("WaitLock", fmt.Sprintf(consul.KvHostsSAPSystemPath, host)).Return(nil)

	consulInst.On("KV").Return(kv)

	systems, _ := Load(consulInst, host)

	expectedSystems := map[string]*SAPSystem{
		"PRD": {
			SID:  "PRD",
			Type: Database,
			Instances: map[string]*SAPInstance{
				"HDB00": &SAPInstance{
					Name: "HDB00",
					Type: Database,
					Host: "host1",
					SAPControl: &SAPControl{
						Processes: map[string]*sapcontrol.OSProcess{
							"enserver": {
								Name:        "enserver",
								Description: "foobar",
								Dispstatus:  sapcontrol.STATECOLOR_GREEN,
								Textstatus:  "Running",
								Starttime:   "1",
								Elapsedtime: "2",
								Pid:         30787,
							},
							"msg_server": {
								Name:        "msg_server",
								Description: "foobar2",
								Dispstatus:  sapcontrol.STATECOLOR_YELLOW,
								Textstatus:  "Stopping",
								Starttime:   "43",
								Elapsedtime: "",
								Pid:         30786,
							},
						},
						Properties: map[string]*sapcontrol.InstanceProperty{
							"INSTANCE_NAME": {
								Property:     "INSTANCE_NAME",
								Propertytype: "string",
								Value:        "HDB00",
							},
							"SAPSYSTEMNAME": {
								Property:     "SAPSYSTEMNAME",
								Propertytype: "string",
								Value:        "PRD",
							},
							"HANA Roles": {
								Property:     "HANA Roles",
								Propertytype: "type3",
								Value:        "some hana value",
							},
						},
						Instances: map[string]*sapcontrol.SAPInstance{
							"host1": {
								Hostname:      "host1",
								InstanceNr:    0,
								HttpPort:      50013,
								HttpsPort:     50014,
								StartPriority: "0.3",
								Features:      "some features",
								Dispstatus:    sapcontrol.STATECOLOR_GREEN,
							},
							"host2": {
								Hostname:      "host2",
								InstanceNr:    1,
								HttpPort:      50113,
								HttpsPort:     50114,
								StartPriority: "0.3",
								Features:      "some other features",
								Dispstatus:    sapcontrol.STATECOLOR_YELLOW,
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expectedSystems, systems)
}