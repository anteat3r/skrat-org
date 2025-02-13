package src

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

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

    var res map[string]any
    err = json.Unmarshal([]byte(resp), &res)
    if err != nil { return err }

    err = StoreData(
      app, datacoll,
      endp, PRIVATE,
      user.Id,
      res, resp,
    )
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

    if time == GetTTime() {
      err = StoreData(
        app,
        datacoll,
        name, ttype, "",
        parsedtt, stringtt,
      )
      if err != nil { return err }
    }

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
  datacoll *core.Collection,
) func(*core.RequestEvent) error {
  return func(e *core.RequestEvent) error {
    weekdays := e.Request.URL.Query().Get("day")
    weekday, err := strconv.Atoi(weekdays)
    if err != nil { return err }

    if weekday < 1 || weekday > 5 { return e.Error(401, "invalid day param", weekday) }

    ttype := e.Request.PathValue("ttype")
    if ttype != TEACHER && ttype != CLASS && ttype != ROOM {
      return e.Error(401, "invalid ttype", ttype)
    }

    classsrcs, err := app.FindRecordsByFilter(
      SOURCES,
      TYPE + ` = "` + ttype + `"`,
      `created`, -1, 0,
    )
    if err != nil { return err }

    res := struct{
      Data map[string]TimeTableDay `json:"data"`
      Hours []TimeTableHourTitle `json:"hours"`
    }{ Data: make(map[string]TimeTableDay), }

    if len(classsrcs) < 1 { return e.JSON(200, res) }

    for _, classsrc := range classsrcs {
      tt, err := QueryData[TimeTable](
        app,
        classsrc.GetString(NAME),
        ttype, "",
      )
      if err != nil {
        if err == sql.ErrNoRows { continue }
        return err
      }

      if res.Hours == nil { res.Hours = tt.Hours }

      if len(tt.Days) < weekday { continue }
      day := tt.Days[weekday - 1]

      res.Data[classsrc.GetString(DESC)] = day
    }

    return e.JSON(200, res)
  }
}
