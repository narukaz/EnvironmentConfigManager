package connection

import (
	"fmt"
	"net/http"

	"github.com/narukaz/EnvironmentConfigManager/pkg/operation"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func ServerConnect(host string, port int, env string, client *mongo.Client) {

	var mux = http.NewServeMux()
	var handler = &operation.Handler{Client: client}

	mux.HandleFunc("/create", handler.InsertItem)
	mux.HandleFunc("/delete", handler.DeleteOne)
	mux.HandleFunc("/getall", handler.GetAllEmployee)

	var formatAddress = fmt.Sprintf("%s:%d", host, port)
	fmt.Println("connecting at \n", formatAddress)
	http.ListenAndServe(formatAddress, mux)
}
