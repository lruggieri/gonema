package elastic

import (
	"github.com/olivere/elastic/v7"
)

type Connection struct{
	connection *elastic.Client
}