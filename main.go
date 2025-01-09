package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
    app := pocketbase.New()

    app.OnServe().BindFunc(func(se *core.ServeEvent) error {

        se.Router.GET("/", apis.Static(os.DirFS("/root/skrat-org/web"), false))
        se.Router.GET("/ss", func(e *core.RequestEvent) error {
          return app.RunInTransaction(func(txApp core.App) error {
            key := e.Request.URL.Query().Get("k")
            value := e.Request.URL.Query().Get("v")
            if key == "" || value == "" { return fmt.Errorf("invalid params") }
            _, err := txApp.FindFirstRecordByData("urls", "name", key)
            if err == nil { return fmt.Errorf("already defined") }
            coll, _ := txApp.FindCollectionByNameOrId("urls")
            nrec := core.NewRecord(coll)
            nrec.Set("name", key)
            nrec.Set("value", value)
            return txApp.Save(nrec)
          })
        })
        se.Router.GET("/s/:key", func(e *core.RequestEvent) error {
          key := e.Request.PathValue("key")
          if key == "" { return fmt.Errorf("invalid param") }
          rec, err := app.FindFirstRecordByData("urls", "name", key)
          if err != nil { return err }
          return e.Redirect(301, rec.GetString("value"))
        })

        return se.Next()
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
