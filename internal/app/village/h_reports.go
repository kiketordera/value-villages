package village

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/signintech/gopdf"
	qrcode "github.com/skip2/go-qrcode"
	"gopkg.in/mgo.v2/bson"
)

// Shows the screen with all QR codes
func (VS *Server) showQRCodes(c *gin.Context) {
	PDFQR := VS.getAllPDFQR(c, "h_reports", "showQRCodes", "PDFQR := VS.getAllPDFQR")
	PDFQR = revertPDF(PDFQR)
	PDFHTML := fromPDFstoPDFHTMLs(PDFQR)
	JSON, err := json.Marshal(PDFHTML)
	VS.checkOperation(c, "h_date", "showQRCodes", "JSON, err := json.Marshal(uHTML)", err)
	render(c, gin.H{
		"PDFQR": PDFHTML,
		"JSON":  string(JSON),
	}, "qr-codes.html")
}

// Shows the screen with all QR codes
func (VS *Server) showNewQRCodesForm(c *gin.Context) {
	render(c, gin.H{}, "form-qr.html")
}

// The POST method that created the new PDF and goes back to the list of QR PDFs
func (VS *Server) generatePDFNewQRCodes(c *gin.Context) {
	idPDF := bson.NewObjectId()
	pagesNumber, b := VS.getIntFromHTML(c, "pages", "PDF QR")
	if !b {
		return
	}
	typeOfItem, b := VS.getStringFromHTML(c, "typeofitem", "PDF QR")
	if !b {
		return
	}
	VS.generatePDFfromURL(c, URLbase+"reports/generate-new-pdf/"+idPDF.Hex()+"/"+typeOfItem+"/"+strconv.Itoa(pagesNumber), "/local-resources/pdf/qr", idPDF)
	c.Redirect(http.StatusFound, "/reports/show-qr-codes/")
}

// Shows the PDF in the browser
func (VS *Server) showPDF(c *gin.Context) {
	pathToPDF := basePath + "/local-resources/pdf/qr/"
	fileName, b := VS.getParamFromURL(c, "name")
	if !b {
		return
	}
	targetPath := filepath.Join(pathToPDF, fileName)
	// This check is for example, I not sure is it can prevent all possible filename attacks - will be much better if real filename will not come from user side. I not even tryed this code
	if !strings.HasPrefix(filepath.Clean(targetPath), pathToPDF) {
		c.String(403, "Look like you attacking me")
		return
	}
	// Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	// Add this two headers to download instead of see in browser
	// c.Header("Content-Disposition", "attachment; filename="+fileName )
	// c.Header("Content-Type", "application/octet-stream")
	c.File(targetPath)
}

//******* OLD METHOD, this is related with wkhtmltopdf but gives problems with docker, even tho the methods works in macOS perfectly with wkhtmltopdf installed  ************//
// PDFqr shows the HTML that shows the QR codes whenever is called and saves the ID in the DB.
// This is transparent and should be only called for converting it to PDF
func (VS *Server) newPDFqrOld(c *gin.Context) {
	// Get parameters to make the PDF
	id, b := VS.getIDParamFromURL(c, "id")
	if !b {
		return
	}
	typeOfI, b := VS.getParamFromURL(c, "type")
	typeOfItem := TypeOfItem(typeOfI)
	if !b {
		return
	}
	pagesN, b := VS.getParamFromURL(c, "pages")
	pagesNumber, _ := strconv.Atoi(pagesN)
	if !b {
		return
	}
	// Change the color of the QR depending on the type of Item
	var color color.RGBA
	switch typeOfItem {
	case ServiceProduct:
		color = Pink
	case Tool:
		color = DarkBrown
	case PrimaryMaterial:
		color = Blue
	}
	var PDF PDF
	PDF.Date = time.Now().Unix()
	PDF.ID = id
	PDF.Pages = pagesNumber
	PDF.TypeOfItem = typeOfItem
	PDF.TypeOfPDF = PDFQR
	VS.addElementToDatabaseWithoutRegister(c, PDF, PDF.ID, "h_reports", "newPDFqr", "PDF.TypeOfPDF = PDFQR ")
	// We create every page of the PDF
	var QRs []QR
	var QR QR
	var imagesToHTML []string
	var PDFpages []PDFPageqrHTML
	for page := 0; page < pagesNumber; page++ {
		for i := 0; i < 108; i++ {
			// Create and add the ID
			QR.ID = bson.NewObjectId()
			QR.PDF = PDF.ID
			QR.PageNumber = page + 1
			QR.TypeOfItem = typeOfItem
			QRs = append(QRs, QR)
			VS.addElementToDatabaseWithoutRegister(c, QR, QR.ID, "h_reports", "newPDFqr", "QRs = append(QRs, QR)")
			// Create and add the QR image
			qr, err := qrcode.New(QRs[i].ID.Hex(), qrcode.High)
			qr.ForegroundColor = color
			img, err := qr.PNG(256)
			if err != nil {
				fmt.Print(err)
			}
			// Create the images for the HTML
			img2 := base64.StdEncoding.EncodeToString(img)
			imagesToHTML = append(imagesToHTML, img2)
		}
		PDFpages = append(PDFpages, getPagePDFqr(imagesToHTML, PDF))
		QRs = nil
		imagesToHTML = nil
	}
	// Here we show it in HTML to be able to make the pdf. After showing wee need to transform it in to PDF
	render(c, gin.H{
		"PDFpages": PDFpages,
	}, "new-qr.html")
}

// Creates the PDF for the QR codes
func (VS *Server) newPDFqr(c *gin.Context) {
	typeOfI, b := VS.getStringFromHTML(c, "typeofitem", "PDF")
	typeOfItem := TypeOfItem(typeOfI)
	if !b {
		return
	}
	pagesNumber, b := VS.getIntFromHTML(c, "pages", "PDF")
	if !b {
		return
	}
	// Change the color of the QR depending on the type of Item
	var color color.RGBA
	switch typeOfItem {
	case ServiceProduct:
		color = Pink
	case Tool:
		color = DarkBrown
	case PrimaryMaterial:
		color = Blue
	case Task:
		color = Green
	}
	var PDF PDF
	PDF.Date = time.Now().Unix()
	PDF.ID = bson.NewObjectId()
	PDF.Pages = pagesNumber
	PDF.TypeOfItem = typeOfItem
	PDF.TypeOfPDF = PDFQR
	VS.addElementToDatabaseWithoutRegister(c, PDF, PDF.ID, "h_reports", "newPDFqr", "PDF.TypeOfPDF = PDFQR ")
	// We create every page of the PDF
	var QRs []QR
	var QR QR
	var imagesToHTML []string
	var PDFpages []PDFPageqrHTML
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	for page := 0; page < pagesNumber; page++ {
		pdf.AddPage()
		x, y := 20.0, 10.0
		for i, k := 0, 0; i < 77; i++ {
			// Create and add the ID in the DataBase
			QR.ID = bson.NewObjectId()
			QR.PDF = PDF.ID
			QR.PageNumber = page + 1
			QR.TypeOfItem = typeOfItem
			QRs = append(QRs, QR)
			VS.addElementToDatabaseWithoutRegister(c, QR, QR.ID, "h_reports", "newPDFqr", "QRs = append(QRs, QR)")

			// Create QR images and add to PDF
			qr, err := qrcode.New(QR.ID.Hex(), qrcode.High)
			if err != nil {
				fmt.Print(err)
			}
			qr.ForegroundColor = color
			qr.BackgroundColor = White
			img, err := qr.PNG(140)
			if err != nil {
				fmt.Print(err)
			}
			// convert []byte to image for saving to file
			image, _, _ := image.Decode(bytes.NewReader(img))
			//save the imgByte to file
			os.MkdirAll(basePath+"/local-resources/pdf/test/", os.ModePerm)

			out, err := os.Create(basePath + "/local-resources/pdf/test/" + QR.ID.Hex() + ".png")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = png.Encode(out, image)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			pdf.Image(basePath+"/local-resources/pdf/test/"+QR.ID.Hex()+".png", x, y, nil) //print image
			x += 75
			k++
			if k >= 7 {
				y += 75
				x = 20
				k = 0
			}
		}
		PDFpages = append(PDFpages, getPagePDFqr(imagesToHTML, PDF))
		QRs = nil
		imagesToHTML = nil
	}

	os.MkdirAll(basePath+"/local-resources/pdf/qr/", os.ModePerm)
	err := pdf.WritePdf(basePath + "/local-resources/pdf/qr/" + PDF.ID.Hex() + ".pdf")
	if err != nil {
		os.MkdirAll(basePath+"/local-resources/pdf/qr/", os.ModePerm)
		pdf.WritePdf(basePath + "/local-resources/pdf/qr/" + PDF.ID.Hex() + ".pdf")
	}
	fmt.Print("PDF done")
	// Delete the images of the QR that we created before
	os.RemoveAll(basePath + "/local-resources/pdf/test/")
	VS.goodFeedback(c, "reports/show-qr-codes")
}
