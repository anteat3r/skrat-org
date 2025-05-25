package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
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

      if ok {
        for i, day := range parsedtt.Days {
					dday := time.Now()
					rdayt := day.Title
					rdayst1 := strings.Split(rdayt, "\n")
					if len(rdayst1) == 2 {
						rdst2 := strings.Split(rdayst1[1], ".")
						if len(rdst2) == 3 {
							dn, err := strconv.Atoi(rdst2[0])
							if err != nil { return err }
							mn, err := strconv.Atoi(rdst2[1])
							if err != nil { return err }
							dday = time.Date(dday.Year(), time.Month(mn), dn, 10, 0, 0, 0, time.UTC)
						}
					}
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
        continue 
      }
      day := tt.Days[weekday - 1]

      day.Owner = classsrc.GetString(DESC)

      res.Data[classsrc.GetString(NAME)] = day
    }

    evts, ok, err := QueryData[BakaEvents](app, EVENTS, EVENTS, "")
    if err != nil { return err }

		rdayn := ""
		for k := range res.Data {
			rdayn = k
			break
		}

		dday := time.Now()

		rdayt := res.Data[rdayn].Title
		rdayst1 := strings.Split(rdayt, "\n")
		if len(rdayst1) == 2 {
			rdst2 := strings.Split(rdayst1[1], ".")
			if len(rdst2) == 3 {
				dn, err := strconv.Atoi(rdst2[0])
				if err != nil { return err }
				mn, err := strconv.Atoi(rdst2[1])
				if err != nil { return err }
				dday = time.Date(dday.Year(), time.Month(mn), dn, 10, 0, 0, 0, time.UTC)
			}
		}

    if ok {
      for _, e := range evts.Events {
        if !e.ContainsDay(dday) { continue }
        if ttype == TEACHER {
          for _, teacher := range e.Teachers {
            tday, ok := res.Data[teacher.Id]
            if !ok { continue }
            if tday.JoinedEvents == nil {
              tday.JoinedEvents = []BakaEvent{ e }
            } else {
              tday.JoinedEvents = append(tday.JoinedEvents, e)
            }
            res.Data[teacher.Id] = tday
          }
        }
        if ttype == CLASS {
          for _, class := range e.Classes {
            tday, ok := res.Data[class.Id]
            if !ok { continue }
            if tday.JoinedEvents == nil {
              tday.JoinedEvents = []BakaEvent{ e }
            } else {
              tday.JoinedEvents = append(tday.JoinedEvents, e)
            }
            res.Data[class.Id] = tday
          }
        }
        if ttype == ROOM {
          for _, room := range e.Rooms {
            tday, ok := res.Data[room.Id]
            if !ok { continue }
            if tday.JoinedEvents == nil {
              tday.JoinedEvents = []BakaEvent{ e }
            } else {
              tday.JoinedEvents = append(tday.JoinedEvents, e)
            }
            res.Data[room.Id] = tday
          }
        }
      }
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
    qp := e.Request.URL.Query()

    teacher := qp.Get("teacher")
    class := qp.Get("class")
    room := qp.Get("room")
    str := qp.Get("string")
    student := qp.Get("student")
    date := qp.Get("date")
    var tdate time.Time
    if date != "" {
      var err error
      tdate, err = time.Parse("2006-01-02", date)
      if err != nil { return err }
    }
    cached := qp.Get("cached") != ""

    etype := qp.Get(TYPE)

    var events BakaEvents
    if cached {
      switch etype {
      case "all":
        evts, ok, err := QueryData[BakaEvents](app, EVENTS, EVENTS, "")
        if err != nil { return err }
        if !ok { return fmt.Errorf("data not cached") }
        events = evts
      case "public":
        return fmt.Errorf("data not cached")
      case "my":
        evts, ok, err := QueryData[BakaEvents](app, EVENTS_MY, PRIVATE, user.Id)
        if err != nil { return err }
        if !ok { return fmt.Errorf("data not cached") }
        events = evts
      }
    } else {
      switch etype {
      case "all":
        resp, err := BakaQuery(app, user, GET, EVENTS, "")
        if err != nil { return err }
        sresp := string(resp)

        err = json.Unmarshal(resp, &events)
        if err != nil { return err }

        err = StoreData(
          app, datacoll,
          EVENTS, EVENTS, "",
          events, sresp,
        )
        if err != nil { return err }
      case "public":
        resp, err := BakaQuery(app, user, GET, EVENTS_PUBLIC, "")
        if err != nil { return err }

        err = json.Unmarshal(resp, &events)
        if err != nil { return err }
      case "my":
        resp, err := BakaQuery(app, user, GET, EVENTS_MY, "")
        if err != nil { return err }
        sresp := string(resp)

        err = json.Unmarshal(resp, &events)
        if err != nil { return err }

        err = StoreData(
          app, datacoll,
          EVENTS_MY, PRIVATE, user.Id,
          events, sresp,
        )
        if err != nil { return err }
      }
    }


    fevts := BakaEvents{Events: make([]BakaEvent, 0, len(events.Events))}
    for _, evt := range events.Events {
      if teacher != "" {
        if !slices.ContainsFunc(evt.Teachers, func(e BakaIdExpand) bool {
          return e.Id == teacher || e.Abbrev == teacher || strings.Contains(e.Name, teacher)
        }) { continue }
      }
      if class != "" {
        if !slices.ContainsFunc(evt.Classes, func(e BakaIdExpand) bool {
          return e.Id == class || e.Abbrev == class || strings.Contains(e.Name, class)
        }) { continue }
      }
      if room != "" {
        if !slices.ContainsFunc(evt.Rooms, func(e BakaIdExpand) bool {
          return e.Id == room || e.Abbrev == room || strings.Contains(e.Name, room)
        }) { continue }
      }
      if student != "" {
        if !slices.ContainsFunc(evt.Students, func(e BakaIdExpand) bool {
          return e.Id == student || e.Abbrev == student || strings.Contains(e.Name, student)
        }) { continue }
      }
      if str != "" {
        if !strings.Contains(evt.Title, str) &&
           !strings.Contains(evt.Description, str) &&
           !strings.Contains(evt.Note, str) {
          continue
        }
      }
      if date != "" {
        if !evt.ContainsDay(tdate) { continue }
      }
      fevts.Events = append(fevts.Events, evt)
    }

    fresp, err := json.Marshal(fevts)
    if err != nil { return err }

    return e.String(200, string(fresp))
  }
}

func MyTimeTableHandler(
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

    var ttable BakaTimeTable
    err = json.Unmarshal(resp, &ttable)
    if err != nil { return err }

    err = StoreData(
      app, datacoll,
      TIMETABLE_ACTUAL, PRIVATE, user.Id,
      ttable, sresp,
    )
    if err != nil { return err }

		// evts, ok, err := QueryData[BakaEvents](
		// 	app, EVENTS_MY, PRIVATE, user.Id,
		// )
		// if err != nil { return err }
		// if ok {
		//
		// }

    return e.String(200, sresp)
  }
}

func XKCDHandler(e *core.RequestEvent) error {
	req, err := http.NewRequest( "GET", "https://xkcd.com", nil)
	if err != nil { return err }

	resp, err := http.DefaultClient.Do(req)
	if err != nil { return err }

	resb, err := io.ReadAll(resp.Body)
	if err != nil { return err }

	reg, err := regexp.Compile(`<img src="\/\/imgs\.xkcd\.com\/comics\/(.+)\.png`)
	if err != nil { return err }

	match := reg.FindStringSubmatch(string(resb))
	if match == nil || len(match) != 2 { return e.InternalServerError("idk lil bro", nil) }

	return e.String(200, "https://imgs.xkcd.com/comics/" + match[1])
}
