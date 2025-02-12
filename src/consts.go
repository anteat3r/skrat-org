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
)

func GetTTime() string {
  // TODO: remove
  year, week := time.Now().ISOWeek()
  if year == 2025 && week == 7 { return NEXT }

  
  wd := time.Now().Weekday()
  if wd == time.Sunday || wd == time.Sunday { return NEXT }
  return ACTUAL
}
