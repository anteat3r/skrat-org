package src

import "time"

type BakaIdExpand struct {
  Id string
  Abbrev string
  Name string
}

type BakaEvents []BakaEvent

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
  Note *string
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
