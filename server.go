package main

import (
	rtctokenbuilder "github.com/AgoraIO/Tools/DynamicKey/AgoraDynamicKey/go/src/RtcTokenBuilder"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"encoding/json"
	"strconv"
)

/*
* Created By Mujuzi Moses
*/

var rtc_token string
var int_uid uint32
var channel_name string

var role_num uint32
var role rtctokenbuilder.Role

func generateRtcToken(int_uid uint32, channelName string, role rtctokenbuilder.Role) {
	
	appID := "9998def41eae49c6b7169b2c46a044f1"
	appCertificate := "00e4546550324d38af42d99ecac06ecd"
	expireTimeInSeconds := uint32(50)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	result,err := rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, int_uid, role, expireTimestamp)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Token with uid: %s\n", result)
		fmt.Printf("uid is %d\n", int_uid)
		fmt.Printf("ChannelName is %s\n", channelName)
		fmt.Printf("Role is %d\n", role)
	}
	rtc_token = result	
}

func rtc(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS");
	w.Header().Set("Access-Control-Allow-Headers", "*");

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" && r.Method != "OPTIONS" {
		http.Error(w, "Unsupported method. Please check.", http.StatusNotFound)
		return
	}

	ChannelName, err1 := mux.Vars(r)["channelName"]
	if err1 != true {
		fmt.Printf("ChannelName Failed ::: %s\n", err1)
		return
	} else {
		channel_name = ChannelName
		fmt.Printf("ChannelName Success ::: %s\n", ChannelName)
	}

	Role, err2 := mux.Vars(r)["role"]
	if err2 != true {
		fmt.Printf("Role Failed ::: %s\n", err2)
		return
	} else {
		role = rtctokenbuilder.RolePublisher
		fmt.Printf("Role Success ::: %s\n", Role)
	}

	UidStr, err3 := strconv.ParseInt(mux.Vars(r)["uid"], 10, 64)
	if err3 != nil {
		fmt.Printf("UID Failed ::: %s\n", err3)
		return
	} else {
		int_uid = 1
		fmt.Printf("UID Success ::: %d\n", UidStr)
	}

	generateRtcToken(int_uid, channel_name, role)
	errorResponse(w, rtc_token, http.StatusOK)
	log.Println(w, r)
}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["token"] = message
	resp["code"] = strconv.Itoa(httpStatusCode)
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)

}

func rootPage(w http.ResponseWriter, r *http.Request) {
	
	w.Write([]byte("This is root page"))
}

const port = "8080"

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", rootPage)
	router.HandleFunc("/rtc/{channelName}/{role}/{uid}", rtc)
	fmt.Println("Starting server at http://127.0.0.1:" +port)
	log.Fatal(http.ListenAndServe(":" +port, router))
}
