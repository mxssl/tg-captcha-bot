package main

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
)

type CasResponse struct {
        Ok     bool `json:"ok"`
        Result struct {
                Status string `json:"status"`
        } `json:"result"`
}

func checkUserCas(userID int64) (bool, string, error) {
        resp, err := http.Get(fmt.Sprintf("https://api.cas.chat/check?user_id=%d", userID))
        if err != nil {
                return false, "", err
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                return false, "", err
        }

        var casResponse CasResponse
        err = json.Unmarshal(body, &casResponse)
        if err != nil {
                return false, "", err
        }

        return casResponse.Ok && casResponse.Result.Status == "ok", casResponse.Result.Status, nil
}
