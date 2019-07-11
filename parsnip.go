package parsnip

// Registry stores "immutable" (after initialization) job information:
// - all queue names
// - queue weights, used by workers; when a worker monitors several queues,
//   it should pick up a job to work on probabilistically using the weights
// - which jobs map to which queues (one job can only be sent on a single queue)
// - job options (retry and backoff policy)
// - job resource requests (named resources each job needs, workers will need to
//   use middleware that provides this, if they want to listen on that queue)

const defaultQueue = "default"
