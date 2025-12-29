package pipeline

import "github.com/example/block-indexer/pkg/models"

type DagBlockPayload struct {
	Block   models.DagBlock
	Parents []string
}
