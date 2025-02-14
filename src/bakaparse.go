package src

import "time"

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
