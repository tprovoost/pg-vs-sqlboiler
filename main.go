package main

import (
	"time"

	_ "github.com/lib/pq"

	"github.com/tprovoost/pg-vs-sqlboiler/modules"
	"github.com/tprovoost/pg-vs-sqlboiler/modules/shared"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func main() {
	boil.DebugMode = false

	var suite shared.BenchmarkSuite

	// First clean up the database.
	modules.BoilerRunBenchmark(modules.BoilerCleanUp, 1)

	suite.ReadOne = []shared.Benchmark{
		modules.BoilerRunBenchmark(modules.BoilerReadOne, 4000),
		modules.PGRunBenchmark(modules.PGReadOne, 4000),
	}

	suite.Insert = []shared.Benchmark{
		modules.BoilerRunBenchmark(modules.BoilerInsert, 1000),
		modules.PGRunBenchmark(modules.PGInsert, 1000),
	}

	suite.FetchIn = []shared.Benchmark{
		modules.BoilerRunBenchmark(modules.BoilerFetchIn, 1000),
		modules.PGRunBenchmark(modules.PGFetchIn, 1000),
	}

	suite.Print()
}

func averageDuration(d time.Duration, N int) int64 {
	return d.Microseconds() / int64(N)
}
