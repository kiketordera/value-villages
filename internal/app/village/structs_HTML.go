package village

import (
	"gopkg.in/mgo.v2/bson"
)

/*************************  VILLAGES & SERVICES *************************/

// SynchronizationHTML are the different synchronizations made to the villages
type SynchronizationHTML struct {
	ID                  bson.ObjectId `json:"id"`
	Date                string        `json:"date"`
	VillageEmitterName  string        `json:"villageemitername"`
	VillageReceiverName string        `json:"villagereceiverame"`
	IDVillageEmitter    bson.ObjectId `json:"idvillageemitter"`
	IDVillageReceiver   bson.ObjectId `json:"idvillagereceiver"`
	Admin               string        `json:"admin"`
	Path                string        `json:"path"`
	IsDone              bool          `json:"isdone"`
}

// ServiceHTML is the HTML representation of the values of the Service to be Human readable
type ServiceHTML struct {
	ID         bson.ObjectId `json:"id"`
	Name       string        `json:"name"`
	Village    string        `json:"village"`
	VillageID  string        `json:"villageid"`
	Balance    float64       `json:"balance"`
	Type       string        `json:"type"`
	TypeID     bson.ObjectId `json:"typeid"`
	Photo      string        `json:"photo"`
	DeliveryID bson.ObjectId `json:"deliveryid"`
}

// ServiceTypeHTML is the HTML representation of the values of the ServiceType to be Human readable
type ServiceTypeHTML struct {
	ID                bson.ObjectId `json:"id"`
	Name              string        `json:"name"`
	Icon              string        `json:"icon"`
	AllowNoQRPurchase bool          `json:"allownoqrpurchase"`
	Link              string        `json:"link"`
}

// AccessHTML is the HTML representation of the Accesses for users to be Human readable
type AccessHTML struct {
	ID          bson.ObjectId `json:"id"`
	NameUser    string        `json:"nameuser"`
	Username    string        `json:"username"`
	PhotoUser   string        `json:"photouser"`
	NameVillage string        `json:"namevillage"`
	Services    []ServiceHTML `json:"services"`
}

// SelectServiceHTML is all the services of the system ordered by village
type SelectServiceHTML struct {
	VillageName string        `json:"villagename"`
	Services    []ServiceHTML `json:"services"`
}

// ItemsPicker is to choose in a TouchPicker the item separated and ordered by category
type ItemsPicker struct {
	CategoryName string     `json:"categoryname"`
	Items        []ItemHTML `json:"stocks"`
}

// StockPicker is to choose in a TouchPicker the item separated and ordered by category
type StockPicker struct {
	CategoryName string      `json:"categoryname"`
	Stocks       []StockHTML `json:"stocks"`
}

// ItemHTML is Item with Human readable and the Photo
type ItemHTML struct {
	ID           bson.ObjectId `json:"id"`
	Name         string        `json:"name"`
	CategoryName string        `json:"categoryname"`
	Photo        string        `json:"photo"`
	Stock        []StockHTML   `json:"stock"`
	UnitType     UnitType      `json:"unittype"`
	IsTrackable  bool          `json:"istrackable"`
	Quantity     float64       `json:"quantity"`
}

/*************************  USERS  *************************/

// UserHTML is the HTML representation of the values of the User to be Human readable
type UserHTML struct {
	ID        bson.ObjectId `json:"id"`
	Name      string        `json:"name"`
	Surname   string        `json:"surname"`
	Age       int           `json:"age"`
	Gender    string        `json:"gender"`
	Photo     string        `json:"photo"`
	Tribe     string        `json:"tribe"`
	Role      string        `json:"role"`
	Story     string        `json:"story"`
	Username  string        `json:"username"`
	IDVillage bson.ObjectId `json:"villageid"`
	Balance   float64       `json:"balance"`
	Village   string        `json:"village"`
}

/*************************  DATABASE  *************************/

// AuditDayHTML is all the audits divided by day to show
type AuditDayHTML struct {
	Date     string      `json:"date"`
	Created  []AuditHTML `json:"created"`
	Modified []AuditHTML `json:"modified"`
	Deleted  []AuditHTML `json:"deleted"`
}

// AuditHTML is the Audit with the Humand readable information
type AuditHTML struct {
	Link  string `json:"link"`
	Icon  string `json:"icon"`
	Class string `json:"class"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	// This would be the Quantity in transaction and payments, Item in assignments and product Name in WorkerOrder
	SecondLine string `json:"secondline"`
}

/*************************  CATEGORIES  *************************/

// CategoryHTML is the Category with the Humand readable information
type CategoryHTML struct {
	ID           bson.ObjectId `json:"id"`
	Name         string        `json:"name"`
	Photo        string        `json:"icon"`
	Type         string        `json:"type"`
	TimeChecking string        `json:"timecheking"`
	Services     []string      `json:"services"`
	TypeItem     string        `json:"typeitem"`
}

/*************************  SYNC  *************************/

// DeliveryHTML is the DeliveryItem with the Humand readable information
type DeliveryHTML struct {
	ID                  bson.ObjectId `json:"id"`
	ServiceEmitterID    bson.ObjectId `json:"serviceemitterid"`
	ServiceEmitterName  string        `json:"serviceemittername"`
	ServiceReceiverID   bson.ObjectId `json:"servicereceiverid"`
	ServiceReceiverName string        `json:"servicereceivername"`
	Date                string        `json:"date"`
	// The destiny
	VillageName string `json:"villagename"`
	// This parameter will be filled on the way back of the sync, in every local village
	ManagerName string `json:"managername"`
	IsSent      bool   `json:"issent"`
	IsComplete  bool   `json:"iscomplete"`
}

/*************************  WORKSHOP  *************************/

// PaymentHTML is the Payment with the Humand readable information
type PaymentHTML struct {
	ID       bson.ObjectId `json:"id"`
	Worker   string        `json:"worker"`
	Date     string        `json:"date"`
	Quantity float64       `json:"quantity"`
	DateUnix int64         `json:"dateunix"`
}

// WorkerOrderHTML is the Worker Order with the Humand readable information
type WorkerOrderHTML struct {
	ID          bson.ObjectId `json:"id"`
	Photo       string        `json:"photoitem"`
	Date        string        `json:"date"`
	Quantity    int           `json:"quantity"`
	AlreadyMade int           `json:"alreadymade"`
	Status      string        `json:"status"`
	WorkerPhoto string        `json:"workerphoto"`
	ItemID      bson.ObjectId `json:"itemid"`
	ItemName    string        `json:"itemname"`
	Name        string        `json:"name"`
	DateUnix    int64         `json:"dateunix"`
}

// SaleHTML is the Sale with the Humand readable information
type SaleHTML struct {
	ID           bson.ObjectId `json:"id"`
	ServiceName  string        `json:"servicename"`
	ItemName     string        `json:"itemname"`
	ItemPhoto    string        `json:"itemphoto"`
	CategoryName string        `json:"categoryname"`
	CategoryType string        `json:"categorytype"`
	UserPhoto    string        `json:"userphoto"`
	UserName     string        `json:"username"`
	Date         string        `json:"date"`
	Price        float64       `json:"price"`
	Quantity     float64       `json:"quantity"`
	DateUnix     int64         `json:"dateunix"`
}

// VideoCourseHTML is the Video with the Humand readable information
type VideoCourseHTML struct {
	ID              bson.ObjectId `json:"id"`
	VideoCourseName string        `json:"name"`
	Description     string        `json:"description"`
	Price           float64       `json:"price"`
	Photo           string        `json:"photos"`
	Steps           int           `json:"numbersteps"`
	StepsID         []string      `json:"stepsid"`
	// For different implementations of this structure in the views
	VideoCourseID bson.ObjectId `json:"videoid"`
	LinkPhoto     string        `json:"linkphoto"`
}

// ServiceOrderHTML is the ServiceOrder with the Humand readable information
type ServiceOrderHTML struct {
	ID           bson.ObjectId `json:"id"`
	ProductName  string        `json:"productname"`
	Photo        string        `json:"productphoto"`
	Village      string        `json:"village"`
	WorkshopName string        `json:"workshopname"`
	Deadline     string        `json:"deadline"`
	Window       string        `json:"window"`
	Quantity     int           `json:"quantity"`
	AlreadyMade  int           `json:"alreadyMade"`
	Assigned     int           `json:"assigned"`
	Status       string        `json:"status"`
}

// Wallet is the movements in the account of a Wallet
type Wallet struct {
	ID          string  `json:"id"`
	ItemName    string  `json:"itemname"`
	ItemPhoto   string  `json:"itemphoto"`
	Date        string  `json:"date"`
	Price       float64 `json:"price"`
	Balance     float64 `json:"balance"`
	ServiceName string  `json:"servicename"`
	DateUnix    int64   `json:"dateunix"`
}

/*************************  COMUNICATION *************************/

// MessageHTML are the messages send in a conversation for human readable
type MessageHTML struct {
	ID             bson.ObjectId `json:"id"`
	User           string        `json:"username"`
	IDConversation bson.ObjectId `json:"idconversation"`
	Text           string        `json:"text"`
	Photo          string        `json:"photo"`
	Audio          string        `json:"audio"`
	Class          string        `json:"class"`
	Date           string        `json:"date"`
	Time           string        `json:"time"`
}

// MessageDayHTML is all the audits divided by day to show
type MessageDayHTML struct {
	Date     string        `json:"date"`
	Messages []MessageHTML `json:"messages"`
}

// ConversationHTML is the conversation between two or more persons in the app for showing in the table
type ConversationHTML struct {
	ID          bson.ObjectId `json:"id"`
	LastMessage string        `json:"lastmessage"`
	Users       []UserHTML    `json:"users"`
}

// ReportHTML is the report of the users with Human readable information
type ReportHTML struct {
	ID      bson.ObjectId `json:"id"`
	User    string        `json:"user"`
	Text    string        `json:"text"`
	Photo   string        `json:"photo"`
	Audio   string        `json:"audio"`
	Date    string        `json:"date"`
	Type    string        `json:"type"`
	IsClose bool          `json:"isclose"`
	Village string        `json:"village"`
}

// StockHTML is the stock of Items and products that the services has
type StockHTML struct {
	ID                    bson.ObjectId `json:"id"`
	VideoOrItemName       string        `json:"videooritemname"`
	VideoOrItemID         bson.ObjectId `json:"videooritemid"`
	ServiceName           string        `json:"serviceid"`
	Price                 float64       `json:"price"`
	Photo                 string        `json:"photo"`
	UserPhoto             string        `json:"userphoto"`
	ServiceIDLocation     bson.ObjectId `json:"servicelocation"`
	ServiceIDLocationName string        `json:"servicelocationname"`
	ItemOrVideoCourseID   bson.ObjectId `json:"itemid"`
	// Only if applicable
	Quantity  float64       `json:"quantity"`
	UnitType  UnitType      `json:"unittype"`
	UserID    bson.ObjectId `json:"userid"`
	ManagerID bson.ObjectId `json:"managerid"`
	Date      string        `json:"date"`
	DateUnix  int64         `json:"dateunix"`
	// This are the users that comes into the whole process of the product
	UsersRelated []bson.ObjectId `json:"usersrelated"`
	// ServiceCreated is where it was produced
	ServiceCreated bson.ObjectId `json:"servicecreated"`
	IsTrackable    bool          `json:"istrackable"`
	CategoryName   string        `json:"categoryname"`
}

/***********************  TOOLS ASSIGNMENTS ***********************/

// AssignmentHTML is tools assigned to the Workers with human readable information
type AssignmentHTML struct {
	ID          bson.ObjectId `json:"id"`
	Worker      string        `json:"worker"`
	WorkerPhoto string        `json:"workerphoto"`
	Service     string        `json:"serviceid"`
	Manager     string        `json:"manager"`
	Item        string        `json:"item"`
	PhotoItem   string        `json:"photoItem"`
	Date        string        `json:"date"`
	// This will be true when the user returns the Tool to the Manager
	IsBack bool `json:"isback"`
}

// ToDoHTML is the item to check
type ToDoHTML struct {
	ID              bson.ObjectId      `json:"id"`
	TitleToDo       string             `json:"titletodo"`
	DescriptionToDo string             `json:"description"`
	IDvillage       bson.ObjectId      `json:"idvillage"`
	NameVillage     string             `json:"namevillage"`
	IDuser          bson.ObjectId      `json:"iduser"`
	IsTrackable     bool               `json:"istrackable"`
	NameUser        string             `json:"nameuser"`
	TimeChecking    TimecheckingType   `json:"timecheking"`
	Descriptions    []FieldDescription `json:"descriptions"`
	Checkboxes      []FieldCheckbox    `json:"checkboxes"`
	Photos          []FieldPhoto       `json:"photos"`
	Numbers         []FieldNumber      `json:"numbers"`
	ToDoChecks      []ToDoChecked      `json:"checks"`
}

/***********************  PDF ***********************/

// PDFPageqrHTML is on of the pages of the PDF
type PDFPageqrHTML struct {
	ID         string   `json:"id"`
	Date       string   `json:"date"`
	TypeOfItem string   `json:"typeofitem"`
	QRs        []string `json:"qrs"`
}

// PDFHTML is the PDF generated and listed for the table
type PDFHTML struct {
	ID         bson.ObjectId `json:"id"`
	Date       string        `json:"date"`
	TypeOfItem string        `json:"typeofitem"`
	Pages      int           `json:"pages"`
	TypeOfPDF  string        `json:"typeofpdf"`
}

// Pagination is the abstract form to be able to put the correct number in the table pagination for the human view
type Pagination struct {
	Active   bool `json:"active"`
	Fakenumb int  `json:"fakenum"`
}
