// Copyright 2016 Eleme Inc. All rights reserved.

/*

Package health implements the health statistic aggregation.

	// Total
	aggregationInterval           // Health aggregation interval, 300s (5min)
	numIndexTotal                 // Number of metric indexes total.
	numClients                    // Number of detector clients.
	// Aggregation
	detectionCost                 // Average of detection time cost in last interval.
	numMetricIncomed              // Number of metrics incomed in last interval.
	numMetricDetected             // Number of metrics detected in last interval.
	numAlertingEvents             // Number of alerting events in last interval.

You can visit /api/info for health info.

*/
package health
