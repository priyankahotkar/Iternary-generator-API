# Vigovia Itinerary PDF Generator

This project generates beautiful, branded travel itinerary PDFs from JSON requests using Go, Gin, and gofpdf.

## Features
- Modular, visually appealing PDF output
- Custom branding and logo
- Day-wise itinerary, trip summary, flights, hotels, notes, payment, and more

---

## Prerequisites
- [Go](https://golang.org/dl/) 1.18 or newer
- [Git](https://git-scm.com/)

## Clone the Repository
```sh
git clone https://github.com/priyankahotkar/Iternary-generator-API.git
```

## Install Dependencies
```sh
go mod tidy
```

## Project Structure
- `main.go` — Gin server and API endpoint
- `pdf_generator.go` — PDF generation logic
- `public/logo.png` — Your company logo (used in the PDF header)
- `pdfs/` — Generated PDFs are saved here

## Add Your Logo
Place your logo at:
```
public/logo.png
```
(Recommended size: ~200x100px, PNG format)

## Run the Server Locally
```sh
go run .
```
The server will start on [http://localhost:8080](http://localhost:8080)

## Generate an Itinerary PDF
You can use [Postman](https://www.postman.com/) to easily test the API:

1. **Open Postman** (or any similar API client).
2. **Set the request method** to `POST` and the URL to:
   ```
   http://localhost:8080/generate-itinerary
   ```
3. **Go to the "Body" tab**, select `raw`, and choose `JSON` from the dropdown.
4. **Paste your sample JSON request** (see the chat or your sample_request.json file).
5. **Go to the "Headers" tab** and ensure you have:
   - Key: `Content-Type`  Value: `application/json`
6. **Click "Send"**.
7. The response will include a URL to download or view the generated PDF.

## View Generated PDFs
All PDFs are saved in the `pdfs/` directory. You can also access them via:
```
http://localhost:8080/pdfs/itinerary_<name>_<startdate>.pdf
```

## Troubleshooting
- **Logo not showing?** Make sure `public/logo.png` exists and is a valid PNG.
- **Port in use?** Change the port in `main.go` if needed.
- **PDF not generating?** Check the server logs for errors.
- **Unicode/emoji issues?** gofpdf may have limited emoji support; use PNG icons for full compatibility.

## License
MIT (or your preferred license) 