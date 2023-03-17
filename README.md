# Image Resizing Service

This is a simple HTTP server that resizes images based on a scale factor parameter and returns the resulting image as a response.

## Getting Started

To start the service, run the following command:
`go run cmd/app/main.go`


This will start the server on port 8080.

## Endpoints

The server provides the following endpoints:

### `POST /images`
This endpoint allows you to upload a JPEG image to the server. The image will be saved in the ./images directory with a randomly generated ID. The endpoint expects the image data to be sent in the request body as multipart/form-data.

### `GET /image/{id}`

Returns the image with the specified `id`. The `id` parameter should be a string that represents the name of the image file without the file extension.

### Query Parameters

The following query parameters are supported:

- `quality` (optional): A string that represents the scale factor for the resized image. Valid values are "25", "50", "75", and "100". If not specified, the scale factor is 1.0 (i.e., no resizing is performed).

### Example

To request an image with id `example` that is resized to 50% quality, you can make a GET request to the following URL: `http://localhost:8080/image/example?quality=50`


## Testing

To test the service, you can use a tool such as [Postman](https://www.postman.com/) to send HTTP requests to the server.

First, start the server as described above.

Then, create a new request in Postman and set the URL to one of the endpoints described above. Add any necessary query parameters.

Finally, send the request and inspect the response.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.


