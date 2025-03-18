package src

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

var (
  czechDayNames = []string{"po", "út", "st", "čt", "pá"}
  czechIsoDayNames = []string{"ne", "po", "út", "st", "čt", "pá", "so"}
)

type BakaIdExpand struct {
  Id string
  Abbrev string
  Name string
}

type BakaIdExpandGroup struct {
  ClassId string
  Id string
  Abbrev string
  Name string
}

type BakaEvents struct {
  Events []BakaEvent
}

type BakaEvent struct {
  Id string
  Title string
  Description string
  Times []BakaEventTime
  EventType BakaIdExpand
  Classes []BakaIdExpand
  ClassSets []BakaIdExpand
  Teachers []BakaIdExpand
  TeacherSets []BakaIdExpand
  Rooms []BakaIdExpand
  RoomSets []BakaIdExpand
  Students []BakaIdExpand
  Note string
  DateChanged time.Time
}

type BakaEventTime struct {
  WholeDay bool
  StartTime time.Time
  EndTime time.Time
  IntervalStartTime *time.Time
  IntervalEndTime *time.Time
  IntervalDate string
}

type BakaAbsence struct {
  PercentageTreshold float64
  Absences []BakaAbsenceUnit
  AbsencesPerSubject []BakaAbsenceSubjectUnit
}

type BakaAbsenceUnit struct {
  Date time.Time
  Unsolved int
  Ok int
  Missed int
  Late int
  Soon int
  School int
  DistanceTeaching int
}

type BakaAbsenceSubjectUnit struct {
  SubjectName string
  LessonsCount int
  Base int
  Late int
  Soon int
  School int
  DistanceTeaching int
}

type BakaMarks struct {
  Subjects []BakaMarksSubject
  MarkOptions []BakaIdExpand
}

type BakaMarksSubject struct {
  Marks []BakaMark
  Subject BakaIdExpand
  AverageText string
  TemporaryMark string
  SubjectNote string
  TemporaryMarkNote string
  PointsOnly bool
  MarkPredictionEnabled bool
}

type BakaMark struct {
  MarkDate time.Time
  EditDate time.Time
  Caption string
  Theme string
  MarkText string
  IsInvalidDate bool
  TeacherId string
  Type string
  TypeNote string
  Weight int
  SubjectId string
  IsNew bool
  IsPoints bool
  CalculatedMarkText string
  ClassRankText string
  Id string
  PointsText string
  MaxPoints int
  ConfirmedWhen time.Time
  ConfirmedBy string
  MarkConfirmationState string
}

type BakaTimeTable struct {
  Hours []BakaTimeTableHour
  Days []BakaTimeTableDay
  Classes []BakaIdExpand
  Groups []BakaIdExpandGroup
  Subjects []BakaIdExpand
  Teachers []BakaIdExpand
  Rooms []BakaIdExpand
  Cycles []BakaIdExpand
  Students []BakaIdExpand
}

type BakaTimeTableHour struct {
  Id int
  Caption string
  BeginTime string
  EndTime string
}

type BakaTimeTableDay struct {
  Atoms []BakaTimeTableAtom
  AssistanceAtoms []any
  DayOfWeek int
  Date time.Time
  DayDescription string
  DayType string
}

type BakaTimeTableAtom struct {
  HourId int
  GroupIds []string
  SubjectId string
  TeacherId string
  RoomId string
  IsLastRoomLesson bool
  CycleIds []string
  Change *BakaTimeTableChange
  HomeworkIds []string
  Homeworks []any
  Theme string
  Assistants []any
  Notice string
  LessonRelease string
}

type BakaTimeTableChange struct {
  ChangeSubject any
  Day time.Time
  Hours string
  ChangeType string
  Description string
  Time string
  TypeAbbrev string
  TypeName string
  AtomType string
}

type MarkNotif struct {
  Title string
  Text string
  Num string
}

type Notif interface {
  JSONEncode() string
}

func (n MarkNotif) JSONEncode() string {
  imageStr := map[string]string{
    "1": "https://finakademie.cz/wp-content/uploads/2016/09/jednicka.jpg",
    "2": "http://www2.rozhlas.cz/podcast/img/podcast-dvojka.jpg",
    "3": "https://www.good-drinks.ch/media/image/2f/0f/28/295.png",
    "4": "https://prod-ripcut-delivery.disney-plus.net/v1/variant/disney/1F0D2A03B3EC79C06674F6CEF4D264A8EA2A31A8712385069789066FF594D9A2/scale?width=1200&aspectRatio=1.78&format=jpeg",
    "5": "https://onehotbook.cz/cdn/shop/products/spravna_petka-4_web_grande.jpg?v=1607959695",
  }[n.Num]
  if imageStr != "" {
    imageStr = `,"image":"` + imageStr + `"`
  }
  return `{"type":"notif","title":"` + n.Title + `","options":{"body":"` + n.Text + `" ` + imageStr + `}}`
}

func (m BakaMark) Notif(subj BakaMarksSubject) Notif {
  return MarkNotif{
    Title: "nová známka: " + subj.Subject.Abbrev,
    Text: m.Caption + ": " + m.MarkText + " (" + m.TypeNote + ")\n" + subj.Subject.Name + ": " + subj.AverageText,
    Num: strings.TrimSuffix(m.MarkText, "-"),
  }
}

func CompareBakaMarks(oldm, newm BakaMarks) []Notif {
  res := make([]Notif, 0)
  for _, subj := range newm.Subjects {
    idx := slices.IndexFunc(oldm.Subjects, func(s BakaMarksSubject) bool {
      return strings.TrimSpace(s.Subject.Id) == strings.TrimSpace(subj.Subject.Id)
    })
    if idx == -1 {
      for _, mark := range subj.Marks {
        res = append(res, mark.Notif(subj))
      }
      continue
    }
    oldsubj := oldm.Subjects[idx]
    for _, mark := range subj.Marks {
      if slices.ContainsFunc(oldsubj.Marks, func(m BakaMark) bool { return m.Id == mark.Id }) { continue }
      res = append(res, mark.Notif(subj))
    }
  }
  return res
}

type TimeTableNotif struct {
  Title string
  Text string
}

func (n TimeTableNotif) JSONEncode() string {
  return `{"type":"notif","title":"` + n.Title + `","options":{"body":"` + n.Text + `"}}`
}

func FindBakaExpand(list []BakaIdExpand, id string) BakaIdExpand {
  idx := slices.IndexFunc(list, func(e BakaIdExpand) bool { return e.Id == id })
  if idx == -1 { return BakaIdExpand{ "?", "?", "?" } }
  return list[idx]
}

func FindBakaExpandGroup(list []BakaIdExpandGroup, id string) BakaIdExpandGroup {
  idx := slices.IndexFunc(list, func(e BakaIdExpandGroup) bool { return e.Id == id })
  if idx == -1 { return BakaIdExpandGroup{ "?", "?", "?", "?" } }
  return list[idx]
}

func CompareBakaTimeTables(oldt, newt BakaTimeTable) []Notif {
  res := make([]Notif, 0)
  hours := make(map[int]BakaTimeTableHour)
  for _, hour := range oldt.Hours {
    hours[hour.Id] = hour
  }
  for _, hour := range newt.Hours {
    hours[hour.Id] = hour
  }
  for i := range 5 {
    olddi := slices.IndexFunc(oldt.Days, func(d BakaTimeTableDay) bool { return d.DayOfWeek == i })
    newdi := slices.IndexFunc(newt.Days, func(d BakaTimeTableDay) bool { return d.DayOfWeek == i })
    if olddi == -1 && newdi == -1 { continue }
    if olddi == -1 {
      newd := newt.Days[newdi]
      res = append(res, TimeTableNotif{
        Title: "změna rozvrhu " + czechDayNames[newd.DayOfWeek],
        Text: "den přidán",
      })
      continue
    }
    if newdi == -1 {
      oldd := oldt.Days[olddi]
      res = append(res, TimeTableNotif{
        Title: "změna rozvrhu " + czechDayNames[oldd.DayOfWeek],
        Text: "den odebrán",
      })
      continue
    }
    oldd := oldt.Days[olddi]
    newd := newt.Days[newdi]
    daynotifs := make([]string, 0)
    for _, hour := range hours {
      oldhi := slices.IndexFunc(oldd.Atoms, func(a BakaTimeTableAtom) bool { return a.HourId == hour.Id })
      newhi := slices.IndexFunc(newd.Atoms, func(a BakaTimeTableAtom) bool { return a.HourId == hour.Id })
      if oldhi == -1 && newhi == -1 { continue }
      if oldhi == -1 {
        newh := newd.Atoms[newhi]
        groups := make([]string, len(newh.GroupIds))
        for i, gid := range newh.GroupIds {
          groups[i] = FindBakaExpandGroup(newt.Groups, gid).Abbrev
        }
        daynotifs = append(daynotifs, fmt.Sprintf(
          "%d. hodina přidána: %s s %s v %s g %s",
          newh.HourId,
          FindBakaExpand(newt.Subjects, newh.SubjectId).Abbrev,
          FindBakaExpand(newt.Teachers, newh.TeacherId).Abbrev,
          FindBakaExpand(newt.Rooms, newh.RoomId).Abbrev,
          strings.Join(groups, ", "),
        ))
        continue
      }
      if newhi == -1 {
        oldh := oldd.Atoms[oldhi]
        daynotifs = append(daynotifs, fmt.Sprintf(
          "%d. hodina odebrána (%s)",
          oldh.HourId,
          FindBakaExpand(oldt.Subjects, oldh.SubjectId),
        ))
        continue
      }
      oldh := oldd.Atoms[oldhi]
      newh := newd.Atoms[newhi]
      nstr := FindBakaExpand(newt.Subjects, newh.SubjectId).Abbrev
      ostr := nstr
      if oldh.SubjectId != newh.SubjectId {
        nstr += " (" + FindBakaExpand(oldt.Subjects, oldh.SubjectId).Abbrev + ")"
      }
      if oldh.TeacherId != newh.TeacherId {
        nstr += "s " + FindBakaExpand(newt.Teachers, newh.TeacherId).Abbrev +
                " (" + FindBakaExpand(oldt.Teachers, oldh.TeacherId).Abbrev + ")"
      }
      if oldh.RoomId != newh.RoomId {
        nstr += "v " + FindBakaExpand(newt.Teachers, newh.TeacherId).Abbrev +
                " (" + FindBakaExpand(oldt.Teachers, oldh.TeacherId).Abbrev + ")"
      }
      if !slices.Equal(oldh.GroupIds, newh.GroupIds) {
        oldgs := make([]string, len(oldh.GroupIds))
        for i, gid := range oldh.GroupIds {
          oldgs[i] = FindBakaExpandGroup(newt.Groups, gid).Abbrev
        }
        newgs := make([]string, len(newh.GroupIds))
        for i, gid := range newh.GroupIds {
          newgs[i] = FindBakaExpandGroup(newt.Groups, gid).Abbrev
        }
        nstr += "g " + strings.Join(newgs, ", ") + " (" + strings.Join(oldgs, ", ") + ")"
      }
      if nstr == ostr { continue }
      daynotifs = append(daynotifs, fmt.Sprintf(
        "%d. hodina změna: %s",
        newh.HourId, nstr,
      ))
    }
    if len(daynotifs) == 0 { continue }
    res = append(res, TimeTableNotif{
      Title: "změna rozvrhu " + czechDayNames[oldd.DayOfWeek],
      Text: strings.Join(daynotifs, "\n"),
    })
  }
  return res
}

type AbsenceNotif struct {
  Title string
  Text string
}

func (n AbsenceNotif) JSONEncode() string {
  return `{"type":"notif","title":"` + n.Title + `","options":{"body":"` + n.Text + `"}}`
}

func FindPendingAbsences(abs BakaAbsence) []Notif {
  res := make([]Notif, 0)
  for _, au := range abs.Absences {
    if au.Unsolved > 0 {
      res = append(res, AbsenceNotif{
        Title: "neomluvená absence",
        Text: fmt.Sprintf(
          "%s %d. %d. %d: %d h",
          czechIsoDayNames[au.Date.Weekday()],
          au.Date.Day(),
          au.Date.Month(),
          au.Date.Year(),
          au.Unsolved,
        ),
      })
      continue
    }
    if au.Late > 0 && time.Since(au.Date) < ABSENCE_STATUE_OF_REPOSE {
      res = append(res, AbsenceNotif{
        Title: "neomluvený pozdní příchod",
        Text: fmt.Sprintf(
          "%s %d. %d. %d: %d h",
          czechIsoDayNames[au.Date.Weekday()],
          au.Date.Day(),
          au.Date.Month(),
          au.Date.Year(),
          au.Unsolved,
        ),
      })
      continue
    }
  }
  return res
}

type BakaInvalidNotif struct {}

func (n BakaInvalidNotif) JSONEncode() string {
  return `{"type":"notif","title":"vypršela cookieska 🍪","options":{"body":"přihlaš se prosím znovu nebo ti vytrhám všechny zuby díky 🦷"}}`
}

func (e BakaEvent) ContainsDay(day time.Time) bool {
  for _, t := range e.Times {
    std, stm, sty := t.StartTime.Date()
    diddy, dm, dy := day.Date()
    if diddy == std && stm == dm && sty == dy { return true }
  }
  return false
}

func BakaIdExpandListContainsId(lst []BakaIdExpand, id string) bool {
  for _, e := range lst {
    if e.Id == id { return true }
  }
  return false
}
