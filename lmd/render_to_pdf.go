package lmd

import (
    "bytes"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "mime/multipart"
    "net/http"
    "os"
)

// curl -F "image"=@"IMAGEFILE" -F "key"="KEY" URL
// curl -v -F "data=@00.load_write_no_action.md" -o test.pdf https://lmdpdfgen.herokuapp.com/
//
// Thanks for the help: @attila-o && @burfl from Stack Overflow, per ...
// http://stackoverflow.com/questions/20205796/golang-post-data-using-the-content-type-multipart-form-data
// and
// http://stackoverflow.com/questions/16311232/how-to-pipe-http-response-to-a-file-in-golang
func writeToPdf(contents string, outputFile string) {

    tempFile, err := ioutil.TempFile(os.TempDir(), "lmd-test-")
    if err != nil {
        log.Fatal(err)
    }
    defer os.Remove(tempFile.Name())
    writeAFile(tempFile.Name(), contents)

    var b bytes.Buffer
    w := multipart.NewWriter(&b)
    f, err := os.Open(tempFile.Name())
    if err != nil {
        log.Fatal(err)
    }

    fw, err := w.CreateFormFile("data", tempFile.Name())
    if err != nil {
        log.Fatal(err)
    }
    if _, err = io.Copy(fw, f); err != nil {
        log.Fatal(err)
    }

    // Don't forget to close the multipart writer.
    // If you don't close it, your request will be missing the terminating boundary.
    w.Close()

    // Now that you have a form, you can submit it to your handler.
    req, err := http.NewRequest("POST", "https://lmdpdfgen.herokuapp.com/", &b)
    if err != nil {
        log.Fatal(err)
    }
    // Don't forget to set the content type, this will contain the boundary.
    req.Header.Set("Content-Type", w.FormDataContentType())

    // Submit the request
    client := &http.Client{}
    res, err := client.Do(req)

    // Check the response
    if res.StatusCode != http.StatusOK {
        err = fmt.Errorf("bad status: %s", res.Status)
    }
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()

    out, err := os.Create(outputFile)
    if err != nil {
        log.Fatal(err)
    }
    defer out.Close()
    io.Copy(out, res.Body)
}
