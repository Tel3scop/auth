package metrics

import (
	helperMetrics "github.com/Tel3scop/helpers/metrics"
)

// Init инициализация базовых метрик
func Init(namespace, appName, subsystem string, bucketsStart, bucketsFactor float64, bucketsCount int) error {
	err := helperMetrics.Init(namespace, appName, subsystem, bucketsStart, bucketsFactor, bucketsCount)
	if err != nil {
		return err
	}

	return nil
}
