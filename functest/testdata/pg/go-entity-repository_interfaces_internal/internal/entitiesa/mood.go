package entitiesa

type Mood string

const (
    MoodUndefined Mood = ""
    MoodOk Mood = "ok"
    MoodHappy Mood = "happy"
)

func (e Mood) Valid() bool {
    switch e {
    case MoodUndefined:
        return true
    case MoodOk:
        return true
    case MoodHappy:
        return true
    default:
        return false
    }
}
