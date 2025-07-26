## Description

The project parses videos based on the [SURF algorithm](https://en.wikipedia.org/wiki/Speeded_up_robust_features), extracting images that match some reference image provided.

## Run

 `docker run -p 7001:7001 ghcr.io/vkhobor/go-opencv-extractor:latest`

 Note: the image at first pull is quite large yet, because of debian base.

## Tests

Before running tests add the example folder below to `./samples` folder.

Then run:

`go test ./...`

The tests are minimal given the scope of the project


## Examples

You can find and example video with reference pictures here: https://drive.google.com/drive/folders/1D8G2S-EWgcTO-FMfbxrASbkY3pg-4oNR?usp=sharing

## Screenshot

<img width="2470" height="1804" alt="image" src="https://github.com/user-attachments/assets/968731ec-8fe5-4f2d-bf97-985029cc2f3f" />
