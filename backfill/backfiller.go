package backfill

type BackFiller interface {
	BackFill(startingBlock int) error
}
