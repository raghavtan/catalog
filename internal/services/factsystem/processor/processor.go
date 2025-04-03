package processor

import (
	"context"
	"fmt"
	"sync"

	"github.com/motain/of-catalog/internal/services/factsystem/aggregators"
	"github.com/motain/of-catalog/internal/services/factsystem/dtos"
	"github.com/motain/of-catalog/internal/services/factsystem/extractors"
	"github.com/motain/of-catalog/internal/services/factsystem/validators"
	"github.com/motain/of-catalog/internal/utils/transformers"
)

type ProcessorInterface interface {
	Process(ctx context.Context, tasks []*dtos.Task) (float64, error)
}

type Processor struct {
	Mu         sync.RWMutex
	Aggregator aggregators.AggregatorInterface
	Validator  validators.ValidatorInterface
	Extractor  extractors.ExtractorInterface
}

func NewProcessor(
	aggregator aggregators.AggregatorInterface,
	validator validators.ValidatorInterface,
	extractor extractors.ExtractorInterface,
) *Processor {
	return &Processor{
		Aggregator: aggregator,
		Validator:  validator,
		Extractor:  extractor,
	}
}

func (p *Processor) Process(ctx context.Context, tasks []*dtos.Task) (float64, error) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))

	mappedTasks := make(map[string]*dtos.Task)
	for _, task := range tasks {
		mappedTasks[task.ID] = task
	}

	var result interface{}
	for _, task := range tasks {
		task.DoneCh = make(chan dtos.TaskResult, 1)
		for _, dependsOn := range task.DependsOn {
			if _, ok := mappedTasks[dependsOn]; !ok {
				continue
			}
			task.Dependencies = append(task.Dependencies, mappedTasks[dependsOn])
		}

		go p.execute(ctx, task, &wg, &result)
	}

	wg.Wait()

	p.Mu.RLock()
	defer p.Mu.RUnlock()

	// Grab the result from the last task and try to convert it to a float64
	return transformers.Interface2Float64(result)
}

func (p *Processor) execute(ctx context.Context, task *dtos.Task, wg *sync.WaitGroup, result *interface{}) {
	defer wg.Done()
	defer close(task.DoneCh)

	// Wait for all dependencies to finish
	for _, dep := range task.Dependencies {
		<-dep.DoneCh
	}

	switch dtos.TaskType(task.Type) {
	case dtos.ExtractType:
		extractErr := p.handleExtract(ctx, task)
		if extractErr != nil {
			fmt.Printf("%s: error extracting data: %v\n", task.ID, extractErr)
		}
	case dtos.ValidateType:
		validationErr := p.handleValidate(task)
		if validationErr != nil {
			fmt.Printf("%s: error validating data: %v\n", task.ID, validationErr)
		}
	case dtos.AggregateType:
		aggregateErr := p.handleAggregate(ctx, task)
		if aggregateErr != nil {
			fmt.Printf("%s: error aggregating data: %v\n", task.ID, aggregateErr)
		}
	default:
		fmt.Printf("%s: unknown task type: %s\n", task.ID, task.Type)
	}

	p.Mu.Lock()
	defer p.Mu.Unlock()

	*result = task.Result
	task.DoneCh <- dtos.TaskResult{Result: task.ID}
}

func (p *Processor) handleExtract(ctx context.Context, task *dtos.Task) error {
	var deps []*dtos.Task
	for _, dep := range task.Dependencies {
		if dtos.TaskType(dep.Type) == dtos.ExtractType {
			deps = append(deps, dep)
		}
	}

	extractErr := p.Extractor.Extract(ctx, task, deps)
	if extractErr != nil {
		return extractErr
	}

	return nil
}

func (p *Processor) handleValidate(task *dtos.Task) error {
	err := p.Validator.Check(task, task.Dependencies)
	if err != nil {
		return fmt.Errorf("%s: error validating data: %v", task.ID, err)
	}

	return nil
}

func (p *Processor) handleAggregate(ctx context.Context, task *dtos.Task) error {
	err := p.Aggregator.Combine(ctx, task, task.Dependencies)
	if err != nil {
		return fmt.Errorf("%s: error aggregating data: %v", task.ID, err)
	}

	return nil
}
