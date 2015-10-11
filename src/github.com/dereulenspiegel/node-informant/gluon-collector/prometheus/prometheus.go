package prometheus

import (
	log "github.com/Sirupsen/logrus"
	"github.com/dereulenspiegel/node-informant/gluon-collector/data"
	stat "github.com/prometheus/client_golang/prometheus"
)

var (
	TotalClientCounter = stat.NewGauge(stat.GaugeOpts{
		Name: "total_clients",
		Help: "Total number of connected clients",
	})

	TotalNodes = stat.NewGauge(stat.GaugeOpts{
		Name: "total_nodes",
		Help: "Total number of Nodes",
	})

	TotalNodeTrafficRx = stat.NewCounter(stat.CounterOpts{
		Name: "total_traffic_rx",
		Help: "Total accumulated received traffic as reported by Nodes",
	})

	TotalNodeTrafficTx = stat.NewCounter(stat.CounterOpts{
		Name: "total_traffic_tx",
		Help: "Total accumulated transmitted traffic as reported by Nodes",
	})

	TotalNodeMgmtTrafficRx = stat.NewCounter(stat.CounterOpts{
		Name: "total_traffic_mgmt_rx",
		Help: "Total accumulated received management traffic as reported by Nodes",
	})

	TotalNodeMgmtTrafficTx = stat.NewCounter(stat.CounterOpts{
		Name: "total_traffic_mgmt_tx",
		Help: "Total accumulated transmitted management traffic as reported by Nodes",
	})
)

func init() {
	stat.MustRegister(TotalClientCounter)
	stat.MustRegister(TotalNodes)
	stat.MustRegister(TotalNodeTrafficRx)
	stat.MustRegister(TotalNodeTrafficTx)
	stat.MustRegister(TotalNodeMgmtTrafficRx)
	stat.MustRegister(TotalNodeMgmtTrafficTx)
}

func initTotalClientsGauge(store data.Nodeinfostore) {
	TotalClientCounter.Set(0.0)
	var totalClients int = 0
	for _, stats := range store.GetAllStatistics() {
		status, err := store.GetNodeStatusInfo(stats.NodeId)
		if err != nil {
			log.WithFields(log.Fields{
				"error":  err,
				"nodeId": stats.NodeId,
			}).Warn("Didn't found node status in store")
		}
		if status.Online {
			totalClients = totalClients + stats.Clients.Total
			log.Debugf("Adding %d clients", stats.Clients.Total)
			TotalClientCounter.Add(float64(stats.Clients.Total))
		} else {
			log.Debugf("Node %s was offline", status.NodeId)
		}
	}
	log.Debugf("Initialised prometheus with %d clients", totalClients)
}

func ProcessStoredValues(store data.Nodeinfostore) {
	TotalNodes.Set(float64(len(store.GetNodeStatusInfos())))
	initTotalClientsGauge(store)
}
