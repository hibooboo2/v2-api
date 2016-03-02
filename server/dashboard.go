package server

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/rancher/go-rancher/api"
	"github.com/rancher/go-rancher/client"
	"github.com/rancher/v2-api/model"
	"net/http"
	"strconv"
)

func (s *Server) DashBoard(rw http.ResponseWriter, r *http.Request) error {
	hosts, err := s.getHosts(r)
	if err != nil {
		return err
	}
	stacks, err := s.getStacks(r)
	if err != nil {
		//return  err
	}
	containers, err := s.getContainers(r)
	if err != nil {
		//return  err
	}
	services, err := s.getServices(r)
	if err != nil {
		//return  err
	}
	processes, err := s.getProcesses(r)
	if err != nil {
		//return  err
	}
	auditlogs, err := s.getAuditLogs(r)
	if err != nil {
		//return  err
	}

	dashBoard := model.DashBoard{
		Hosts:      hosts,
		Stacks:     stacks,
		Containers: containers,
		Services:   services,
		Processes:  processes,
		AuditLogs:  auditlogs,
	}

	rw.Header().Set("X-API-Schemas", api.GetApiContext(r).UrlBuilder.Collection("schema"))

	if rw.Header().Get("Content-Type") == "" {
		rw.Header().Set("Content-Type", "application/json")
	}

	enc := json.NewEncoder(rw)
	return enc.Encode(dashBoard)
}

const numBuckets = 10

func (s *Server) getHosts(r *http.Request) (*model.Hosts, error) {
	hosts := []model.Host{}
	err := s.DB.Select(&hosts, `SELECT data, state FROM host`)
	if err != nil {
		logrus.Debugf("Error getting hosts %#v", err)
		return nil, err
	}
	howManyBuckets := float64(numBuckets)
	if r.URL.Query().Get("buckets") != "" {
		howManyBuckets, err = strconv.ParseFloat(r.URL.Query().Get("buckets"), 64)
		if err != nil {
			return nil, err
		}
	}
	hostsInfo := &model.Hosts{
		CPU: makeBuckets(howManyBuckets),
	}
	for _, host := range hosts {
		logrus.Debugf("Host: %#v", host)
		hostdata := model.HostData{}
		err = json.Unmarshal([]byte(host.Data), &hostdata)
		if err != nil {
			return nil, err
		}
		hostdata.FormattedID = s.obfuscate(r, "host", strconv.FormatInt(host.ID, 10))
		if err = hostToHostInfo(&hostdata, hostsInfo); err != nil {
			return nil, err
		}
	}
	return hostsInfo, nil
}

func hostToHostInfo(host *model.HostData, hostsInfo *model.Hosts) error {
	for _, val := range host.Fields.Info.CPUInfo.CPUCoresPercentages {
		hostsInfo.CPU.AddValue(val, host.FormattedID)
	}
	return nil
}

func makeBuckets(numBuckets float64) *model.HostsBucket {
	buckets := []*model.BucketOfIDs{}
	sizeOfBucket := float64(100) / numBuckets
	for x := float64(0); x < numBuckets; x++ {
		buckets = append(buckets, &model.BucketOfIDs{
			IDs:        []string{},
			RangeStart: x * sizeOfBucket,
			RangeEnd:   (x + 1) * sizeOfBucket,
		})
	}
	return &model.HostsBucket{
		Buckets: buckets,
	}
}

func (s *Server) getStacks(r *http.Request) (*model.Stacks, error) {
	return nil, fmt.Errorf("Implement this method.")
}
func (s *Server) getContainers(r *http.Request) (*model.Containers, error) {
	return nil, fmt.Errorf("Implement this method.")
}
func (s *Server) getServices(r *http.Request) (*model.Services, error) {
	return nil, fmt.Errorf("Implement this method.")
}
func (s *Server) getProcesses(r *http.Request) (*model.Processes, error) {
	return nil, fmt.Errorf("Implement this method.")
}
func (s *Server) getAuditLogs(r *http.Request) ([]client.AuditLog, error) {
	return nil, fmt.Errorf("Implement this method.")
}
