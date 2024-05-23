package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "strings"

    "cloud.google.com/go/asset/apiv1"
    assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1"
    "google.golang.org/api/iterator"
    "google.golang.org/api/option"
    "golang.org/x/oauth2/google"
)

type Resource struct {
    Name       string `json:"name"`
    AssetType  string `json:"asset_type"`
    Project    string `json:"project"`
    ResourceID string `json:"resource_id"`
}

func main() {
    http.HandleFunc("/assets", assetsHandler)
    http.Handle("/", http.FileServer(http.Dir("./static")))
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func assetsHandler(w http.ResponseWriter, r *http.Request) {
    assetType := r.URL.Query().Get("type")
    assets, err := getAssets()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if assetType != "" {
        filteredAssets := []Resource{}
        for _, asset := range assets {
            if strings.Contains(asset.AssetType, assetType) {
                filteredAssets = append(filteredAssets, asset)
            }
        }
        assets = filteredAssets
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(assets)
}

func getAssets() ([]Resource, error) {
    ctx := context.Background()

    defaultCredentials, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/cloud-platform")
    if err != nil {
        return nil, err
    }

    assetClient, err := asset.NewClient(ctx, option.WithCredentials(defaultCredentials))
    if err != nil {
        return nil, err
    }
    defer assetClient.Close()

    searchRequest := &assetpb.SearchAllResourcesRequest{
        Scope: "projects/deans-playground",
    }

    resourceIterator := assetClient.SearchAllResources(ctx, searchRequest)
    var resources []Resource

    for {
        resourceResponse, err := resourceIterator.Next()
        if err == iterator.Done {
            break
        }
        if err != nil {
            return nil, err
        }

        // Extract the name from the full resource path
        parts := strings.Split(resourceResponse.Name, "/")
        name := parts[len(parts)-1]

        resource := Resource{
            Name:       name,
            AssetType:  resourceResponse.AssetType,
            Project:    "deans-playground",
            ResourceID: resourceResponse.AssetType,
        }
        resources = append(resources, resource)
    }

    return resources, nil
}

