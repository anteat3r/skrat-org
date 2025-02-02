package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
    app := pocketbase.New()

    app.OnServe().BindFunc(func(se *core.ServeEvent) error {

        se.Router.GET("/{path...}", apis.Static(os.DirFS("/root/skrat-org/web/dist"), false))
        
        appsv, err := os.ReadFile("/root/skrat-org/web/src/App.svelte")
        if err != nil { return err }
        fmt.Println(string(appsv))
        re, _ := regexp.Compile(`<Route path="(.+)"`)
        ms := re.FindAllSubmatch(appsv, 0)
        if ms == nil { ms = [][][]byte{} }
        fmt.Println(ms)
        for _, m := range ms {
          if m == nil { continue }
          if len(m) < 2 { continue }
          route := string(m[1])
          se.Router.GET(route, func(e *core.RequestEvent) error {
            return e.FileFS(os.DirFS("/root/skrat-org/web/dist"), "index.html")
          })
        }
        
        se.Router.GET("/ss", func(e *core.RequestEvent) error {
          return app.RunInTransaction(func(txApp core.App) error {
            key := e.Request.URL.Query().Get("k")
            value := e.Request.URL.Query().Get("v")
            if key == "" || value == "" { return e.Error(400, "invalid params", nil) }
            if !strings.HasPrefix(value, "https://") {
              value = "https://" + value
            }
            _, err := txApp.FindFirstRecordByData("urls", "name", key)
            if err == nil { return e.Error(400, "already defined", nil ) }
            coll, _ := txApp.FindCollectionByNameOrId("urls")
            nrec := core.NewRecord(coll)
            nrec.Set("name", key)
            nrec.Set("value", value)
            err = txApp.Save(nrec)
            if err != nil { return err }
            return e.String(200, "")
          })
        })
        se.Router.GET("/s/{key}", func(e *core.RequestEvent) error {
          key := e.Request.PathValue("key")
          if key == "" { return e.Error(400, "invalid key", nil ) }
          rec, err := app.FindFirstRecordByData("urls", "name", key)
          if err != nil { return err }
          return e.Redirect(301, rec.GetString("value"))
        })
        se.Router.GET("/s/{key}/", func(e *core.RequestEvent) error {
          return e.Redirect(301, "https://skrat.org/s/" + e.Request.PathValue("key"))
        })
        se.Router.GET("/f/{key}", func(e *core.RequestEvent) error {
          key := e.Request.PathValue("key")
          if key == "" { return e.Error(400, "invalid key", nil ) }
          rec, err := app.FindFirstRecordByData("files", "name", key)
          if err != nil { return err }
          return e.Redirect(301, "https://skrat.org/api/files/files/" + rec.Id + "/" + rec.GetString("file"))
        })


        return se.Next()
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
