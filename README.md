# X

[![Join the chat at https://gitter.im/tapeja/x](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/tapeja/x?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

**WARNING: This is a work in progress!**

X is an image processing service written in Go. It is designed to be easy to
configure, deploy, and scale.

It compiles into a small static binary which can run on your local machine
during development or across thousands of containers in the cloud.

## Configure

To get started you must configure some endpoints for your image processing
service via a yaml file, for example a `config.yml` could contain :

```yaml
---

thumbs:
  sizes:
    - square: 80
    - square: 200
  formats:
    - webp
    - png
    - jpg
  stores:
    - s3: mythumbs
      aws_access_key: randomstuff
      aws_secret_key: randomstuffs

gallery:
# ... and on and on, you can define as many endpoints as you wish.

```

Each endpoint configuration starts with a hash key with the endpoint name, i.e.
`thumbs`. The name can be any string composed of `a-z`, `A-Z`, `0-9`, `-`,  and
`_`.  The endpoint holds a hash with the keys `sizes`, `formats`, and `stores`
which configure the processing and storage for the endpoint.

### `sizes`

List of hashes, defining the different size constraints on the processed image.
A different image will be created for each constraint. All contstraints maintain
aspect ratios and avoid enlarging the images. The parameters are defined in
pixels.

Valid constraint values are :

- `square`

	Creates a square image with the dimensions of the passed parameter.
First crops the larger dimension equally from both sides to match the smaller
dimension then shrinks the square down to the passed parameter. Image isn't
resized if both dimensions are smaller than the passed parameter.

- `max`

	Determines the larger dimension of the image and shrinks it to match the
passed parameter retaining the original aspect ratio. Image isn't resized if
both dimensions are smaller than the passed parameter.

- `max_width` and `max_height`

	Just like `max` but instead of determining the largest dimension it is
fixed on either the width or height. Image isn't resized if the fixed dimension
is smaller than the passed parameter.

Constraint examples, rows are the modes, the columns are the original dimensions
expressed in pixels with width x height, and the cells are the output dimensions
:

|                  | 80x160 | 160x80 |   20x80   |   20x40   |
|------------------|:------:|:------:|:---------:|:---------:|
| `square: 40`     |  40x40 |  40x40 |   20x40   | no resize |
| `max: 40`        |  20x40 |  40x20 |   10x40   | no resize |
| `max_height: 40` |  20x40 |  80x40 |   10x40   | no resize |
| `max_width: 40`  |  40x80 |  40x20 | no resize | no resize |

### `formats`

List of formats to output, valid values are :

- `jpg`
- `png`
- `webp`

### `stores`

List of hashes, configuring the different places you would like to store the
generated images. Currently available stores are local file storage and S3.

- `file`

	Stores files in the specified local directory.

- `s3`

	Stores files on [S3](http://aws.amazon.com/s3/) in the given bucket.
Must also provide the appropriate AWS credentials with `aws_access_key` and
`aws_secret_key`.

Example stores configuration :

```yaml
stores:

  # Configure an S3 bucket store
  - s3: mys3bucket
    aws_access_key: randomstuff
    aws_secret_key: randomstuffs

  # Configure a local file store
  - file: /mnt/mydirectory
```

## Process

To begin processing images, start up the server by running the binary :

```
./x -c config.yml -p 8080
```

The image processing service will now be running on port 8080. You can send one
or more images for processing via a `POST` request to the appropriate endpoint.
The endpoints are prefixed with `v1`, the current version of the API.

Example request :

```
curl -F "first_image=@avatar.jpeg" -F "second_image=@cat.jpeg" localhost:8080/v1/thumbs
```

A response will come back with a JSON hash associating the `name` attribute of
each file to the SHA of the provided file or an error code. If one of the images
in a request produces an error the whole response will fail with the appropriate
status code.

```json
{
  "first_image": "9fefe3e2d4c5e849f3a8f0136db973c2c012afcd",
  "second_image": "7638417db6d59f3c431d3e1f261cc637155684cd",
}
```

The images will then be available in each configured datastore at the following
path(s) :

```
{{ size key }}/{{ size parameter }}/{{ sha of original }}.{{ format }}
```

For example the `thumbs` endpoint configured in the example above would
output files at the following locations :

```
square/80/9fefe3e2d4c5e849f3a8f0136db973c2c012afcd.jpg
square/200/9fefe3e2d4c5e849f3a8f0136db973c2c012afcd.jpg
square/80/9fefe3e2d4c5e849f3a8f0136db973c2c012afcd.png
square/200/9fefe3e2d4c5e849f3a8f0136db973c2c012afcd.png
square/80/9fefe3e2d4c5e849f3a8f0136db973c2c012afcd.webp
square/200/9fefe3e2d4c5e849f3a8f0136db973c2c012afcd.webp
```

## Build

For a truly static binary you need to use Go with a version greater than 1.2 and
build with the [`netgo`](https://golang.org/doc/go1.2) tag and
[`cgo`](http://golang.org/cmd/cgo/) disabled :

```bash
CGOENABLED=0 GOOS=linux go build -a -o x --tags netgo .
```
