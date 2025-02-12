package src

import (
	"encoding/json"
	"fmt"
	"strconv"

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
    if ttype != TEACHER && ttype != CLASS && ttype != ROOM {
      return e.Error(401, "invalid ttype", ttype)
    }
    status, res, err := BakaTimeTableQuery(app, user, time, ttype, name)
    if err != nil { return err }
    if status != 200 { return fmt.Errorf("bad status code: %v", status) }

    parsedtt, err := ParseTimeTableWeb(res)
    if err != nil { return err }

    jsontt, err := json.Marshal(parsedtt)
    if err != nil { return err }

    stringtt := string(jsontt)

    if time != GetTTime() { return e.String(200, stringtt) }

    var datarec *core.Record
    datarec, err = app.FindFirstRecordByFilter(
      DATA,
      `owner = "" && name = {:name}`,
      dbx.Params{"name": name},
    )
    if err != nil {
      datarec = core.NewRecord(datacoll)
      datarec.Set(OWNER, "")
      datarec.Set(NAME, name)
      datarec.Set(TYPE, ttype)
    }
    datarec.Set(DATA, stringtt)
    if datarec.GetString(TYPE) == "" {
      datarec.Set(TYPE, ttype)
    }

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
    
    status, html, err := BakaWebQuery(app, user, TIMETABLE_PUBLIC)
    if status != 200 { return fmt.Errorf("bad status code: %v", status) }

    srcs, err := ParseSourcesWeb(html)
    if err != nil { return err }

    return e.JSON(200, srcs)
  }
}

func DayOverviewHandler(
  app *pocketbase.PocketBase,
) func(*core.RequestEvent) error {
  return func(e *core.RequestEvent) error {
    weekdays := e.Request.URL.Query().Get("day")
    weekday, err := strconv.Atoi(weekdays)
    if err != nil { return err }

    if weekday < 1 || weekday > 5 {
      return e.Error(401, "invalid day param", weekday)
    }

    classsrcs, err := app.FindRecordsByFilter(
      SOURCES,
      TYPE + ` = "` + CLASS + `"`,
      `created`,
      -1, 0,
    )
    if err != nil { return err }

    app.Logger().Info(fmt.Sprintf("%#v", classsrcs))

    res := struct{
      Data map[string][]TimeTableHour `json:"data"`
      Hours []TimeTableHourTitle `json:"hours"`
    }{
      Data: make(map[string][]TimeTableHour),
    }

    if len(classsrcs) < 1 { return e.JSON(200, res) }

    for _, classsrc := range classsrcs {
      datarecs, err := app.FindRecordsByFilter(
        DATA,
        NAME + ` = "` + classsrc.GetString(NAME) + `" && ` + OWNER + ` = ""`,
        `created`,
        1, 0,
      )
      if err != nil { return err }

      if len(datarecs) < 1 { continue }
      datarec := datarecs[0]

      var tt TimeTable
      err = json.Unmarshal([]byte(datarec.GetString(DATA)), &tt)
      if err != nil { return err }

      if len(tt.Days) < weekday - 2 { continue }

      res.Data[datarec.GetString(DESC)] = tt.Days[weekday - 1].Hours
    }

    return e.JSON(200, res)
  }
}
