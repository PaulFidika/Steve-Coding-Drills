package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/creasty/defaults"
	"github.com/samber/mo"
)

// type JobRequest struct {
// 	ModelId string
// 	Prompt mo.Option[string] // `default:"mo.None[string]()"`
// 	NumImages int `default:"4"`
// 	AspectRatio string `default:"1/1"`
// }

type AspectRatioEnum string

const (
	AspectRatio1x1  AspectRatioEnum = "1/1"
	AspectRatio16x9 AspectRatioEnum = "16/9"
	AspectRatio9x16 AspectRatioEnum = "9/16"
)

func (a AspectRatioEnum) IsValid() bool {
	switch a {
	case AspectRatio1x1, AspectRatio16x9, AspectRatio9x16:
		return true
	}
	return false
}

func (a AspectRatioEnum) String() string {
	return string(a)
}


type JobRequest struct {
	ModelId string
	Prompt mo.Option[string]
	NumImages mo.Option[int] 
	AspectRatio mo.Option[string]
}

func (job *JobRequest) SetDefaults() {
	// if _, exists := job.Prompt.Get(); !exists {
	// 	job.Prompt = mo.None[string]()
	// }
	if job.NumImages.IsAbsent() {
		job.NumImages = mo.Some(4)
	}
	if job.AspectRatio.IsAbsent() {
		job.AspectRatio = mo.Some(AspectRatio1x1.String())
	}
}

func RunIt() {
	job := JobRequest{ModelId: "SDXL"}
	// defaults.MustSet(&job)
	job.SetDefaults()

	fmt.Printf("Job details: %+v\n", job)
	fmt.Printf("The number of images: %X\n", job.NumImages.MustGet())

	option1 := mo.Some(42)

	option2 := mo.None[int]()

	if option1.IsPresent() {
		fmt.Println("Option1 is", option1.MustGet())
	}

	fmt.Println("Option2 is", option2.OrElse(15))

	manualBenchmark()
}

func BenchmarkLoops(b *testing.B) {
	b.Run("Simple increment", func(b *testing.B) {
		var count uint32
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			count++
		}
	})

	b.Run("mo.Some assignment", func(b *testing.B) {
		var moCount = mo.Some(0)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			value, _ := moCount.Get()
			moCount = mo.Some(value + 1)
		}
	})
}

// If you need to run this outside of the testing framework:
func manualBenchmark() {
	// Version 1: Simple increment
	simpleTimer := time.NewTimer(time.Second)
	iterations := 0
	for {
		select {
		case <-simpleTimer.C:
			fmt.Printf("Simple increment: %d iterations in 1 second\n", iterations)
			goto nextBenchmark
		default:
			job := JobRequest{ModelId: "SDXL"}
			job.SetDefaults()
			// if _, exists := benchmarkJob.Prompt.Get(); !exists {
			// 	benchmarkJob.Prompt = mo.None[string]()
			// }
			// if benchmarkJob.NumImages == 0 {
			// 	benchmarkJob.NumImages = 4
			// }
			// if benchmarkJob.AspectRatio == "" {
			// 	benchmarkJob.AspectRatio = "1/1"
			// }
			iterations++
		}
	}

nextBenchmark:
	// Version 2: Using mo.Some
	moTimer := time.NewTimer(time.Second)
	iterations = 0
	for {
		select {
		case <-moTimer.C:
			fmt.Printf("mo.Some assignment: %d iterations in 1 second\n", iterations)
			return
		default:
			job := JobRequest{ModelId: "SDXL"}
			defaults.MustSet(&job)
			iterations++
		}
	}
}
