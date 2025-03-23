package src

import (
	"encoding/json"
	"strconv"
	"time"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
)

var RequireBakaValid = &hook.Handler[*core.RequestEvent]{
  Id: "requireBakaValid",
  Func: func(re *core.RequestEvent) error {
    user := re.Auth
    if user == nil { return re.UnauthorizedError("not auth token", nil) }
    if !user.GetBool(BAKAVALID) && user.GetString(RANK) != SPERL { return re.UnauthorizedError("not bakavalid", nil) }
    return re.Next()
  },
}

func EndpHandler(
  app *pocketbase.PocketBase,
  datacoll *core.Collection,
) func(*core.RequestEvent) error {
  return func(e *core.RequestEvent) error {
    user := e.Auth
    endp := e.Request.URL.Query().Get("endp")
    resp, err := BakaQuery(app, user, "GET", endp, "")
    if err != nil { return err }

    return e.String(200, string(resp))
  }
}

func MarksHandler(
  app *pocketbase.PocketBase,
  datacoll *core.Collection,
) func(*core.RequestEvent) error {
  return func(e *core.RequestEvent) error {
    user := e.Auth

    resp, err := BakaQuery(app, user, GET, MARKS, "")
    if err != nil { return err }
    sresp := string(resp)

    var marks BakaMarks
    err = json.Unmarshal(resp, &marks)
    if err != nil { return err }

    err = StoreData(
      app, datacoll,
      MARKS, PRIVATE, user.Id,
      marks, sresp,
    )
    if err != nil { return err }

    return e.String(200, sresp)
  }
}

func TimeTableHandler(
  app *pocketbase.PocketBase,
  datacoll *core.Collection,
) func(*core.RequestEvent) error {
  return func(e *core.RequestEvent) error {
    user := e.Auth

    date := e.Request.URL.Query().Get("date")

    qparam := ""
    if date != "" {
      qparam = "?date=" + date
    }

    resp, err := BakaQuery(app, user, GET, TIMETABLE_ACTUAL + qparam, "")
    if err != nil { return err }
    sresp := string(resp)

    var marks BakaMark
    err = json.Unmarshal(resp, &marks)
    if err != nil { return err }

    err = StoreData(
      app, datacoll,
      TIMETABLE_ACTUAL, PRIVATE, user.Id,
      marks, sresp,
    )
    if err != nil { return err }

    return e.String(200, sresp)
  }
}

func StoreVapidEndpoint(
  app *pocketbase.PocketBase,
) func(*core.RequestEvent) error {
  return func(e *core.RequestEvent) error {
    body := struct{
      Vapid string `json:"vapid"`
    }{}
    err := e.BindBody(&body)
    if err != nil { return err }

    user := e.Auth
    user.Set(WANTS_REFRESH, true)
    user.Set(VAPID, body.Vapid)

    err = app.Save(user)
    if err != nil { return err }

    return e.String(200, "")
  }
}

func VapidTestHandler(
  app *pocketbase.PocketBase,
) func(*core.RequestEvent) error {
  return func(e *core.RequestEvent) error {
    user := e.Auth

    vapid := user.GetString(VAPID)

    s := &webpush.Subscription{}
    err := json.Unmarshal([]byte(vapid), s)
    if err != nil { return err }

    _, err = webpush.SendNotification([]byte(`{"title":"test","type":"notif"}`), s, &webpush.Options{
      Subscriber: user.GetString("email"),
      VAPIDPublicKey: VAPID_PUBKEY,
      VAPIDPrivateKey: VAPID_PRIVKEY,
    })
    if err != nil { return err }

    // resb, err := io.ReadAll(resp.Body)
    // defer resp.Body.Close()
    // if err != nil { return err }

    // app.Logger().Info(string(resb))
    // app.Logger().Info(fmt.Sprintf("%#v", s))

    return e.String(200, "")
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

    ttime := e.Request.PathValue("time")
    ttype := e.Request.PathValue("ttype")
    name := e.Request.PathValue("name")
    if ttype != TEACHER && ttype != CLASS && ttype != ROOM {
      return e.Error(401, "invalid ttype", ttype)
    }
    parsedtt, err := BakaTimeTableQuery(app, user, ttime, ttype, name)
    if err != nil { return err }

    if ttime != PERMANENT {
      evts, ok, err := QueryData[BakaEvents](app, EVENTS, EVENTS, "")
      if err != nil { return err }
      nw := time.Now()
      nweek := nw.AddDate(0, 0, -int(nw.Weekday() - 1))
      if ttime == NEXT {
        nweek = nweek.AddDate(0, 0, 7)
      }
      if ok {
        for i, day := range parsedtt.Days {
          dday := nweek.AddDate(0, 0, i) 
          for _, e := range evts.Events {
            if ttype == TEACHER && !BakaIdExpandListContainsId(e.Teachers, name) { continue }
            if ttype == CLASS && !BakaIdExpandListContainsId(e.Classes, name) { continue }
            if ttype == ROOM && !BakaIdExpandListContainsId(e.Rooms, name) { continue }
            if !e.ContainsDay(dday) { continue }
            parsedtt.Days[i].JoinedEvents = append(day.JoinedEvents, e)
          }
        }
      }
    }

    jsontt, err := json.Marshal(parsedtt)
    if err != nil { return err }

    stringtt := string(jsontt)

    if ttime == GetTTime() {
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
    
    html, err := BakaWebQuery(app, user, TIMETABLE_PUBLIC)
    if err != nil { return err }

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

    // app.Logger().Info(fmt.Sprintf("%#v", DataCache))

    for _, classsrc := range classsrcs {
      tt, ok, err := QueryData[TimeTable](
        app,
        classsrc.GetString(NAME),
        ttype, "",
      )
      if err != nil {
        return err
      }
      if !ok { continue }

      if res.Hours == nil {
        res.Hours = tt.Hours
      } else if len(tt.Hours) > len(res.Hours) {
        res.Hours = tt.Hours
      }

      if len(tt.Days) < weekday {
        // app.Logger().Info(fmt.Sprintf("%#v", tt))
        continue 
      }
      day := tt.Days[weekday - 1]

      res.Data[classsrc.GetString(DESC)] = day
    }

    return e.JSON(200, res)
  }
}

func EventsHandler(
  app *pocketbase.PocketBase,
  datacoll *core.Collection,
) func(*core.RequestEvent) error {
  return func(e *core.RequestEvent) error {
    user := e.Auth

    resp, err := BakaQuery(app, user, GET, EVENTS_ALL, "")
    if err != nil { return err }
    sresp := string(resp)

    var events BakaEvents
    err = json.Unmarshal(resp, &events)
    if err != nil { return err }

    err = StoreData(
      app, datacoll,
      EVENTS, EVENTS, "",
      events, sresp,
    )
    if err != nil { return err }

    return e.String(200, sresp)
  }
}
