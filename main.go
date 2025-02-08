package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/anteat3r/skrat-org/src"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

var lastReloaded = time.Now()

func main() {
    app := pocketbase.New()


    app.OnServe().BindFunc(func(se *core.ServeEvent) error {
        datacoll, _ := app.FindCollectionByNameOrId(src.DATA)

        se.Router.GET("/{path...}", apis.Static(os.DirFS("/root/skrat-org/web/dist"), false))
        
        appsv, err := os.ReadFile("/root/skrat-org/web/src/App.svelte")
        if err != nil { return err }
        fmt.Println(string(appsv))
        re, err := regexp.Compile(`<Route path="(.+)"`)
        if err != nil { return err}
        ms := re.FindAllSubmatch(appsv, -1)
        if ms == nil { ms = [][][]byte{} }
        fmt.Println(ms)
        for _, m := range ms {
          if m == nil { continue }
          if len(m) < 2 { continue }
          route := string(m[1])
          if route == "/" { continue }
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

        se.Router.GET("/api/system/reload", func(e *core.RequestEvent) error {
          nw := time.Now()
          if nw.Sub(lastReloaded) > time.Minute {
            lastReloaded = nw
            err := exec.Command("sh", "/root/skrat-org/reload.sh").Run()
            if err != nil { return e.Error(400, err.Error(), nil) }
            return e.String(200, "")
          }
          return e.Error(401, "reloaded a minute ago, slow down", nil)
        })

        se.Router.POST(
          "/api/kleo/login",
          src.LoginHandler(app),
        ).Bind(apis.RequireAuth(src.USERS))

        se.Router.GET(
          "/api/kleo/endp",
          src.EndpHandler(app, datacoll),
        ).Bind(apis.RequireAuth(src.USERS))

        se.Router.GET(
          "/api/kleo/web/{time}/{ttype}/{name}",
          src.WebTimeTableHandler(app, datacoll),
        ).Bind(apis.RequireAuth(src.USERS))

        se.Router.GET(
          "/api/kleo/websrcs",
          src.WebSourcesHandler(app),
        ).Bind(apis.RequireAuth(src.USERS))

        return se.Next()
    })

    if err := app.Start(); err != nil { log.Fatal(err) }
}
