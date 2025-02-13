package src

import (
	"encoding/json"
	"fmt"

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
)

func StoreData(
  app core.App,
  datacoll *core.Collection,
  name, ttype, owner string,
  parseddata any,
  stringdata string,
) error {
  datarec, err := app.FindFirstRecordByFilter(
    DATA,
    OWNER + ` = {:owner} && ` + NAME + ` = {:name} && ` + TYPE + ` = {:type}`,
    dbx.Params{"name": name, TYPE: ttype, OWNER: owner},
  )
  if err != nil {
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

  DataCache[ResourceName{
    Name: name,
    Type: ttype,
    Owner: owner,
  }] = parseddata

  return nil
}

func QueryData[T any](
  app core.App,
  name, ttype, owner string,
) (data T, err error) {
  resname := ResourceName{
    Name: name,
    Type: ttype,
    Owner: owner,
  }
  cdata, ok := DataCache[resname]

  if ok {
    tdata, ok := cdata.(T)
    if ok { return tdata, nil }
  }

  var rec *core.Record
  rec, err = app.FindFirstRecordByFilter(
    DATA,
    OWNER + ` = {:owner} && ` + NAME + ` = {:name} && ` + TYPE + ` = {:type}`,
    dbx.Params{"name": name, TYPE: ttype, OWNER: owner},
  )
  if err != nil { return }

  sdata := rec.GetString(DATA)

  rest := new(T)
  err = json.Unmarshal([]byte(sdata), rest)
  if err != nil { return }

  DataCache[resname] = *rest

  return *rest, nil
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
        status, resp, err := BakaQuery(txApp, user, "GET", "events", "")
        if err != nil { return err }
        if status != 200 { return fmt.Errorf("invalid status code: %v %v", status, resp) }

        jresp = resp

        evts := BakaEvents{}
        json.Unmarshal([]byte(resp), &evts)

        tresp = evts
      } else {
        status, resp, err := BakaTimeTableQuery(txApp, user, GetTTime(), src.GetString(TYPE), src.GetString(NAME))
        if err != nil { return err }
        if status != 200 { return fmt.Errorf("invalid status code: %v %v", status, resp) }

        tt, err := ParseTimeTableWeb(resp)
        if err != nil { return err }

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
    
    if err != nil { app.Logger().Error(err.Error(), err) }
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

      status, resp, err := BakaWebQuery(txApp, user, TIMETABLE_PUBLIC)
      if err != nil { return err }
      if status != 200 { return fmt.Errorf("invalid status code: %v %v", status, resp) }

      srcs, err := ParseSourcesWeb(resp)
      if err != nil { return err }

      txApp.Logger().Info(fmt.Sprint(srcs.AsMap()))

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

    fmt.Print(err)

    if err != nil { app.Logger().Error(err.Error(), err) }
  }
}
