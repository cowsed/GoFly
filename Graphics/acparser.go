package graphics

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/go-gl/mathgl/mgl32"
)

type ACModel struct {
	version   int
	materials []ACMat
	obj       ACObj
}
type ACObj struct {
	name      string
	modeltype string

	loc     mgl32.Vec3
	mesh    ACMesh
	numkids int
	kids    []ACObj
}
type ACMat struct {
	name  string
	rgb   mgl32.Vec3
	amb   mgl32.Vec3
	emis  mgl32.Vec3
	spec  mgl32.Vec3
	shi   int
	trans float32
}
type ACMesh struct {
	verts []mgl32.Vec3
	faces []ACFace
}
type ACFace struct {
	vertIndices []int
	matIndex    int
}

func LoadACFile(fname string) (*ACModel, error) {
	f1, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	src, err := io.ReadAll(f1)
	if err != nil {
		return nil, err
	}
	m, err := ParseACFile(string(src))
	if err != nil {
		return nil, err
	}
	return m, err
}

func ParseACFile(src string) (*ACModel, error) {
	m := ACModel{}
	m.materials = []ACMat{}
	src = strings.ReplaceAll(src, "\r", "")
	lines := strings.Split(src, "\n")

	header := lines[0]
	if header[0:4] != "AC3D" {
		return nil, fmt.Errorf("incorrect AC3D header")
	}

	m.version = HexChartoInt(header[4])
	//log.Println("ver", m.version)
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		parts := strings.Split(line, " ")
		tokeni := parts[0]
		switch tokeni {
		case "MATERIAL":
			mat := ACMat{}
			fmt.Sscanf(line, "MATERIAL %s rgb %f %f %f  amb %f %f %f  emis %f %f %f  spec %f %f %f  shi %d  trans %f",
				&mat.name, &mat.rgb[0], &mat.rgb[1], &mat.rgb[2], &mat.amb[0], &mat.amb[1], &mat.amb[2], &mat.emis[0], &mat.emis[1], &mat.emis[2],
				&mat.spec[0], &mat.spec[1], &mat.spec[2], &mat.shi, &mat.trans)
			mat.name = strings.ReplaceAll(mat.name, "\"", "")
			m.materials = append(m.materials, mat)
		case "OBJECT":

			var nextLineDelta int
			m.obj, nextLineDelta = ParseACObject(lines[i:])
			i += nextLineDelta
			//log.Println("===========THIS SHOULD BE THE END OF IT===========")
		}
	}
	return &m, nil
}

//Parses OBJECT...kids and turns it into an ACObj
//perhaps later pass in the index of the first line so that error reporting can say line x
func ParseACObject(lines []string) (ACObj, int) {
	o := ACObj{}
	fmt.Sscanf(lines[0], "OBJECT %s", &o.modeltype)
	for i := 1; i < len(lines); i++ {
		tokeni := strings.Split(lines[i], " ")[0]
		switch tokeni {
		case "OBJECT":
			child, nextLineDelta := ParseACObject(lines[i:])
			o.kids = append(o.kids, child)
			i += nextLineDelta
			//parsed all kids
			if len(o.kids) == o.numkids {
				return o, i
			}

		case "name":
			fmt.Sscanf(lines[i], "name %s", &o.name)
			//log.Println("Naming object", o.name)
		case "loc":
			fmt.Sscanf(lines[i], "loc %f %f %f", &o.loc[0], &o.loc[1], &o.loc[2])
			//log.Println("Getting location", o.loc, "for", o.name)
		case "numvert":
			//Start Mesh Parsing
			numverts := 0
			fmt.Sscanf(lines[i], "numvert %d", &numverts)
			o.mesh.verts = ParseMeshVerts(lines[i:i+numverts+1], numverts)
			//log.Println("Finished parsing verts for", o.name)
			i += numverts //Skip over vertices
		case "numsurf":
			faces, nextLineDelta := ParseMeshSurfaces(lines[i:])
			o.mesh.faces = faces
			i += nextLineDelta - 1
			//log.Println("parsed surfaces of ", o.name)

		case "kids":
			fmt.Sscanf(lines[i], "kids %d", &o.numkids)
			//log.Println("Ending Object", o.name)
			if o.numkids == 0 {
				//log.Println("No Kids")
				return o, i
			}
			//Multiple kids

			//not all kids have been parsed, keep parsing
			//log.Println("Hit Kids")

		default:
			log.Println("Omitting line: ", lines[i])
		}
	}
	//log.Println("***********Last Object**********", o.name)
	return o, len(lines)
}

//Parses the surfaces of an ac file and returns where it ends (index of last line kids)
func ParseMeshSurfaces(lines []string) ([]ACFace, int) {
	numsurf := 0
	fmt.Sscanf(lines[0], "numsurf %d", &numsurf)
	faces := make([]ACFace, numsurf)
	faceIndex := -1
	numIndices := 0
	indexIndex := 0

	for i := 1; i < len(lines); i++ {
		tokeni := strings.Split(lines[i], " ")[0]
		switch tokeni {
		case "kids":
			//End of the line, return
			return faces, i
		case "SURF":
			faceIndex++
			indexIndex = 0

		case "mat":
			fmt.Sscanf(lines[i], "mat %d", &faces[faceIndex].matIndex)

		case "refs":
			fmt.Sscanf(lines[i], "refs %d", &numIndices)
			faces[faceIndex].vertIndices = make([]int, numIndices)

		default: //is an index specification
			var vIndex int
			var texInd0, texInd1 int
			//tex coordinates
			fmt.Sscanf(lines[i], "%d %d %d", &vIndex, &texInd0, &texInd1)
			faces[faceIndex].vertIndices[indexIndex] = vIndex
			indexIndex++
		}

	}
	log.Println("Should never get here face version")
	return faces, len(lines)
}
func ParseMeshVerts(lines []string, numverts int) []mgl32.Vec3 {
	//First line should be numverts
	verts := make([]mgl32.Vec3, numverts)
	for i := 0; i < numverts; i++ {
		line := lines[i+1]
		fmt.Sscanf(line, "%f %f %f", &verts[i][0], &verts[i][1], &verts[i][2])
	}
	return verts
}

func HexChartoInt(c byte) int {
	bs, _ := hex.DecodeString("0" + string(c))
	return int(bs[0])
}
