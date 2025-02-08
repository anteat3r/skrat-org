package src

import (
	"encoding/json"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func EndpHandler(
  app *pocketbase.PocketBase,
  datacoll *core.Collection,
) func(*core.RequestEvent) error {
  return func(e *core.RequestEvent) error {
    user := e.Auth
    if !user.GetBool(BAKAVALID) {
      return e.UnauthorizedError("your baka login is not valid", nil)
    }
    endp := e.Request.URL.Query().Get("endp")
    status, resp, err := BakaQuery(app, user, "GET", endp, "")
    if err != nil { return err }

    var datarec *core.Record
    datarec, err = app.FindFirstRecordByFilter(
      DATA,
      "owner = {:owner} && name = {:name}",
      dbx.Params{"owner": user.Id, "name": endp},
    )
    if err != nil {
      datarec = core.NewRecord(datacoll)
      datarec.Set(OWNER, user.Id)
      datarec.Set(NAME, endp)
    }
    datarec.Set(DATA, resp)

    err = app.Save(datarec)
    if err != nil { return err }

    return e.String(status, resp)
  }
}

func LoginHandler(
  app *pocketbase.PocketBase,
) func(*core.RequestEvent) error {
  return func(e *core.RequestEvent) error {
    body := struct{
      Username string `json:"username"`
      Password string `json:"password"`
    }{}
    err := e.BindBody(&body)
    if err != nil { return err }

    return BakaLoginPass(app, e.Auth, body.Username, body.Password)
  }
}

func WebTimeTableHandler(
  app *pocketbase.PocketBase,
  datacoll *core.Collection,
) func(*core.RequestEvent) error {
  return func(e *core.RequestEvent) error {
    user := e.Auth
    time := e.Request.PathValue("time")
    ttype := e.Request.PathValue("ttype")
    name := e.Request.PathValue("name")
    status, res, err := BakaTimeTableQuery(app, user, time, ttype, name)
    if err != nil { return err }
    if status != 200 { return fmt.Errorf("bad status code: %v", status) }

    parsedtt, err := ParseTimeTableWeb(res)
    if err != nil { return err }

    jsontt, err := json.Marshal(parsedtt)
    if err != nil { return err }

    stringtt := string(jsontt)

    var datarec *core.Record
    datarec, err = app.FindFirstRecordByFilter(
      DATA,
      `owner = "" && name = {:name}`,
      dbx.Params{"name": name},
    )
    if err != nil {
      datarec = core.NewRecord(datacoll)
      datarec.Set(OWNER, user.Id)
      datarec.Set(NAME, name)
    }
    datarec.Set(DATA, stringtt)

    err = app.Save(datarec)
    if err != nil { return err }

    return e.String(200, stringtt)
  }
}

func WebSourcesHandler(
  app *pocketbase.PocketBase,
) func(*core.RequestEvent) error {
  return func(e *core.RequestEvent) error {
    user := e.Auth
    
    status, html, err := BakaWebQuery(app, user, "TimeTable/Public")
    if status != 200 { return fmt.Errorf("bad status code: %v", status) }

    srcs, err := ParseSourcesWeb(html)
    if err != nil { return err }

    return e.JSON(200, srcs)
  }
}
