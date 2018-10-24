package primitive

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/otiai10/copy"
	"github.com/pkg/errors"

	"github.com/unchartedsoftware/distil-ingest/metadata"
	"github.com/unchartedsoftware/distil-ingest/pipeline"
	"github.com/unchartedsoftware/distil-ingest/primitive/compute"
	"github.com/unchartedsoftware/distil-ingest/primitive/compute/description"
	"github.com/unchartedsoftware/distil-ingest/primitive/compute/result"
	"github.com/unchartedsoftware/plog"
)

const (
	// D3MSchemaPathRelative is the standard name of the schema document.
	D3MSchemaPathRelative = "datasetDoc.json"
	// D3MDataPathRelative is the standard name of the data file.
	D3MDataPathRelative = "tables/learningData.csv"

	denormFieldName = "filename"
)

// FeatureRequest captures the properties of a request to a primitive.
type FeatureRequest struct {
	SourceVariableName  string
	FeatureVariableName string
	Variable            *metadata.Variable
	Step                *pipeline.PipelineDescription
}

type IngestStep struct {
	client *compute.Client
}

func NewIngestStep(client *compute.Client) *IngestStep {
	return &IngestStep{
		client: client,
	}
}

func (s *IngestStep) submitPrimitive(dataset string, step *pipeline.PipelineDescription) (string, error) {

	res, err := s.client.ExecutePipeline(context.Background(), dataset, step)
	if err != nil {
		return "", errors.Wrap(err, "unable to dispatch mocked pipeline")
	}
	resultURI := strings.Replace(res.ResultURI, "file://", "", -1)
	return resultURI, nil
}

func (s *IngestStep) readCSVFile(filename string, hasHeader bool) ([][]string, error) {
	// open the file
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open data file")
	}
	defer csvFile.Close()
	reader := csv.NewReader(csvFile)

	lines := make([][]string, 0)

	// skip the header as needed
	if hasHeader {
		_, err = reader.Read()
		if err != nil {
			return nil, errors.Wrap(err, "failed to read header from file")
		}
	}

	// read the raw data
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.Wrap(err, "failed to read line from file")
		}

		lines = append(lines, line)
	}

	return lines, nil
}

func (s *IngestStep) appendFeature(dataset string, d3mIndexField int, hasHeader bool, feature *FeatureRequest, lines [][]string) ([][]string, error) {
	datasetURI, err := s.submitPrimitive(dataset, feature.Step)
	if err != nil {
		return nil, errors.Wrap(err, "unable to run pipeline primitive")
	}
	log.Infof("parsing primitive result from '%s'", datasetURI)

	// parse primitive response (new field contains output)
	res, err := result.ParseResultCSV(datasetURI)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse pipeline primitive result")
	}

	// find the field with the feature output
	labelIndex := 1
	for i, f := range res[0] {
		if f == feature.FeatureVariableName {
			labelIndex = i
		}
	}

	// build the lookup for the new field
	features := make(map[string]string)
	for i, v := range res {
		// skip header
		if i > 0 {
			d3mIndex := v[0].(string)
			labels := v[labelIndex].(string)
			features[d3mIndex] = labels
		}
	}

	// add the new feature to the raw data
	for i, line := range lines {
		if i > 0 || !hasHeader {
			d3mIndex := line[d3mIndexField]
			feature := features[d3mIndex]
			line = append(line, feature)
			lines[i] = line
		}
	}

	return lines, nil
}

func getFeatureVariables(meta *metadata.Metadata, prefix string) ([]*FeatureRequest, error) {
	mainDR := meta.GetMainDataResource()
	features := make([]*FeatureRequest, 0)
	for _, v := range mainDR.Variables {
		if v.RefersTo != nil && v.RefersTo["resID"] != nil {
			// get the refered DR
			resID := v.RefersTo["resID"].(string)

			res := getDataResource(meta, resID)

			// check if needs to be featurized
			if res.CanBeFeaturized() {
				// create the new resource to hold the featured output
				indexName := fmt.Sprintf("%s%s", prefix, v.Name)

				// add the feature variable
				v := metadata.NewVariable(len(mainDR.Variables), indexName, "label", v.Name, "string", "string", "", "", []string{"attribute"}, metadata.VarRoleMetadata, nil, mainDR.Variables, false)

				// create the required pipeline
				step, err := description.CreateCrocPipeline("leather", "", []string{denormFieldName}, []string{indexName})
				if err != nil {
					return nil, errors.Wrap(err, "unable to create step pipeline")
				}

				features = append(features, &FeatureRequest{
					SourceVariableName:  denormFieldName,
					FeatureVariableName: indexName,
					Variable:            v,
					Step:                step,
				})
			}
		}
	}

	return features, nil
}

func getClusterVariables(meta *metadata.Metadata, prefix string) ([]*FeatureRequest, error) {
	mainDR := meta.GetMainDataResource()
	features := make([]*FeatureRequest, 0)
	for _, v := range mainDR.Variables {
		if v.RefersTo != nil && v.RefersTo["resID"] != nil {
			// get the refered DR
			resID := v.RefersTo["resID"].(string)

			res := getDataResource(meta, resID)

			// check if needs to be featurized
			if res.CanBeFeaturized() || res.ResType == "timeseries" {
				// create the new resource to hold the featured output
				indexName := fmt.Sprintf("%s%s", prefix, v.Name)

				// add the feature variable
				v := metadata.NewVariable(len(mainDR.Variables), indexName, "group", v.Name, "string", "string", "", "", []string{"attribute"}, metadata.VarRoleMetadata, nil, mainDR.Variables, false)

				// create the required pipeline
				var step *pipeline.PipelineDescription
				var err error
				if res.CanBeFeaturized() {
					step, err = description.CreateUnicornPipeline("horned", "", []string{denormFieldName}, []string{indexName})
				} else {
					step, err = description.CreateSlothPipeline("leaf", "", []string{denormFieldName}, []string{indexName})
				}
				if err != nil {
					return nil, errors.Wrap(err, "unable to create step pipeline")
				}

				features = append(features, &FeatureRequest{
					SourceVariableName:  denormFieldName,
					FeatureVariableName: indexName,
					Variable:            v,
					Step:                step,
				})
			}
		}
	}

	return features, nil
}

func getD3MIndexField(dr *metadata.DataResource) int {
	d3mIndexField := -1
	for _, v := range dr.Variables {
		if v.Name == metadata.D3MIndexName {
			d3mIndexField = v.Index
		}
	}

	return d3mIndexField
}

func toStringArray(in []interface{}) []string {
	strArr := make([]string, 0)
	for _, v := range in {
		strArr = append(strArr, v.(string))
	}
	return strArr
}

func toFloat64Array(in []interface{}) ([]float64, error) {
	strArr := make([]float64, 0)
	for _, v := range in {
		strFloat, err := strconv.ParseFloat(v.(string), 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert interface array to float array")
		}
		strArr = append(strArr, strFloat)
	}
	return strArr, nil
}

func getDataResource(meta *metadata.Metadata, resID string) *metadata.DataResource {
	// main data resource has d3m index variable
	for _, dr := range meta.DataResources {
		if dr.ResID == resID {
			return dr
		}
	}

	return nil
}

func getRelativePath(rootPath string, filePath string) string {
	relativePath := strings.TrimPrefix(filePath, rootPath)
	relativePath = strings.TrimPrefix(relativePath, "/")

	return relativePath
}

func copyResourceFiles(sourceFolder string, destinationFolder string) error {
	// if source contains destination, then go folder by folder to avoid
	// recursion problem

	if strings.HasPrefix(destinationFolder, sourceFolder) {
		// copy every subfolder that isn't the destination folder
		files, err := ioutil.ReadDir(sourceFolder)
		if err != nil {
			return errors.Wrapf(err, "unable to read source data '%s'", sourceFolder)
		}
		for _, f := range files {
			name := path.Join(sourceFolder, f.Name())
			if name != destinationFolder {
				err = copyResourceFiles(name, destinationFolder)
				if err != nil {
					return err
				}
			}
		}
	} else {
		err := copy.Copy(sourceFolder, destinationFolder)
		if err != nil {
			return errors.Wrap(err, "unable to copy source data")
		}
	}

	return nil
}
