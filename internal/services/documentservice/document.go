package documentservice

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"

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
	document, extractErr := ds.extractData(repo)
	if extractErr != nil {
		return nil, extractErr
	}

	repoURL := ds.gitHubService.GetRepoURL(repo)
	properties, propErr := ds.gitHubService.GetRepoProperties(repo)
	if propErr != nil {
		return nil, propErr
	}

	documentLinks := make(map[string]string)
	ds.processDocuments(document.Nav, documentLinks, repoURL, properties["DefaultBranch"], "")

	return documentLinks, nil
}

func (ds *DocumentService) extractData(repo string) (dtos.Document, error) {
	fileContent, fileErr := ds.getRemoteDocument(repo)
	if fileErr != nil {
		return dtos.Document{}, fileErr
	}

	decoder := yaml.NewDecoder(bytes.NewReader([]byte(fileContent)))
	var result dtos.Document
	for {
		decodeErr := decoder.Decode(&result)
		if decodeErr != nil {
			if decodeErr == io.EOF {
				break
			}
			return dtos.Document{}, decodeErr
		}
	}

	return result, nil
}

func (ds *DocumentService) getRemoteDocument(repo string) (string, error) {
	// Let's assume the standard is to use the docs folder
	fileContent, docsFileErr := ds.gitHubService.GetFileContent(repo, "docs/mkdocs.yaml")
	if docsFileErr == nil {
		return fileContent, nil
	}

	// Fallback to the root folder
	rootFileContent, rootFileErr := ds.gitHubService.GetFileContent(repo, "mkdocs.yaml")
	if rootFileErr == nil {
		return rootFileContent, nil
	}

	// Fallback to the .of folder
	ofFileContent, ofFileErr := ds.gitHubService.GetFileContent(repo, ".of/mkdocs.yaml")
	if ofFileErr == nil {
		return ofFileContent, nil
	}

	return "", errors.New("error getting file content from remote repository looking for mkdocs.yaml or docs/mkdocs.yaml")
}

func (ds *DocumentService) processDocuments(
	docs []dtos.NavItem,
	documentLinks map[string]string,
	repoURL, defaultBranch, parentName string,
) {
	for _, doc := range docs {
		var title string
		if parentName == "" {
			title = doc.Title
		} else {
			title = fmt.Sprintf("%s/%s", parentName, doc.Title)
		}

		if len(doc.SubItems) > 0 {
			ds.processDocuments(doc.SubItems, documentLinks, repoURL, defaultBranch, title)
			continue
		}

		documentLinks[title] = fmt.Sprintf("%s/blob/%s/docs/%s", repoURL, defaultBranch, doc.File)
	}
}
