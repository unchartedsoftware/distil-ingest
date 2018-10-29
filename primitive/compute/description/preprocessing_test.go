package description

import (
	"io/ioutil"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/unchartedsoftware/distil-ingest/metadata"
)

func TestCreateUserDatasetPipeline(t *testing.T) {

	pipeline, err := CreateUserDatasetPipeline(
		"test_user_pipeline", "a test user pipeline", "test_target")

	// assert 1st is a semantic type update
	hyperParams := pipeline.GetSteps()[0].GetPrimitive().GetHyperparams()
	assert.Equal(t, []int64{1, 3}, ConvertToIntArray(hyperParams["add_columns"].GetValue().GetData().GetRaw().GetList()))
	assert.Equal(t, []string{"http://schema.org/Integer"}, ConvertToStringArray(hyperParams["add_types"].GetValue().GetData().GetRaw().GetList()))
	assert.Equal(t, []int64{}, ConvertToIntArray(hyperParams["remove_columns"].GetValue().GetData().GetRaw().GetList()))
	assert.Equal(t, []string{""}, ConvertToStringArray(hyperParams["remove_types"].GetValue().GetData().GetRaw().GetList()))

	// assert 2nd is a semantic type update
	hyperParams = pipeline.GetSteps()[1].GetPrimitive().GetHyperparams()
	assert.Equal(t, []int64{}, ConvertToIntArray(hyperParams["add_columns"].GetValue().GetData().GetRaw().GetList()))
	assert.Equal(t, []string{""}, ConvertToStringArray(hyperParams["add_types"].GetValue().GetData().GetRaw().GetList()))
	assert.Equal(t, []int64{1, 3}, ConvertToIntArray(hyperParams["remove_columns"].GetValue().GetData().GetRaw().GetList()))
	assert.Equal(t, []string{"https://metadata.datadrivendiscovery.org/types/CategoricalData"},
		ConvertToStringArray(hyperParams["remove_types"].GetValue().GetData().GetRaw().GetList()))

	// assert 3rd step is column remove and index two was remove
	hyperParams = pipeline.GetSteps()[2].GetPrimitive().GetHyperparams()
	assert.Equal(t, "0", hyperParams["resource_id"].GetValue().GetData().GetRaw().GetString_(), "0")
	assert.Equal(t, []int64{2}, ConvertToIntArray(hyperParams["columns"].GetValue().GetData().GetRaw().GetList()))

	assert.NoError(t, err)
	t.Logf("\n%s", proto.MarshalTextString(pipeline))
}

func TestCreateUserDatasetPipelineMappingError(t *testing.T) {

	pipeline, err := CreateUserDatasetPipeline(
		"test_user_pipeline", "a test user pipeline", "test_target")
	assert.Error(t, err)
	t.Logf("\n%s", proto.MarshalTextString(pipeline))
}

func TestCreateUserDatasetEmpty(t *testing.T) {

	pipeline, err := CreateUserDatasetPipeline(
		"test_user_pipeline", "a test user pipeline", "test_target")

	assert.Nil(t, pipeline)
	assert.Nil(t, err)

	t.Logf("\n%s", proto.MarshalTextString(pipeline))
}

func TestCreatePCAFeaturesPipeline(t *testing.T) {
	pipeline, err := CreatePCAFeaturesPipeline("pca_features_test", "test pca feature ranking pipeline")
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/create_pca_features.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateSimonPipeline(t *testing.T) {
	pipeline, err := CreateSimonPipeline("simon_test", "test simon classification pipeline")
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/create_simon.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateCrocPipeline(t *testing.T) {
	pipeline, err := CreateCrocPipeline("croc_test", "test croc object detection pipeline", []string{"filename"}, []string{"objects"})
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/create_croc.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateUnicornPipeline(t *testing.T) {
	pipeline, err := CreateUnicornPipeline("unicorn test", "test unicorn image detection pipeline", []string{"filename"}, []string{"objects"})
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/create_unicorn.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateSlothPipeline(t *testing.T) {
	timeSeriesVariables := []*metadata.Variable{
		{Name: "time", Index: 0},
		{Name: "value", Index: 1},
	}

	pipeline, err := CreateSlothPipeline("sloth_test", "test sloth object detection pipeline", "time", "value", timeSeriesVariables)
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/create_sloth.pln", data, 0644)
	assert.NoError(t, err)
}

func TestCreateDukePipeline(t *testing.T) {
	pipeline, err := CreateDukePipeline("duke_test", "test duke data summary pipeline")
	assert.NoError(t, err)

	data, err := proto.Marshal(pipeline)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	err = ioutil.WriteFile("/tmp/create_duke.pln", data, 0644)
	assert.NoError(t, err)
}