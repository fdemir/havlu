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
	"text/scanner"
	"unicode"

	"github.com/jaswdr/faker"
)

const RECORD_COUNT = 10

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

			// is valid faker method

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

func generateFake(entities []Entity) map[string]*[]interface{} {
	resultList := make(map[string]*[]interface{})

	fake := faker.New()

	methodCache := make(map[string]reflect.Value)
	fakerValue := reflect.ValueOf(fake)

	for _, entity := range entities {
		entityResults := make([]interface{}, 0, RECORD_COUNT)
		resultList[entity.Name] = &entityResults

		for i := 0; i < RECORD_COUNT; i++ {
			fakeValues := make(map[string]any)

			for _, attr := range entity.Attributes {
				methodKey := attr.Type
				method, exists := methodCache[methodKey]

				if !exists {
					typeParts := strings.Split(attr.Type, ".")
					fakerMethod := fakerValue.MethodByName(typeParts[0])
					method = fakerMethod.Call([]reflect.Value{})[0].MethodByName(typeParts[1])
					methodCache[methodKey] = method
				}

				fakeValues[attr.Name] = method.Call([]reflect.Value{})[0].Interface()
			}

			*resultList[entity.Name] = append(*resultList[entity.Name], fakeValues)
		}
	}

	return resultList
}

func Generate(
	src io.Reader,
) map[string]*[]interface{} {
	entities := parseSource(src)
	return generateFake(entities)
}
