package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	countStaleConfigErrors = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "count_stale_config_errors_total",
		Help:      "The total number of times that threads hit stale config exception. Since a stale config exception triggers a refresh of the metadata, this number is roughly proportional to the number of metadata refreshes.",
	})
	countDonorMoveChunkStarted = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "count_donor_move_chunk_started_total",
		Help:      "The total number of times that the moveChunk command has started on the shard, of which this node is a member, as part of a chunk migration process. This increasing number does not consider whether the chunk migrations succeed or not.",
	})
	totalDonorChunkCloneTimeMillis = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "total_donor_chunk_clone_time_milliseconds",
		Help:      "The cumulative time, in milliseconds, taken by the clone phase of the chunk migrations from this shard, of which this node is a member. Specifically, for each migration from this shard, the tracked time starts with the moveChunk command and ends before the destination shard enters a catch-up phase to apply changes that occurred during the chunk migrations.",
	})
	totalCriticalSectionCommitTimeMillis = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "total_critical_section_commit_time_milliseconds",
		Help:      "The cumulative time, in milliseconds, taken by the update metadata phase of the chunk migrations from this shard, of which this node is a member. During the update metadata phase, all operations on the collection are blocked.",
	})
	totalCriticalSectionTimeMillis = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "total_critical_section_time_milliseconds",
		Help:      "The cumulative time, in milliseconds, taken by the catch-up phase and the update metadata phase of the chunk migrations from this shard, of which this node is a member.",
	})

	catalogCacheNumDatabaseEntries = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "catalog_cache_num_database_entries",
		Help:      "The total number of database entries that are currently in the catalog cache.",
	})
	catalogCacheNumCollectionEntries = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "catalog_cache_num_collection_entries",
		Help:      "The total number of collection entries (across all databases) that are currently in the catalog cache.",
	})
	catalogCacheContStaleConfigErrors = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "catalog_cache_count_stale_config_errors",
		Help:      "The total number of times that threads hit stale config exception. A stale config exception triggers a refresh of the metadata.",
	})
	catalogCacheTotalRefreshWaitTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "catalog_cache_total_refresh_wait_time_microseconds",
		Help:      "The cumulative time, in microseconds, that threads had to wait for a refresh of the metadata.",
	})
	catalogCacheNumActiveIncrementalRefreshes = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "catalog_cache_num_active_incremental_refreshes",
		Help:      "The number of incremental catalog cache refreshes that are currently waiting to complete.",
	})
	catalogCacheCountIncrementalRefreshesStarted = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "catalog_cache_count_incremental_refreshes_started",
		Help:      "    The cumulative number of incremental refreshes that have started.",
	})
	catalogCacheNumActiveFullRefreshes = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "catalog_cache_num_active_full_refreshes",
		Help:      "The number of full catalog cache refreshes that are currently waiting to complete.",
	})
	catalogCacheCountFullRefreshesStarted = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "catalog_cache_count_full_refreshes_started",
		Help:      "The cumulative number of full refreshes that have started.",
	})
	catalogCacheCountFailedRefreshes = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "sharding_statistics",
		Name:      "catalog_cache_count_failed_refreshes",
		Help:      "The cumulative number of full or incremental refreshes that have failed.",
	})
)

// Cursors are the cursor metrics
type ShardingStatistics struct {
	CountStaleConfigErrors               int64 `bson:"countStaleConfigErrors"`
	CountDonorMoveChunkStarted           int64 `bson:"countDonorMoveChunkStarted"`
	TotalDonorChunkCloneTimeMillis       int64 `bson:"totalDonorChunkCloneTimeMillis"`
	TotalCriticalSectionCommitTimeMillis int64 `bson:"totalCriticalSectionCommitTimeMillis"`
	TotalCriticalSectionTimeMillis       int64 `bson:"totalCriticalSectionTimeMillis"`
	CatalogCache                         `bson:"catalogCache"`
}

type CatalogCache struct {
	NumDatabaseEntries               int64 `bson:"numDatabaseEntries"`
	NumCollectionEntries             int64 `bson:"numCollectionEntries"`
	CountStaleConfigErrors           int64 `bson:"countStaleConfigErrors"`
	TotalRefreshWaitTimeMicros       int64 `bson:"totalRefreshWaitTimeMicros"`
	NumActiveIncrementalRefreshes    int64 `bson:"numActiveIncrementalRefreshes"`
	CountIncrementalRefreshesStarted int64 `bson:"countIncrementalRefreshesStarted"`
	NumActiveFullRefreshes           int64 `bson:"numActiveFullRefreshes"`
	CountFullRefreshesStarted        int64 `bson:"countFullRefreshesStarted"`
	CountFailedRefreshes             int64 `bson:"countFailedRefreshes"`
}

// Export exports the data to prometheus.
func (s *ShardingStatistics) Export(ch chan<- prometheus.Metric) {
	countStaleConfigErrors.Set(float64(s.CountStaleConfigErrors))
	countDonorMoveChunkStarted.Set(float64(s.CountDonorMoveChunkStarted))
	totalDonorChunkCloneTimeMillis.Set(float64(s.TotalDonorChunkCloneTimeMillis))
	totalCriticalSectionCommitTimeMillis.Set(float64(s.TotalCriticalSectionCommitTimeMillis))
	totalCriticalSectionTimeMillis.Set(float64(s.TotalCriticalSectionTimeMillis))
	catalogCacheNumDatabaseEntries.Set(float64(s.CatalogCache.NumDatabaseEntries))
	catalogCacheNumCollectionEntries.Set(float64(s.CatalogCache.NumCollectionEntries))
	catalogCacheContStaleConfigErrors.Set(float64(s.CatalogCache.CountStaleConfigErrors))
	catalogCacheTotalRefreshWaitTime.Set(float64(s.CatalogCache.TotalRefreshWaitTimeMicros))
	catalogCacheNumActiveIncrementalRefreshes.Set(float64(s.CatalogCache.NumActiveIncrementalRefreshes))
	catalogCacheCountIncrementalRefreshesStarted.Set(float64(s.CatalogCache.CountIncrementalRefreshesStarted))
	catalogCacheNumActiveFullRefreshes.Set(float64(s.CatalogCache.NumActiveFullRefreshes))
	catalogCacheCountFullRefreshesStarted.Set(float64(s.CatalogCache.CountFullRefreshesStarted))
	catalogCacheCountFailedRefreshes.Set(float64(s.CatalogCache.CountFailedRefreshes))

	countStaleConfigErrors.Collect(ch)
	countDonorMoveChunkStarted.Collect(ch)
	totalDonorChunkCloneTimeMillis.Collect(ch)
	totalCriticalSectionCommitTimeMillis.Collect(ch)
	totalCriticalSectionTimeMillis.Collect(ch)
	catalogCacheNumDatabaseEntries.Collect(ch)
	catalogCacheNumCollectionEntries.Collect(ch)
	catalogCacheContStaleConfigErrors.Collect(ch)
	catalogCacheTotalRefreshWaitTime.Collect(ch)
	catalogCacheNumActiveIncrementalRefreshes.Collect(ch)
	catalogCacheCountIncrementalRefreshesStarted.Collect(ch)
	catalogCacheNumActiveFullRefreshes.Collect(ch)
	catalogCacheCountFullRefreshesStarted.Collect(ch)
	catalogCacheCountFailedRefreshes.Collect(ch)
}

// Describe describes the metrics for prometheus
func (s *ShardingStatistics) Describe(ch chan<- *prometheus.Desc) {
	countStaleConfigErrors.Describe(ch)
	countDonorMoveChunkStarted.Describe(ch)
	totalDonorChunkCloneTimeMillis.Describe(ch)
	totalCriticalSectionCommitTimeMillis.Describe(ch)
	totalCriticalSectionTimeMillis.Describe(ch)
	catalogCacheNumDatabaseEntries.Describe(ch)
	catalogCacheNumCollectionEntries.Describe(ch)
	catalogCacheContStaleConfigErrors.Describe(ch)
	catalogCacheTotalRefreshWaitTime.Describe(ch)
	catalogCacheNumActiveIncrementalRefreshes.Describe(ch)
	catalogCacheCountIncrementalRefreshesStarted.Describe(ch)
	catalogCacheNumActiveFullRefreshes.Describe(ch)
	catalogCacheCountFullRefreshesStarted.Describe(ch)
	catalogCacheCountFailedRefreshes.Describe(ch)
}
