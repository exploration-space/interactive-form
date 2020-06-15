package storage

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gofrs/uuid"
	"golang.org/x/net/context"
)

type azureBlob struct {
	key, accountName, serviceEndpoint, container string
}

func (s azureBlob) Connect() error {
	return nil
}

func (s azureBlob) Listdir() []string {
	files := []string{}
	return files
}

func (s azureBlob) Touch() error {
	return nil
}

func (s azureBlob) Upload(content []byte) error {
	url, err := uploadBytesToBlob(content, s.key, s.accountName, s.serviceEndpoint, s.container)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(url)
	fmt.Println(string(content))
	return nil
}

func AzureBlobStorage(key string, accountName string, container string) azureBlob {
	_, errC := azblob.NewSharedKeyCredential(accountName, key)
	if errC != nil {
		fmt.Println("cred", errC)
		panic("credential issues")
	}

	serviceEndpoint := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)
	s := azureBlob{
		key,
		accountName,
		serviceEndpoint,
		container,
	}
	return s
}

func readFile(filePath string) ([]byte, error) {
	dat, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	} else {
		return dat, nil
	}
}

func uploadBytesToBlob(b []byte, key string, accountName string, endPoint string, container string) (string, error) {
	u, _ := url.Parse(fmt.Sprint(endPoint, container, "/", getBlobName()))
	credential, errC := azblob.NewSharedKeyCredential(accountName, key)
	if errC != nil {
		fmt.Println("here", errC)
		return "", errC
	}

	blockBlobUrl := azblob.NewBlockBlobURL(*u, azblob.NewPipeline(credential, azblob.PipelineOptions{}))

	ctx := context.Background()
	o := azblob.UploadToBlockBlobOptions{
		BlobHTTPHeaders: azblob.BlobHTTPHeaders{
			ContentType: "application/json",
		},
	}

	_, errU := azblob.UploadBufferToBlockBlob(ctx, b, blockBlobUrl, o)
	fmt.Println("errU", errU)
	return blockBlobUrl.String(), errU
}

func getBlobName() string {
	t := time.Now()
	uuid, _ := uuid.NewV4()

	return fmt.Sprintf("%s-%v.json", t.Format("20060102"), uuid)
}