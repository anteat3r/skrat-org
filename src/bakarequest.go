package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

func BakaQuery(
  app *pocketbase.PocketBase,
  user *core.Record,
  method, endpoint, body string,
) (status int, res string, err error) {

  defer func(){
    if err == nil { return }
    user.Set(BAKAVALID, false)
    err = app.Save(user)
  }()

  attempts := 0
  var resp *http.Response
  var req *http.Request
  var bodybuf *strings.Reader
  var resb []byte
  var nw time.Time

  if user.GetDateTime(BAKATOKEN_EXPIRES).Time().Before(time.Now()){
    goto try_refresh
  }

  try_access:
    bodybuf = strings.NewReader(body)
    app.Logger().Info(BAKA_PATH + "api/v3/" + endpoint)
    req, err = http.NewRequest(
      method,
      BAKA_PATH + "api/v3/" + endpoint,
      bodybuf,
    )
    if err != nil { return }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization", "Bearer " + user.GetString(BAKATOKEN))

    nw = time.Now()
    resp, err = http.DefaultClient.Do(req)
    if err != nil { return }
    app.Logger().Info(time.Since(nw).String())

    if resp.StatusCode != 401 {
      resb, err = io.ReadAll(resp.Body)
      if err != nil { return }
      resp.Body.Close()
      return resp.StatusCode, string(resb), nil
    }

  try_refresh:
    req, err = http.NewRequest(
      "POST",
      BAKA_PATH + "api/login",
      strings.NewReader("client_id=ANDR&grant_type=refresh_token&refresh_token=" + user.GetString(BAKAREFRESTOKEN)),
    )
    if err != nil { return }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    resp, err = http.DefaultClient.Do(req)

    if resp.StatusCode == 200 {
      resb, err = io.ReadAll(resp.Body)
      if err != nil { return }
      resp.Body.Close()
      jres := BakaLoginResponse{}
      err = json.Unmarshal(resb, &jres)
      if err != nil { return }

      user.Set(BAKATOKEN, jres.AccessToken)
      user.Set(BAKAREFRESTOKEN, jres.RefreshToken)
      var date types.DateTime
      date, err = types.ParseDateTime(time.Now().Add(time.Second * time.Duration(jres.ExpiresIn)))
      if err != nil { return }
      user.Set(BAKATOKEN_EXPIRES, date)

      err = app.Save(user)
      if err != nil { return }

      if attempts > 0 {
        err = fmt.Errorf("invalid token for user %v\n", user.Id)
        return
      }
      attempts++
      goto try_access
    }

  return 0, "", nil
}

type BakaLoginResponse struct {
  UserId string `json:"bak:UserId"`
  AccessToken string `json:"access_token"`
  TokenType string `json:"token_type"`
  ExpiresIn int `json:"expires_in"`
  Scope string `json:"scope"`
  RefreshToken string `json:"refresh_token"`
  ApiVersion string `json:"bak:ApiVersion"`
  AppVersion string `json:"bak:AppVersion"`
}

func BakaLoginPass(
  app *pocketbase.PocketBase,
  user *core.Record,
  username, password string,
) (err error) {
  req, err := http.NewRequest(
    "POST",
    BAKA_PATH + "/api/login",
    strings.NewReader("client_id=ANDR&grant_type=password&username=" + username + "&password=" + password),
  )
  if err != nil { return err }
  req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
  resp, err := http.DefaultClient.Do(req)
  if err != nil { return err }

  body, err := io.ReadAll(resp.Body)
  if err != nil { return err }

  res := BakaLoginResponse{}
  err = json.Unmarshal(body, &res)
  if err != nil { return err }

  user.Set(BAKATOKEN, res.AccessToken)
  user.Set(BAKAREFRESTOKEN, res.RefreshToken)
  
  date, _ := types.ParseDateTime(time.Now().Add(
    time.Second * time.Duration(res.ExpiresIn)))
  user.Set(BAKATOKEN_EXPIRES, date)

  user.Set(BAKAVALID, true)

  err = app.Save(user)

  return
}
