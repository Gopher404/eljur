package schedules

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

const token = "vk1.a.A4lLFx-zbGgGOAYN8IIxoZPSpJ4_Ye4qj4psKCvR7RAkZNrhHD713uz3JtFyd6oQpVXQdWCe7N-USI2SaosdrjIZAXmfO1PB9_DbL0Z6JbS0WA0BlzM8VpmAVUHiwl0ob4-qtvPI_zODFYDPM10flVNvJc5AqPspsZymVGH2G6mzoqa7fKeMzKic-CLgS64AWSJzuApx2Te1d-LtqqXMog"

func TestGetWeekFromDocument(t *testing.T) {
	parser := newParser(-193901024, token)
	docs, err := parser.getListDocuments()
	require.NoError(t, err, "get list docs error")

	doc, err := parser.getDocument(docs[0])
	require.NoError(t, err, "get document error")

	week := parser.getWeekFromDocument(doc)
	fmt.Println("week:", week)
}

func TestGetChangesFromDocument(t *testing.T) {

	parser := newParser(-193901024, token)
	//testS := "Код будущего Куликовский М.Ю. Инф. Куликовский М.Ю. Минеев-Ли В.Е. 302,305"
	//_, _, i1, i2 := parser.extractNameS(testS)
	//fmt.Println(testS[i1:i2])
	//return

	docs, err := parser.getListDocuments()
	require.NoError(t, err, "get list docs error")

	doc, err := parser.getDocument(docs[0])
	require.NoError(t, err, "get document error")

	changes := parser.getChangesFromDocument(doc, "ИС11")
	fmt.Printf("changes: %+v", changes)
}
