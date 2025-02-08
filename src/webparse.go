package src

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/go-shiori/dom"
	"golang.org/x/net/html"
)

type TimeTableCell struct {
  Subject string `json:"subject"`
  Teacher string `json:"teacher"`
  Room string `json:"room"`
  Group string `json:"group"`
  Detail string `json:"detail"`
  Color string `json:"color"`
}

type TimeTableHour struct {
  Cells []TimeTableCell `json:"cells"`
}

type TimeTableDay struct {
  Title string `json:"title"`
  Special string `json:"special"`
  Hours []TimeTableHour `json:"hours"`
}

type TimeTableHourTitle struct {
  Idx int `json:"idx"`
  Dur string `json:"dur"`
}

type TimeTable struct {
  Hours []TimeTableHourTitle `json:"hours"`
  Days []TimeTableDay `json:"days"`
}

func NewlineInnerText(el *html.Node) string {
  return strings.ReplaceAll(dom.InnerText(el), "<br>", "\n")
}

func ParseTimeTableWeb(htmldoc string) (tt TimeTable, err error) {
  defer func() {
    if r := recover(); r != nil { err = fmt.Errorf("%#v", r) }
  }()
  doc, err := html.Parse(strings.NewReader(htmldoc))
  if err != nil { return TimeTable{}, err }

  bkmain := dom.QuerySelector(doc, ".bk-timetable-main")
  bkrows := dom.QuerySelectorAll(bkmain, ".bk-timetable-row")

  timetable := TimeTable{
    Hours: make([]TimeTableHourTitle, 0),
    Days: make([]TimeTableDay, 0),
  }

  hourtitles := dom.QuerySelectorAll(bkmain, ".bk-hour-wrapper")
  for _, htitle := range hourtitles {
    hourtres := TimeTableHourTitle{}

    numel := dom.QuerySelector(htitle, ".num")
    timespans := dom.QuerySelectorAll(htitle, "span")

    if len(timespans) == 3 {
      hourtres.Dur = dom.InnerHTML(timespans[0])+ 
        " - " + dom.InnerHTML(timespans[2])
    }

    idx, _ := strconv.Atoi(dom.InnerHTML(numel))
    hourtres.Idx = idx
    timetable.Hours = append(timetable.Hours, hourtres)
  }

  for _, row := range bkrows {
    cells := dom.QuerySelectorAll(row, ".bk-timetable-cell")
    rowres := TimeTableDay{
      Hours: make([]TimeTableHour, 0),
    }

    dayel := dom.QuerySelector(row, ".bk-day-day")
    dateel := dom.QuerySelector(row, ".bk-day-date")

    rowres.Title = dom.InnerHTML(dayel) + " " + dom.InnerHTML(dateel)

    if len(cells) == 1 {
      titleel := dom.QuerySelector(cells[0], "span")
      rowres.Special = dom.InnerHTML(titleel)
      goto appendday
    }

    for _, cell := range cells {
      hours := dom.QuerySelectorAll(cell, ".day-item-hover")
      cellres := TimeTableHour{
        Cells: make([]TimeTableCell, 0),
      }
      for _, hour := range hours {
        hourres := TimeTableCell{}

        roomel := dom.QuerySelector(hour, ".right > div")
        if roomel != nil { hourres.Room = NewlineInnerText(roomel) }
        groupel := dom.QuerySelector(hour, ".left > div")
        if groupel != nil { hourres.Group = NewlineInnerText(groupel) }
        subjectel := dom.QuerySelector(hour, ".middle")
        if subjectel != nil { hourres.Subject = NewlineInnerText(subjectel) }
        teacherel := dom.QuerySelector(hour, ".bottom > span")
        if teacherel != nil { hourres.Teacher = NewlineInnerText(teacherel) }
        classes := strings.Split(dom.ClassName(hour), " ")
        hourres.Color = "white"
        if slices.Contains(classes, "pink") {
          hourres.Color = "pink"
        } else if slices.Contains(classes, "green") {
          hourres.Color = "green"
        }

        cellres.Cells = append(cellres.Cells, hourres)
      }
      rowres.Hours = append(rowres.Hours, cellres)
    }
    appendday: timetable.Days = append(timetable.Days, rowres)
  }

  return timetable, nil
}

type WebSourcePair struct {
  Id string `json:"id"`
  Name string `json:"name"`
}

type WebSources struct {
  Teachers []WebSourcePair `json:"teachers"`
  Rooms []WebSourcePair `json:"rooms"`
  Classes []WebSourcePair `json:"classes"`
}

func ParseSourcesWeb(htmldoc string) (WebSources, error) {
  doc, err := html.Parse(strings.NewReader(htmldoc))
  if err != nil { return WebSources{}, err }
  sources := WebSources{
    Teachers: make([]WebSourcePair, 0),
    Rooms: make([]WebSourcePair, 0),
    Classes: make([]WebSourcePair, 0),
  }

  teachers := dom.QuerySelector(doc, "#selectedTeacher")
  rooms := dom.QuerySelector(doc, "#selectedRoom")
  classes := dom.QuerySelector(doc, "#selectedClass")

  teacheropts := dom.QuerySelectorAll(teachers, "option")
  for _, opt := range teacheropts {
    if dom.InnerHTML(opt) == "" { continue }
    sources.Teachers = append(sources.Teachers, WebSourcePair{
      Id: dom.GetAttribute(opt, "value"),
      Name: dom.InnerHTML(opt),
    })
  }

  roomopts := dom.QuerySelectorAll(rooms, "option")
  for _, opt := range roomopts {
    if dom.InnerHTML(opt) == "" { continue }
    sources.Rooms = append(sources.Rooms, WebSourcePair{
      Id: dom.GetAttribute(opt, "value"),
      Name: dom.InnerHTML(opt),
    })
  }
  
  classopts := dom.QuerySelectorAll(classes, "option")
  for _, opt := range classopts {
    if dom.InnerHTML(opt) == "" { continue }
    sources.Classes = append(sources.Classes, WebSourcePair{
      Id: dom.GetAttribute(opt, "value"),
      Name: dom.InnerHTML(opt),
    })
  }

  return sources, nil
}
