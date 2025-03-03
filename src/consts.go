package src

import "time"

const (
  BAKATOKEN = "bakatoken"
  BAKAREFRESTOKEN = "bakarefreshtoken"
  BAKAVALID = "bakavalid"
  BAKATOKEN_EXPIRES = "bakatoken_expires"
  USERS = "users"
  BAKA_PATH = "https://bakalari.gchd.cz/"
  DATA = "data"
  OWNER = "owner"
  NAME = "name"
  BAKACOOKIE = "bakacookie"
  BAKACOOKIE_EXPIRES = "bakacookie_expires"
  LAST_USED = "last_used"
  LAST_UPDATED = "last_updated"
  TYPE = "type"
  SOURCES = "sources"
  TIMETABLE_PUBLIC = "TimeTable/Public"
  DESC = "desc"

  TEACHER = "Teacher"
  CLASS = "Class"
  ROOM = "Room"
  EVENTS = "Events"

  ACTUAL = "Actual"
  NEXT = "Next"
  PERMANENT = "Permanent"

  PRIVATE = "Private"
  VAPID = "vapid"
  WANTS_REFRESH = "wants_refresh"

  GET = "GET"
  POST = "POST"

  MARKS = "marks"
  TIMETABLE_ACTUAL = "timetable/actual"
  TIMETABLE_PERMANENT = "timetable/permanent"
)

var (
  VAPID_PRIVKEY = ""
  VAPID_PUBKEY = "BGl8lG0dFZxVzpEwgnPQlHaqDuaBojbFJHJzh2CMYi8mZshivG7RRkGDLKAC6E23E6ELtp3ikBXuepRJBMRlbwc"
)

func GetTTime() string {
  // TODO: remove
  year, week := time.Now().ISOWeek()
  if year == 2025 && week == 7 { return NEXT }

  
  wd := time.Now().Weekday()
  if wd == time.Sunday || wd == time.Sunday { return NEXT }
  return ACTUAL
}

func GetCDate() string {
  nw := time.Now()
  return nw.Add(time.Hour * -24 * time.Duration(nw.Weekday() - 1)).Format("2006-01-02")
}
