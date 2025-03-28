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

var doNotRedirectClient = http.Client{
  CheckRedirect: func(
    req *http.Request, via []*http.Request,
  ) error {
    return http.ErrUseLastResponse
  },
}

var _ error = (*BakaInvalidError)(nil)

type BakaInvalidError struct {
  user string
  t time.Time
}

func (e BakaInvalidError) Error() string {
  return fmt.Sprintf("baka cookie expired for user %s (%v)", e.user, e.t)
}

func BakaQuery(
  app core.App,
  user *core.Record,
  method, endpoint, body string,
) (res []byte, err error) {

  // defer func(){
  //   if err == nil { return }
  //   user.Set(BAKAVALID, false)
  //   fmt.Printf("%v %#v %T\n", err, err, err)
  //   err = app.Save(user)
  //   if err != nil { return }
  // }()

  attempts := 0
  var resp *http.Response
  var req *http.Request
  var bodybuf *strings.Reader
  var resb []byte

  if user.GetDateTime(BAKATOKEN_EXPIRES).Time().Before(time.Now()){
    goto try_refresh
  }

  try_access:
    bodybuf = strings.NewReader(body)
    req, err = http.NewRequest(
      method,
      BAKA_PATH + "api/3/" + endpoint,
      bodybuf,
    )
    if err != nil { return }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("Authorization", "Bearer " + user.GetString(BAKATOKEN))

    resp, err = http.DefaultClient.Do(req)
    if err != nil { return }

    if resp.StatusCode != 401 {
      resb, err = io.ReadAll(resp.Body)
      if err != nil { return }
      resp.Body.Close()

      if resp.StatusCode != 200 {
        return nil, fmt.Errorf("invalid status code: %v %v", resp.StatusCode, string(resb))
      }
      return resb, nil
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
    if err != nil { return }

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
        user.Set(BAKAVALID, false)
        err = app.Save(user)
        if err != nil { return }
        err = fmt.Errorf("invalid token for user %v\n", user.Id)
        return
      }
      attempts++
      goto try_access
    }

  return nil, fmt.Errorf("idk bro tohle by se stát nemělo")
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
    BAKA_PATH + "api/login",
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

  req, err = http.NewRequest(
    "GET",
    BAKA_PATH + "login",
    strings.NewReader(""),
  )
  if err != nil { return err }

  resp, err = http.DefaultClient.Do(req)
  if err != nil { return err }

  rawcookie := resp.Header.Get("Set-Cookie")
  rwckl := strings.Split(rawcookie, ";")
  if len(rwckl) < 2 { return fmt.Errorf("invalid cookie: %v\n", rawcookie) }
  cookie := rwckl[0]

  payl := "username=" + username + "&password=" + password + "&persistent=true&returnUrl="

  req, err = http.NewRequest(
    "POST",
    BAKA_PATH + "Login",
    strings.NewReader(payl),
  )
  if err != nil { return err }

  req.Header.Add("Cookie", cookie)
  req.Header.Add(
    "Content-Type",
    "application/x-www-form-urlencoded",
  )
  req.Header.Add(
    "Accept", 
    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/png,image/svg+xml,*/*;q=0.8",
  )

  resp, err = doNotRedirectClient.Do(req)
  if err != nil { return err }

  rawcookie = resp.Header.Get("Set-Cookie")
  rwckl = strings.Split(rawcookie, ";")
  if len(rwckl) < 2 { return fmt.Errorf("invalid cookie: %v\n", rawcookie) }
  cookie = rwckl[0]

  user.Set(BAKACOOKIE, cookie)
  nextweek, err := types.ParseDateTime(
    time.Now().Add(time.Hour * time.Duration(24 * 6.5)),
  )
  if err != nil { return err }
  user.Set(BAKACOOKIE_EXPIRES, nextweek)
  
  user.Set(BAKAVALID, true)

  err = app.Save(user)

  return
}

func BakaWebQuery(
  app core.App,
  user *core.Record,
  endpoint string,
) (res string, err error) {
  if user.GetDateTime(BAKACOOKIE_EXPIRES).Time().Before(time.Now()) {
    // user.Set(BAKAVALID, false)
    // err = app.Save(user)
    // if err != nil { return }
    // err = fmt.Errorf("user cookie expired %s", user.GetString(NAME))
    err = BakaInvalidError{user: user.GetString(NAME), t: time.Now()}
    return
  }
  
  var req *http.Request
  req, err = http.NewRequest(
    "GET",
    BAKA_PATH + endpoint,
    nil,
  )
  if err != nil { return }
  req.Header.Set("Cookie", user.GetString(BAKACOOKIE))

  var resp *http.Response
  resp, err = http.DefaultClient.Do(req)
  if err != nil { return }

  var resb []byte
  resb, err = io.ReadAll(resp.Body)
  if err != nil { return }

  if resp.StatusCode != 200 {
    err = fmt.Errorf("invalid status code %v %v", resp.StatusCode, resp)
    return
  }

  return string(resb), nil
}

func BakaTimeTableQuery(
  app core.App,
  user *core.Record,
  time, ttype, name string,
) (tt TimeTable, err error) {
  defer func(){
    if err != nil {
      user.Set(BAKAVALID, false)
      errs := app.Save(user)
      if errs != nil { err = errs }
    }
  }()
  resp, err := BakaWebQuery(
    app, user,
    TIMETABLE_PUBLIC + "/" + time + "/" + ttype + "/" + name,
  )
  if err != nil { return TimeTable{}, err }

  res, err := ParseTimeTableWeb(resp)
  if err != nil { return TimeTable{}, err }

  return res, nil
}
