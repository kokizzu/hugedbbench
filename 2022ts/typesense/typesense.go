package typesense

import (
	"time"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

func main() {
	client := typesense.NewClient(
		typesense.WithServer("http://localhost:8108"),
		typesense.WithAPIKey("local-typesense-api-key"),
		typesense.WithConnectionTimeout(5*time.Second),
		typesense.WithCircuitBreakerMaxRequests(50),
		typesense.WithCircuitBreakerInterval(2*time.Minute),
		typesense.WithCircuitBreakerTimeout(1*time.Minute),
	)

	//Create a collection

	yes := true
	sortField := "num_employees"
	schema := &api.CollectionSchema{
		Name: "companies",
		Fields: []api.Field{
			{
				Name: "company_name",
				Type: "string",
			},
			{
				Name: "num_employees",
				Type: "int32",
			},
			{
				Name:  "country",
				Type:  "string",
				Facet: &yes,
			},
		},
		DefaultSortingField: &sortField,
	}

	client.Collections().Create(schema)

	//Index a document

	document := struct {
		ID           string `json:"id"`
		CompanyName  string `json:"company_name"`
		NumEmployees int    `json:"num_employees"`
		Country      string `json:"country"`
	}{
		ID:           "123",
		CompanyName:  "Stark Industries",
		NumEmployees: 5215,
		Country:      "USA",
	}

	client.Collection("companies").Documents().Create(document)

	//Upserting a document

	newDocument := struct {
		ID           string `json:"id"`
		CompanyName  string `json:"company_name"`
		NumEmployees int    `json:"num_employees"`
		Country      string `json:"country"`
	}{
		ID:           "123",
		CompanyName:  "Stark Industries",
		NumEmployees: 5215,
		Country:      "USA",
	}

	client.Collection("companies").Documents().Upsert(newDocument)

	//Search a collection

	filterBy := "num_employees:>100"
	sortBy := "num_employees:desc"
	searchParameters := &api.SearchCollectionParams{
		Q:        "stark",
		QueryBy:  "company_name",
		FilterBy: &filterBy,
		SortBy:   &sortBy,
	}

	client.Collection("companies").Documents().Search(searchParameters)

	//Retrieve a document

	client.Collection("companies").Document("123").Retrieve()

	//Update a document

	document2 := struct {
		CompanyName  string `json:"company_name"`
		NumEmployees int    `json:"num_employees"`
	}{
		CompanyName:  "Stark Industries",
		NumEmployees: 5500,
	}

	client.Collection("companies").Document("123").Update(document2)

	//Delete an individual document

	client.Collection("companies").Document("123").Delete()

	//Delete a bunch of documents

	filterBy = "num_employees:>100"
	batchSize := 100
	filter := &api.DeleteDocumentsParams{FilterBy: &filterBy, BatchSize: &batchSize}
	client.Collection("companies").Documents().Delete(filter)

	//Retrieve a collection

	client.Collection("companies").Retrieve()

	//Export documents from a collection

	client.Collection("companies").Documents().Export()

	//Import an array of documents:

	documents := []interface{}{
		struct {
			ID           string `json:"id"`
			CompanyName  string `json:"companyName"`
			NumEmployees int    `json:"numEmployees"`
			Country      string `json:"country"`
		}{
			ID:           "123",
			CompanyName:  "Stark Industries",
			NumEmployees: 5215,
			Country:      "USA",
		},
	}
	action := "create"
	batchSize = 40
	params := &api.ImportDocumentsParams{
		Action:    &action,
		BatchSize: &batchSize,
	}

	client.Collection("companies").Documents().Import(documents, params)

	//Import a JSONL file:
	/*
		params := &api.ImportDocumentsParams{
			Action:    "create",
			BatchSize: 40,
		}
		importBody, err := os.Open("documents.jsonl")
		// defer close, error handling ...

		client.Collection("companies").Documents().ImportJsonl(importBody, params)
	*/
	//List all collections

	client.Collections().Retrieve()

	//Drop a collection

	client.Collection("companies").Delete()

	//Create an API Key

	description := "Search-only key."
	expiresAt := time.Now().AddDate(0, 6, 0).Unix()
	keySchema := &api.ApiKeySchema{
		Description: &description,
		Actions:     []string{"documents:search"},
		Collections: []string{"companies"},
		ExpiresAt:   &expiresAt,
	}

	client.Keys().Create(keySchema)

	//Retrieve an API Key

	client.Key(1).Retrieve()

	//List all keys

	client.Keys().Retrieve()

	//Delete API Key

	client.Key(1).Delete()

	//Create or update an override

	override := &api.SearchOverrideSchema{
		Rule: api.SearchOverrideRule{
			Query: "apple",
			Match: "exact",
		},
		Includes: &[]api.SearchOverrideInclude{
			{
				Id:       "422",
				Position: 1,
			},
			{
				Id:       "54",
				Position: 2,
			},
		},
		Excludes: &[]api.SearchOverrideExclude{
			{
				Id: "287",
			},
		},
	}

	client.Collection("companies").Overrides().Upsert("customize-apple", override)

	//List all overrides

	client.Collection("companies").Overrides().Retrieve()

	//Delete an override

	client.Collection("companies").Override("customize-apple").Delete()

	//Create or Update an alias

	body := &api.CollectionAliasSchema{CollectionName: "companies_june11"}
	client.Aliases().Upsert("companies", body)

	//Retrieve an alias

	client.Alias("companies").Retrieve()

	//List all aliases

	client.Aliases().Retrieve()

	//Delete an alias

	client.Alias("companies").Delete()

	//Create or update a multi-way synonym

	synonym := &api.SearchSynonymSchema{
		Synonyms: []string{"blazer", "coat", "jacket"},
	}
	client.Collection("products").Synonyms().Upsert("coat-synonyms", synonym)

	//Create or update a one-way synonym

	root := "blazer"
	synonym = &api.SearchSynonymSchema{
		Root:     &root,
		Synonyms: []string{"blazer", "coat", "jacket"},
	}
	client.Collection("products").Synonyms().Upsert("coat-synonyms", synonym)

	//Retrieve a synonym

	client.Collection("products").Synonym("coat-synonyms").Retrieve()

	//List all synonyms

	client.Collection("products").Synonyms().Retrieve()

	//Delete a synonym

	client.Collection("products").Synonym("coat-synonyms").Delete()

	//Create snapshot (for backups)

	client.Operations().Snapshot("/tmp/typesense-data-snapshot")

	//Re-elect Leader

	client.Operations().Vote()
}
