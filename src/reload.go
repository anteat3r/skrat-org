package src

import (
	"encoding/json"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

func TimeTableReload(app *pocketbase.PocketBase, datacoll *core.Collection) func() {
  return func() {
    err := app.RunInTransaction(func(txApp core.App) error {
      srcs, err := txApp.FindRecordsByFilter(
        SOURCES,
        `id != ""`,
        LAST_UPDATED,
        1, 0, 
      )
      if err != nil { return err }

      if len(srcs) < 0 { return fmt.Errorf("no suitable source found") }
      src := srcs[0]

      users, err := txApp.FindRecordsByFilter(
        USERS,
        BAKAVALID + ` = true`,
        LAST_USED, 1, 0,
      )
      if err != nil { return err }

      if len(users) < 1 { return fmt.Errorf("no suitable user found") }
      user := users[0]

      status, resp, err := BakaTimeTableQuery(txApp, user, GetTTime(), src.GetString(TYPE), src.GetString(NAME))
      if status != 200 { return fmt.Errorf("invalid status code: %v %v", status, resp) }

      tt, err := ParseTimeTableWeb(resp)
      if err != nil { return err }

      resb, err := json.Marshal(tt)
      if err != nil { return err }

      var datarec *core.Record
      datarec, err = app.FindFirstRecordByFilter(
        DATA,
        OWNER + ` = "" && ` + NAME + ` = {:name}`,
        dbx.Params{"name": src.GetString(NAME)},
      )
      if err != nil {
        datarec = core.NewRecord(datacoll)
        datarec.Set(OWNER, "")
        datarec.Set(NAME, src.GetString(NAME))
        datarec.Set(TYPE, src.GetString(TYPE))
      }
      datarec.Set(DATA, string(resb))
      if datarec.GetString(TYPE) == "" {
        datarec.Set(TYPE, src.GetString(TYPE))
      }

      user.Set(LAST_USED, types.NowDateTime())

      err = txApp.Save(user)
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
      _, err := txApp.DB().NewQuery("delete * from " + SOURCES).Execute()
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

      user.Set(LAST_USED, types.NowDateTime())

      err = txApp.Save(user)
      if err != nil { return err }

      return nil
    })

    if err != nil { app.Logger().Error(err.Error(), err) }
  }
}
