package documentservice

import (
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"path/filepath"

	"github.com/motain/of-catalog/internal/services/documentservice/dtos"
	"github.com/motain/of-catalog/internal/services/githubservice"
)

type DocumentServiceInterface interface {
	GetDocuments(repo string) (map[string]string, error)
}

type DocumentService struct {
	gitHubService githubservice.GitHubServiceInterface
}

func NewDocumentService(gitHubService githubservice.GitHubServiceInterface) *DocumentService {
	return &DocumentService{
		gitHubService: gitHubService,
	}
}

func (ds *DocumentService) GetDocuments(repo string) (map[string]string, error) {
	documentLinks := make(map[string]string)
	document, indexLocation, extractErr := ds.extractData(repo)
	if extractErr == nil {
		repoURL := ds.gitHubService.GetRepoURL(repo)
		properties, propErr := ds.gitHubService.GetRepoProperties(repo)
		if propErr == nil {
			uriToDocFile := repoURL + filepath.Join(string(filepath.Separator), "blob", properties["DefaultBranch"], indexLocation, "docs")
			ds.processDocuments(document.Nav, documentLinks, uriToDocFile, "")
		}
	}
	readmeLinks, readmeErr := ds.getReadmeDocuments(repo)
	if readmeErr == nil {
		for title, url := range readmeLinks {
			documentLinks[title] = url
		}
	}

	return documentLinks, nil
}

func (ds *DocumentService) getReadmeDocuments(repo string) (map[string]string, error) {
	readmeFiles := []string{
		"docs/README.md",
		"README.md",
		"index.md",
		"readme.md",
		"docs/readme.md",
		"docs/index.md",
	}

	repoURL := ds.gitHubService.GetRepoURL(repo)
	properties, propErr := ds.gitHubService.GetRepoProperties(repo)
	if propErr != nil {
		return nil, fmt.Errorf("failed to get repo properties: %w", propErr)
	}

	defaultBranch := properties["DefaultBranch"]
	documentLinks := make(map[string]string)

	for _, filePath := range readmeFiles {
		// Try to get the file content to check if it exists
		_, fileErr := ds.gitHubService.GetFileContent(repo, filePath)
		if fileErr == nil {
			// File exists, create the document entry
			title := ds.generateReadmeTitle(filePath)
			url := fmt.Sprintf("%s/blob/%s/%s", repoURL, defaultBranch, filePath)
			documentLinks[title] = url
		}
	}

	// Return error only if no README files were found
	if len(documentLinks) == 0 {
		return nil, errors.New("no README files found")
	}

	return documentLinks, nil
}

func (ds *DocumentService) generateReadmeTitle(filePath string) string {
	return "README" // Simplified title generation for README files
}

func (ds *DocumentService) extractData(repo string) (dtos.Document, string, error) {
	fileContent, indexLocation, fileErr := ds.getRemoteDocument(repo)
	if fileErr != nil {
		return dtos.Document{}, "", fileErr
	}

	decoder := yaml.NewDecoder(bytes.NewReader([]byte(fileContent)))
	var result dtos.Document
	for {
		decodeErr := decoder.Decode(&result)
		if decodeErr != nil {
			if decodeErr == io.EOF {
				break
			}
			return dtos.Document{}, "", decodeErr
		}
	}

	return result, indexLocation, nil
}

func (ds *DocumentService) getRemoteDocument(repo string) (string, string, error) {
	const indexFile = "mkdocs.yaml"
	possibleIndexLocations := []string{
		"",     // Let's assume the standard is to use the root folder
		"docs", // Fallback to the docs folder
		"doc",  // Fallback to the doc folder
		".of",  // Fallback to the .of folder
	}

	for _, folder := range possibleIndexLocations {
		fileContent, fileErr := ds.gitHubService.GetFileContent(repo, filepath.Join(folder, indexFile))
		if fileErr == nil {
			return fileContent, folder, nil
		}
	}

	return "", "", errors.New("error getting file content from remote repository")
}

func (ds *DocumentService) processDocuments(docs []dtos.NavItem, documentLinks map[string]string, uriToDocFile, parentName string) {
	for _, doc := range docs {
		title := doc.Title
		if parentName != "" {
			title = fmt.Sprintf("%s/%s", parentName, doc.Title)
		}

		if len(doc.SubItems) > 0 {
			ds.processDocuments(doc.SubItems, documentLinks, uriToDocFile, title)
			continue
		}

		documentLinks[title] = fmt.Sprintf("%s/%s", uriToDocFile, doc.File)
	}
}
