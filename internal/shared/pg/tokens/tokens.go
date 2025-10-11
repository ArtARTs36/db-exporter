package tokens

const (
	CreateTable = iota + 1
	AlterTable  = iota + 2

	BracketLeft  // (
	BracketRight // )

	NotNull

	Comma

	PrimaryKey

	Semicolon
)
