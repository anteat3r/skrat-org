package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/anteat3r/skrat-org/src"
	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/spf13/cobra"
)

var (
  lastReloaded = time.Now()
)

func main() {
  envmap, err := godotenv.Read()
  if err != nil { panic(err) }

  vapidPriv, ok := envmap["VAPID_PRIVKEY"]
  if !ok { panic("vapid privkey not found") }

  src.VAPID_PRIVKEY = vapidPriv

  app := pocketbase.New()

    app.OnServe().BindFunc(func(se *core.ServeEvent) error {
        datacoll, _ := app.FindCollectionByNameOrId(src.DATA)

        se.Router.GET("/{path...}", apis.Static(os.DirFS("/root/skrat-org/web/dist"), false))
        
        appsv, err := os.ReadFile("/root/skrat-org/web/src/App.svelte")
        if err != nil { return err }
        re, err := regexp.Compile(`<Route path="(.+)"`)
        if err != nil { return err}
        ms := re.FindAllSubmatch(appsv, -1)
        if ms == nil { ms = [][][]byte{} }
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
          err = e.Redirect(301, rec.GetString("value"))
          if err != nil { return err }
          rec.Set("count", rec.GetInt("count") + 1)
          err = app.Save(rec)
          if err != nil { return err }
          return nil
        })
        se.Router.GET("/s", func(e *core.RequestEvent) error {
          if len(e.Request.URL.Query()) != 1 { return e.Error(400, "key not specified", nil) }
          var key string
          for k, _ := range e.Request.URL.Query() { key = k }
          if key == "" { return e.Error(400, "invalid key", nil ) }
          rec, err := app.FindFirstRecordByData("urls", "name", key)
          if err != nil { return err }
          err = e.Redirect(301, rec.GetString("value"))
          if err != nil { return err }
          rec.Set("count", rec.GetInt("count") + 1)
          err = app.Save(rec)
          if err != nil { return err }
          return nil
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
        ).Bind(apis.RequireAuth(src.USERS)).Bind(src.RequireBakaValid)

        se.Router.GET(
          "/api/kleo/websrcs",
          src.WebSourcesHandler(app),
        ).Bind(apis.RequireAuth(src.USERS)).Bind(src.RequireBakaValid)

        se.Router.GET(
          "/api/kleo/daytt/{ttype}",
          src.DayOverviewHandler(app, datacoll),
        ).Bind(apis.RequireAuth(src.USERS)).Bind(src.RequireBakaValid)

        se.Router.GET(
          "/api/kleo/marks",
          src.MarksHandler(app, datacoll),
        ).Bind(apis.RequireAuth(src.USERS)).Bind(src.RequireBakaValid)

        se.Router.GET(
          "/api/kleo/events",
          src.EventsHandler(app, datacoll),
        ).Bind(apis.RequireAuth(src.USERS)).Bind(src.RequireBakaValid)

        se.Router.POST(
          "/api/kleo/setupnotifs",
          src.StoreVapidEndpoint(app),
        ).Bind(apis.RequireAuth(src.USERS))

        se.Router.POST(
          "/api/kleo/vapidtest",
          src.VapidTestHandler(app),
        ).Bind(apis.RequireAuth(src.USERS))

        se.Router.GET(
          "/api/kleo/uncache",
          func(e *core.RequestEvent) error {
            name := e.Request.URL.Query().Get(src.NAME)
            ttype := e.Request.URL.Query().Get(src.TYPE)
            owner := e.Request.URL.Query().Get(src.OWNER)
            src.UnCacheData(app, name, ttype, owner)
            return e.String(200, "")
          },
        )

        se.Router.GET(
          "/api/kleo/mytt",
          src.MyTimeTableHandler(app, datacoll),
        ).Bind(apis.RequireAuth(src.USERS)).Bind(src.RequireBakaValid)

				se.Router.POST(
					"/api/reloadsrcs",
					func(e *core.RequestEvent) error {
						src.TimeTableSourcesReload(app)
						return e.String(200, "")
					},
				).Bind(apis.RequireAuth(src.USERS))


        app.Cron().MustAdd(
          "ttreload",
          "* 5-17 * * 0-6",
          src.TimeTableReload(app, datacoll),
        )

        app.Cron().MustAdd(
          "srcsreload",
          "1 7 * * 6",
          src.TimeTableSourcesReload(app),
        )

        app.Cron().MustAdd(
          "personalreload",
          "* 6-16 * * 0-5",
          src.PersonalReload(app, datacoll),
        )

        app.Cron().MustAdd(
          "eveningrefresh",
          "0 17 * * 0-5",
          src.EveningRefresh(app, datacoll),
        )

        return se.Next()
    })

    app.RootCmd.AddCommand(&cobra.Command{
      Use: "reloadsrcs",
      Run: func(cmd *cobra.Command, args []string) {
        src.TimeTableSourcesReload(app)()
      },
    })

    app.RootCmd.AddCommand(&cobra.Command{
      Use: "cleardata",
      Run: func(cmd *cobra.Command, args []string) {
        _, err := app.DB().NewQuery("delete from data").Execute()
        if err != nil { app.Logger().Error(err.Error(), err) }
      },
    })


    if err := app.Start(); err != nil { log.Fatal(err) }
}
