package pingback

import (
	"fmt"
	"log"
	"strings"

	"github.com/knakk/rdf"
	"github.com/rs/xid"
	"opencoredata.org/ocdGarden/CSDCO/VaultWalker/pkg/utils"
)

// RDFGraph (item, shaval, *rdf)
// In this approach each object gets a named graph.  Perhaps this is not
// needed since each data graph also has a sha ID with it?  Which is all we really
// use in the graph IRI.   ???
func pingtograph(item string, shaval string, ub *utils.Buffer) int {
	var b strings.Builder

	//t := utils.MimeByType(item)
	newctx, _ := rdf.NewIRI(fmt.Sprintf("http://opencoredata.org/objectgraph/id/%s", shaval))
	ctx := rdf.Context(newctx)

	guid := xid.New()
	s := fmt.Sprintf("http://opencoredata.org/id/do/%s", guid)
	//d := fmt.Sprintf("http://opencoredata.org/id/dx/%s", guid) // distribution URL

	_ = iiTriple(s, "http://www.w3.org/1999/02/22-rdf-syntax-ns#type", item, ctx, &b)

	len, err := ub.Write([]byte(b.String()))
	if err != nil {
		log.Printf("error in the buffer write... %v\n", err)
	}

	return len //  we will return the bytes count we write...
}

func iiTriple(s, p, o string, c rdf.Context, b *strings.Builder) error {
	sub, err := rdf.NewIRI(s)
	pred, err := rdf.NewIRI(p)
	obj, err := rdf.NewIRI(o)

	t := rdf.Triple{Subj: sub, Pred: pred, Obj: obj}
	q := rdf.Quad{t, c}

	qs := q.Serialize(rdf.NQuads)
	fmt.Fprintf(b, "%s", qs)
	return err
}

func ilTriple(s, p, o string, c rdf.Context, b *strings.Builder) error {
	sub, err := rdf.NewIRI(s)
	pred, err := rdf.NewIRI(p)
	obj, err := rdf.NewLiteral(o)

	t := rdf.Triple{Subj: sub, Pred: pred, Obj: obj}
	q := rdf.Quad{t, c}

	qs := q.Serialize(rdf.NQuads)
	fmt.Fprintf(b, "%s", qs)
	return err
}
