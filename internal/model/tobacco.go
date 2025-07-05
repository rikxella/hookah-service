package model

import (
	"encoding/csv"
	"io"
	"os"
	"strings"
)

type Tobacco struct {
	BrandName  string
	BrandURL   string
	FlavorName string
	FlavorURL  string
}

type Brand struct {
	BrandName string
}

type SearchIndex struct {
	BrandPrefixes map[string][]string             // префикс -> список брендов
	BrandToNames  map[string]map[string][]*Tobacco // бренд -> (префикс -> табаки)
}

func BrandLoadFromCSV(path string) ([]*Brand, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var brands []*Brand
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		brands = append(brands, &Brand{
			BrandName:  record[0],
		})
	}

	return brands, nil
}

func TobaccoLoadFromCSV(path string) ([]*Tobacco, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var tobaccos []*Tobacco
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		tobaccos = append(tobaccos, &Tobacco{
			BrandName:  record[0],
			BrandURL:   record[1],
			FlavorName: record[2],
			FlavorURL:  record[3],
		})
	}

	return tobaccos, nil
}

func BuildIndex(tobaccos []*Tobacco, brands []*Brand) *SearchIndex {
	index := &SearchIndex{
		BrandPrefixes: make(map[string][]string),
		BrandToNames:  make(map[string]map[string][]*Tobacco),
	}

	brandMap := make(map[string]struct{})
	for _, t := range brands {
		brandMap[t.BrandName] = struct{}{}
	}

	for brand := range brandMap {
		lowerBrand := strings.ToLower(brand)
		for i := 2; i <= len(lowerBrand); i++ {
			prefix := lowerBrand[:i]
			index.BrandPrefixes[prefix] = append(index.BrandPrefixes[prefix], brand)
		}
	}

	for _, t := range tobaccos {
		brand := t.BrandName
		lowerName := strings.ToLower(t.FlavorName)

		if _, exists := index.BrandToNames[brand]; !exists {
			index.BrandToNames[brand] = make(map[string][]*Tobacco)
		}

		for i := 2; i <= len(lowerName); i++ {
			prefix := lowerName[:i]
			index.BrandToNames[brand][prefix] = append(
				index.BrandToNames[brand][prefix], t)
		}
	}

	return index
}