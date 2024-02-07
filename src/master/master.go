package master

import (
	"net/http"
)

func (m *MasterService) Master() {
	http.HandleFunc("/clusterState", m.CS.HandleClusterState)
	http.ListenAndServe(":1986", nil)

	// More Master logic
}
