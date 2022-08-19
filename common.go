package honeycomb

const (
	noApiKeyDetectedMessage         string = "Missing an API Key!\nConfigure via HONEYCOMB_API_KEY environment variable, or in code."
	classicKeyMissingDatasetMessage string = "Honeycomb Classic API Key detected!\nYour API key: %s requires a dataset to be configured.\nConfigure via HONEYCOMB_DATASET or in code."
	dontSetADatasetMessageMessage   string = "Dataset detected! Datasets are a Honeycomb Classic configuration value.\nUnset HONEYCOMB_DATASET or remove configuration code that sets a dataset."
)

func isclassicApiKey(apiKey string) bool {
	return len(apiKey) == 32
}
