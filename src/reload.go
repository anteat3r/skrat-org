package src

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type ResourceName struct {
  Name string
  Type string
  Owner string
}

var (
  DataCache = map[ResourceName]any{}
  DataCacheMu = sync.RWMutex{}
)

func StoreData(
  app core.App,
  datacoll *core.Collection,
  name, ttype, owner string,
  parseddata any,
  stringdata string,
) error {

  var datarecs []*core.Record
  err := app.RecordQuery(DATA).Where(dbx.HashExp{OWNER: owner, NAME: name, TYPE: ttype},).Limit(1).All(&datarecs)
  if err != nil { return err }

  var datarec *core.Record
  if len(datarecs) > 0 {
    datarec = datarecs[0]
  } else {
    datarec = core.NewRecord(datacoll)
    datarec.Set(OWNER, owner)
    datarec.Set(NAME, name)
    datarec.Set(TYPE, ttype)
  }

  if stringdata == "" {
    bytedata, err := json.Marshal(parseddata)
    if err != nil { return err }
    stringdata = string(bytedata)
  }
  
  datarec.Set(DATA, stringdata)
  if datarec.GetString(TYPE) == "" {
    datarec.Set(TYPE, ttype)
  }

  err = app.Save(datarec)
  if err != nil { return err }

  DataCacheMu.Lock()
  DataCache[ResourceName{
    Name: name,
    Type: ttype,
    Owner: owner,
  }] = parseddata
  DataCacheMu.Unlock()

  return nil
}

func QueryData[T any](
  app core.App,
  name, ttype, owner string,
) (data T, dok bool, err error) {
  resname := ResourceName{
    Name: name,
    Type: ttype,
    Owner: owner,
  }
  DataCacheMu.RLock()
  cdata, ok := DataCache[resname]
  DataCacheMu.RUnlock()

  if ok {
    tdata, ok := cdata.(T)
    if ok { return tdata, true, nil }
  }

  var datarecs []*core.Record
  err = app.RecordQuery(DATA).Where(dbx.HashExp{NAME: name, TYPE: ttype, OWNER: owner}).Limit(1).All(&datarecs)
  if err != nil { return }
  if len(datarecs) < 1 {
    dok = false
    return
  }
  datarec := datarecs[0]

  rest := new(T)
  err = json.Unmarshal([]byte(datarec.GetString(DATA)), rest)
  if err != nil { return }

  DataCacheMu.Lock()
  DataCache[resname] = *rest
  DataCacheMu.Unlock()

  return *rest, true, nil
}

func UnCacheData(
  app core.App,
  name, ttype, owner string,
) {
  DataCacheMu.Lock()
  delete(DataCache, ResourceName{Name: name, Type: ttype, Owner: owner})
  DataCacheMu.Unlock()
}

func TimeTableReload(app *pocketbase.PocketBase, datacoll *core.Collection) func() {
  return func() {
    defer func(){
      r := recover()
      if r == nil { return }
      app.Logger().Error(fmt.Sprint(r))
    }()
    err := app.RunInTransaction(func(txApp core.App) error {
      srcs, err := txApp.FindRecordsByFilter(
        SOURCES,
        `id != ""`,
        LAST_UPDATED, 1, 0, 
      )
      if err != nil { return err }

      if len(srcs) < 1 { return fmt.Errorf("no suitable source found") }
      src := srcs[0]

      users, err := txApp.FindRecordsByFilter(
        USERS,
        BAKAVALID + ` = true`,
        LAST_USED, 1, 0,
      )
      if err != nil { return err }

      if len(users) < 1 { return fmt.Errorf("no suitable user found") }
      user := users[0]

      var jresp string
      var tresp any

      if src.GetString(TYPE) == EVENTS {
        resp, err := BakaQuery(txApp, user, "GET", EVENTS_ALL, "")
        if err != nil {
          if _, ok := err.(BakaInvalidError); !ok { return err }
          user.Set(BAKAVALID, false)
          err = txApp.Save(user)
          if err != nil { return err }
          SendNotifs(app, user, []Notif{ BakaInvalidNotif{} })
          return nil
        }

        jresp = string(resp)

        evts := BakaEvents{}
        err = json.Unmarshal(resp, &evts)
        if err != nil {
          if _, ok := err.(*json.SyntaxError); ok {
            app.Logger().Info(fmt.Sprintf("%#v", jresp))
            app.Logger().Info("cant laod evetns")
          }
          return err
        }

        tresp = evts
      } else {
        app.Logger().Info(fmt.Sprintf("ttreload %s %s %s", GetTTime(), src.GetString(TYPE), src.GetString(NAME)))
        tt, err := BakaTimeTableQuery(txApp, user, GetTTime(), src.GetString(TYPE), src.GetString(NAME))
        if err != nil {
          if _, ok := err.(BakaInvalidError); !ok { return err }
          user.Set(BAKAVALID, false)
          err = txApp.Save(user)
          if err != nil { return err }
          SendNotifs(app, user, []Notif{ BakaInvalidNotif{} })
          return nil
        }
        if !user.GetBool(BAKAVALID) && user.GetString(VAPID) != "" {
          err = SendNotifs(app, user, []Notif{ BakaInvalidNotif{} })
          if err != nil { return err }
          return nil
        }

        tresp = tt

        resb, err := json.Marshal(tt)
        if err != nil { return err }

        jresp = string(resb)
      }

      err = StoreData(
        txApp,
        datacoll,
        src.GetString(NAME),
        src.GetString(TYPE),
        "",
        tresp, jresp,
      )
      if err != nil { return err }

      user.Set(LAST_USED, types.NowDateTime())

      err = txApp.Save(user)
      if err != nil { return err }

      src.Set(LAST_UPDATED, types.NowDateTime())

      err = txApp.Save(src)
      if err != nil { return err }

      return nil
    })
    
    if err != nil { app.Logger().Error(err.Error(), slog.Any("error", err)) }
  }
}

func TimeTableSourcesReload(
  app *pocketbase.PocketBase,
) func() {
  return func() {

    err := app.RunInTransaction(func(txApp core.App) error {
      _, err := txApp.DB().NewQuery("delete from " + SOURCES).Execute()
      if err != nil { return err }

      users, err := txApp.FindRecordsByFilter(
        USERS,
        BAKAVALID + ` = true`,
        LAST_USED, 1, 0,
      )
      if err != nil { return err }

      if len(users) < 1 { return fmt.Errorf("no suitable user found") }
      user := users[0]

      resp, err := BakaWebQuery(txApp, user, TIMETABLE_PUBLIC)
      if err != nil {
				ierr, ok := err.(BakaInvalidError)
				if !ok { return err }
        user.Set(BAKAVALID, false)
        err = txApp.Save(user)
        if err != nil { return err }
        SendNotifs(txApp, user, []Notif{ BakaInvalidNotif{} })
        return ierr
      }

      srcs, err := ParseSourcesWeb(resp)
      if err != nil { return err }
			fmt.Println("10207")

      coll, _ := txApp.FindCollectionByNameOrId(SOURCES)
      for name, src := range srcs.AsMap() {
        for _, ttsrc := range src {
          rec := core.NewRecord(coll)
          rec.Set(NAME, ttsrc.Id)
          rec.Set(DESC, ttsrc.Name)
          rec.Set(TYPE, name)
          err := txApp.Save(rec)
          if err != nil { return err }
        }
      }

      rec := core.NewRecord(coll)
      rec.Set(NAME, EVENTS)
      rec.Set(DESC, EVENTS)
      rec.Set(TYPE, EVENTS)
      err = txApp.Save(rec)
      if err != nil { return err }

      user.Set(LAST_USED, types.NowDateTime())

      err = txApp.Save(user)
      if err != nil { return err }

      return nil
    })

    if err != nil { app.Logger().Error(err.Error(), slog.Any("error", err)) }
  }
}

func SendNotifs(
  app core.App,
  user *core.Record,
  notifs []Notif,
) error {
  app.Logger().Info(fmt.Sprintf("sending notifs to %s: %#v", user.GetString(NAME), notifs))
  vapid := user.GetString(VAPID)
  if vapid == "" { return nil }

  s := &webpush.Subscription{}
  err := json.Unmarshal([]byte(vapid), s)
  if err != nil { return err }
  for _, n := range notifs {
    resp, err := webpush.SendNotification([]byte(n.JSONEncode()), s, &webpush.Options{
      Subscriber: user.GetString("email"),
      VAPIDPublicKey: VAPID_PUBKEY,
      VAPIDPrivateKey: VAPID_PRIVKEY,
    })
    if resp != nil {
      bd, err := io.ReadAll(resp.Body)
      defer resp.Body.Close()
      if err != nil { return err }
      app.Logger().Info(fmt.Sprintf("sent notifs user %s", user.GetString(NAME)), slog.String("resp.Body", string(bd)), slog.Int("resp.StatusCode", resp.StatusCode), slog.String("body", n.JSONEncode()))
    }
    if err != nil { return err }
  }
  return nil
}

func PersonalReload(
  app core.App,
  datacoll *core.Collection,
) func() {
  return func() {

    err := func() error {


      users, err := app.FindAllRecords(
        USERS,
        dbx.HashExp{WANTS_REFRESH: true},
      )
      if err != nil { return err }

      for _, user := range users {

        if !user.GetDateTime(LAST_REFRESHED).IsZero() {
          if time.Since(user.GetDateTime(LAST_REFRESHED).Time()).Minutes() < 
            float64(user.GetInt(REFRESH_INTERVAL)) { continue }
        }

        app.Logger().Info(fmt.Sprintf("reloading user %s", user.GetString(NAME)))

        total_notifs := make([]Notif, 0)
        var marks BakaMarks
        var oldmarks BakaMarks
        var ttable BakaTimeTable
        var oldttable BakaTimeTable
        var ok bool
        var sresp string

        resp, err := BakaQuery(app, user, GET, MARKS, "")
        if err != nil {
          if _, ok := err.(BakaInvalidError); !ok { return err }
          user.Set(BAKAVALID, false)
          err = app.Save(user)
          if err != nil { return err }
          total_notifs = append(total_notifs, BakaInvalidNotif{})
          goto sendnotifs
        }
        sresp = string(resp)

        err = json.Unmarshal(resp, &marks)
        if err != nil { return err }

        oldmarks, ok, err = QueryData[BakaMarks](app, MARKS, PRIVATE, user.Id)
        if err != nil { return err }
        
        if ok {
          notifs := CompareBakaMarks(oldmarks, marks)
          app.Logger().Info(fmt.Sprintf("notifs marks user %s: %#v", user.GetString(NAME), notifs), slog.Any("oldmarks", oldmarks), slog.Any("marks", marks))
          total_notifs = append(total_notifs, notifs...)
        }

        err = StoreData(
          app, datacoll,
          MARKS, PRIVATE, user.Id,
          marks, sresp,
        )
        if err != nil { return err }

        resp, err = BakaQuery(app, user, GET, TIMETABLE_ACTUAL, "")
        if err != nil {
          if _, ok := err.(BakaInvalidError); !ok { return err }
          user.Set(BAKAVALID, false)
          err = app.Save(user)
          if err != nil { return err }
          total_notifs = append(total_notifs, BakaInvalidNotif{})
          goto sendnotifs
        }
        sresp = string(resp)

        err = json.Unmarshal(resp, &ttable)
        if err != nil { return err }

        oldttable, ok, err = QueryData[BakaTimeTable](app, TIMETABLE_ACTUAL, PRIVATE, user.Id)
        if err != nil { return err }
        
        if ok {
          if len(oldttable.Days) > 0 && len(ttable.Days) > 0 && oldttable.Days[0].Date == ttable.Days[0].Date {
            notifs := CompareBakaTimeTables(oldttable, ttable)
            total_notifs = append(total_notifs, notifs...)
          }
        }

        err = StoreData(
          app, datacoll,
          TIMETABLE_ACTUAL, PRIVATE, user.Id,
          ttable, sresp,
        )
        if err != nil { return err }

        sendnotifs:
        if len(total_notifs) > 0 {
          err := SendNotifs(app, user, total_notifs)
          if err != nil { return err }
        }

        user.Set(LAST_REFRESHED, types.NowDateTime())
        err = app.Save(user)
        if err != nil { return err }
      }

      return nil
    }()
    
    if err != nil { app.Logger().Error(err.Error(), slog.Any("error", err)) }
  }
}

func EveningRefresh(
  app core.App,
  datacoll *core.Collection,
) func() {
  return func() {
    err := func() error {

      users, err := app.FindAllRecords(
        USERS,
        dbx.HashExp{WANTS_REFRESH: true},
      )
      if err != nil { return err }

      for _, user := range users {
        if !user.GetBool(WANTS_REFRESH) { continue }

        total_notifs := make([]Notif, 0)

        var absence BakaAbsence
				var events BakaEvents
        var sresp string

        resp, err := BakaQuery(app, user, GET, ABSENCE_STUDENT, "")
        if err != nil {
          if _, ok := err.(BakaInvalidError); !ok { return err }
          user.Set(BAKAVALID, false)
          err = app.Save(user)
          if err != nil { return err }
          total_notifs = append(total_notifs, BakaInvalidNotif{})
          goto sendnotifs
        }
        sresp = string(resp)

        err = json.Unmarshal(resp, &absence)
        if err != nil { return err }

        err = StoreData(
          app, datacoll,
          ABSENCE_STUDENT, PRIVATE, user.Id,
          absence, sresp,
        )
        if err != nil { return err }

        resp, err = BakaQuery(app, user, GET, EVENTS_MY, "")
        if err != nil {
          if _, ok := err.(BakaInvalidError); !ok { return err }
          user.Set(BAKAVALID, false)
          err = app.Save(user)
          if err != nil { return err }
          total_notifs = append(total_notifs, BakaInvalidNotif{})
          goto sendnotifs
        }
        sresp = string(resp)

        err = json.Unmarshal(resp, &events)
        if err != nil { return err }

        err = StoreData(
          app, datacoll,
          EVENTS_MY, PRIVATE, user.Id,
          events, sresp,
        )
        if err != nil { return err }
        
        total_notifs = append(total_notifs, FindPendingAbsences(absence)...)

        sendnotifs:
        if len(total_notifs) > 0 {
          err := SendNotifs(app, user, total_notifs)
          if err != nil { return err }
        }

        user.Set(LAST_REFRESHED, types.NowDateTime())
        err = app.Save(user)
        if err != nil { return err }
      }

      return nil
    }()

    if err != nil { app.Logger().Error(err.Error(), slog.Any("error", err)) }
  }
}
