// Format:
// [entity] [entity_name] {
// 	[attribute_name] [attribute_type]
// }

// Tokenizer -> Parser -> Generator

package main

import (
	"io"
	"reflect"
	"strings"
	"sync"
	"text/scanner"
	"unicode"

	"github.com/jaswdr/faker"
)

const RECORD_COUNT = 10
const IMAGE_METHOD_KEY = "Person.Image.Name.Image.Name"

type Attribute struct {
	Name string
	Type string
}

type Entity struct {
	Name       string
	Attributes []Attribute
}

// TODO: should be refactored
func parseSource(src io.Reader) []Entity {
	var s scanner.Scanner
	s.Init(src)
	s.Filename = "input"
	s.Mode ^= scanner.SkipComments

	// Allow dot in identifiers
	s.IsIdentRune = func(ch rune, i int) bool {
		return ch == '.' && i > 0 || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
	}

	var entities []Entity
	var currentEntity *Entity
	var currentAttribute *Attribute

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		text := s.TokenText()
		switch {
		case text == "entity":
			if currentEntity != nil {
				entities = append(entities, *currentEntity)
			}
			currentEntity = &Entity{}
		case text == "{":
			// Start of an entity's attributes
		case text == "}":
			// End of an entity's attributes
			if currentEntity != nil && currentAttribute != nil {
				currentEntity.Attributes = append(currentEntity.Attributes, *currentAttribute)
				currentAttribute = nil
			}
		case currentEntity != nil && currentEntity.Name == "":
			currentEntity.Name = text
		case currentAttribute == nil:
			currentAttribute = &Attribute{Name: text}
		case currentAttribute != nil && currentAttribute.Type == "":
			currentAttribute.Type = text

			typeParts := strings.Split(currentAttribute.Type, ".")

			if len(typeParts) == 2 {
				faker := faker.New()
				fakerValue := reflect.ValueOf(faker)

				fakerMethod := fakerValue.MethodByName(typeParts[0])

				if !fakerMethod.IsValid() {
					panic("Invalid type at attribute " + currentAttribute.Name + " for " + typeParts[0])
				}
			}

			if currentEntity != nil {
				currentEntity.Attributes = append(currentEntity.Attributes, *currentAttribute)
				currentAttribute = nil
			}
		}
	}

	if currentEntity != nil {
		entities = append(entities, *currentEntity)
	}

	return entities
}

type entityResult struct {
	EntityName string
	Data       *[]interface{}
}

// TODO: Map the all available funcitons instead of using reflection. Reflection is performance killer.
func generateFake(entities []Entity) map[string]*[]interface{} {
	resultList := make(map[string]*[]interface{})
	fake := faker.New()
	methodCache := make(map[string]reflect.Value)
	fakerValue := reflect.ValueOf(fake)

	var wg sync.WaitGroup
	resultsChan := make(chan entityResult, len(entities))

	for _, entity := range entities {
		wg.Add(1)
		go func(entity Entity) {
			defer wg.Done()
			lowercasedEntityName := strings.ToLower(entity.Name)
			entityResults := make([]interface{}, 0, RECORD_COUNT)

			for i := 0; i < RECORD_COUNT; i++ {
				fakeValues := make(map[string]any)

				for _, attr := range entity.Attributes {
					methodKey := attr.Type
					method, exists := methodCache[methodKey]

					if !exists {
						typeParts := strings.Split(attr.Type, ".")
						method = fakerValue.MethodByName(typeParts[0])

						for _, part := range typeParts[1:] {
							method = method.Call([]reflect.Value{})[0].MethodByName(part)
							methodKey += "." + part
						}

						methodCache[methodKey] = method
					}

					result := method.Call([]reflect.Value{})[0].Interface()

					if methodKey == IMAGE_METHOD_KEY {
						result = strings.Split(result.(string), "/")[2]
					}

					fakeValues[attr.Name] = result
				}

				entityResults = append(entityResults, fakeValues)
			}

			resultsChan <- entityResult{EntityName: lowercasedEntityName, Data: &entityResults}
		}(entity)
	}

	// Close the channel after all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	for result := range resultsChan {
		resultList[result.EntityName] = result.Data
	}

	return resultList
}

func Generate(
	src io.Reader,
) map[string]*[]interface{} {
	entities := parseSource(src)
	return generateFake(entities)
}
