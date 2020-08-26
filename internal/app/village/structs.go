package village

import (
	"github.com/timshannon/bolthold"
	"gopkg.in/mgo.v2/bson"
)

/*************************  VILLAGES & SERVICES *************************/

// Server is the server that instances to run the DataBase of the village
type Server struct {
	DataBase *bolthold.Store
}

// Village are the villages that we were runing different instances of the app
type Village struct {
	ID   bson.ObjectId `json:"id"`
	Name string        `json:"name"`
	// The prefix is a 4 CAPITAL letter as a shorter name and for putting this prefix in the username of the Workers, only workers
	Prefix string `json:"prefix"`
}

// Service are the different services that the Village has
type Service struct {
	ID        bson.ObjectId `json:"id"`
	Name      string        `json:"name"`
	IDVillage bson.ObjectId `json:"idvillage"`
	Balance   float64       `json:"balance"`
	Type      bson.ObjectId `json:"type"`
}

// ServiceType are the different type of services that we can implement in the villages
type ServiceType struct {
	ID                bson.ObjectId `json:"id"`
	Name              string        `json:"name"`
	Icon              string        `json:"icon"`
	AllowNoQRPurchase bool          `json:"allownoqrpurchase"`
}

// Access are the different accesses that the User have to the services
type Access struct {
	ID        bson.ObjectId `json:"id"`
	IDUser    bson.ObjectId `json:"iduser"`
	IDService bson.ObjectId `json:"idservice"`
	IsActive  bool          `json:"isactive"`
}

// Sale is the sale of Items (material) that the services make to the Users
type Sale struct {
	ID        bson.ObjectId `json:"id"`
	IDService bson.ObjectId `json:"idservice"`
	IDStock   bson.ObjectId `json:"idstock"`
	IDWorker  bson.ObjectId `json:"idworker"`
	IDManager bson.ObjectId `json:"idmanager"`
	Date      int64         `json:"date"`
	Price     float64       `json:"price"`
	Quantity  float64       `json:"quantity"`
}

/*************************  USERS  *************************/

// User is the user of the system
type User struct {
	ID        bson.ObjectId `json:"id"`
	Name      string        `json:"name"`
	Surname   string        `json:"surname"`
	Age       int           `json:"age"`
	Gender    GenderType    `json:"gender"`
	Photo     string        `json:"photo"`
	Tribe     TribeType     `json:"tribe"`
	Role      Role          `json:"role"`
	Story     string        `json:"story"`
	Password  string        `json:"password"`
	Username  string        `json:"username"`
	IDVillage bson.ObjectId `json:"villageid"`
	Balance   float64       `json:"balance"`
}

// Role is the close list of different roles that the users can have
type Role string

// The options for RoleType
const (
	Admin   Role = "admin"
	Manager Role = "manager"
	Worker  Role = "worker"
)

// GenderType is the close list of different gender that the users can have
type GenderType string

// The options for GenderType
const (
	Male   GenderType = "Male"
	Female GenderType = "Female"
)

// TribeType is the close list of different tribes that the users can have
type TribeType string

// The options for TribeType
const (
	Ngidoca         TribeType = "Ngidoca"
	Ngiduya         TribeType = "Ngiduya"
	Ngikadanya      TribeType = "Ngikadanya"
	Ngikalesso      TribeType = "Ngikalesso"
	Ngikatap        TribeType = "Ngikatap"
	Ngikateok       TribeType = "Ngikateok"
	Ngikinom        TribeType = "Ngikinom"
	Ngilelet        TribeType = "Ngilelet"
	Ngilobol        TribeType = "Ngilobol"
	Ngimedeo        TribeType = "Ngimedeo"
	Ngimerpur       TribeType = "Ngimerpur"
	Ngimeturuana    TribeType = "Ngimeturuana"
	Ngingolereto    TribeType = "Ngingolereto"
	Ngiponga        TribeType = "Ngiponga"
	Ngipuco         TribeType = "Ngipuco"
	Ngirarak        TribeType = "Ngirarak"
	Ngisalika       TribeType = "Ngisalika"
	Ngisiger        TribeType = "Ngisiger"
	Ngitarapakolong TribeType = "Ngitarapakolong"
	Ngitengor       TribeType = "Ngitengor"
	Ngiteso         TribeType = "Ngiteso"
	Other           TribeType = "Other"
)

/*************************  DATABASE  *************************/

// Synchronization are the different synchronizations made to the villages
// IsDone will be true when the Village receiver makes the sync
type Synchronization struct {
	ID                bson.ObjectId `json:"id"`
	Date              int64         `json:"date"`
	IDVillageEmitter  bson.ObjectId `json:"idvillageemitter"`
	IDVillageReceiver bson.ObjectId `json:"idvillagereceiver"`
	IDAdmin           bson.ObjectId `json:"idadmin"`
	IsDone            bool          `json:"isdone"`
	DateSyncReceiver  int64         `json:"datesyncreceiver"`
}

// SyncOptions is the close list of different options for Sync
type SyncOptions string

// The options for SyncOptions
const (
	exp SyncOptions = "export"
	imp SyncOptions = "import"
)

// Audit are the changes made in the DataBase
type Audit struct {
	ID                bson.ObjectId     `json:"id"`
	IDItem            bson.ObjectId     `json:"iditem"`
	IDService         bson.ObjectId     `json:"idservice"`
	IDUser            bson.ObjectId     `json:"iduser"`
	Description       AuditType         `json:"description"`
	Date              int64             `json:"date"`
	IDVillage         bson.ObjectId     `json:"idvillage"`
	InformationObject map[string]string `json:"informationobject"`
}

// AuditType is the close list of different types that an audit can have
type AuditType string

// The options for AuditType
const (
	Created  AuditType = "created"
	Modified AuditType = "modified"
	Deleted  AuditType = "deleted"
)

// Change is the change into one of the attributes of the struct, there can be several Changes per Audit, one per attribute modified
type Change struct {
	ID            bson.ObjectId `json:"id"`
	IDAudit       bson.ObjectId `json:"idaudit"`
	NameAttribute string        `json:"nameattribute"`
	Before        string        `json:"before"`
	After         string        `json:"after"`
}

// EventDB are the changes made in the DataBase as an Event driven, is the backup is case of Sync problem
type EventDB struct {
	ID        bson.ObjectId `json:"id"`
	IDAudit   bson.ObjectId `json:"idaudit"`
	IDVillage bson.ObjectId `json:"idvillage"`
	Previous  interface{}   `json:"previous"`
	After     interface{}   `json:"after"`
	Date      int64         `json:"date"`
}

/***********************  CATEGORY ***********************/

// Category are the categories where the items will be classified for the inventory
type Category struct {
	ID             bson.ObjectId    `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Photo          string           `json:"photo"`
	Type           bson.ObjectId    `json:"type"`
	TimeChecking   TimecheckingType `json:"timecheking"`
	Items          []bson.ObjectId  `json:"items"`
	IsTrackable    bool             `json:"istrackable"`
	ServicesAccess []bson.ObjectId  `json:"servicesaccess"`
	TypeOfItem     TypeOfItem       `json:"typeofitem"`
}

// TimecheckingType is the close list of different times for checking that a category can have
type TimecheckingType string

// The options for TimecheckingType
const (
	Never         TimecheckingType = "never"
	OnlyOnce      TimecheckingType = "once"
	EveryMonth    TimecheckingType = "monthly"
	EveryWeek     TimecheckingType = "weekly"
	EveryDay      TimecheckingType = "daily"
	ChangeVillage TimecheckingType = "changevillage"
)

// CategoryType are the different types that a category can have
type CategoryType struct {
	ID   bson.ObjectId `json:"id"`
	Name string        `json:"name"`
}

// Item are the different items that we can register for latter make individual elements out of the item.
type Item struct {
	ID          bson.ObjectId `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Price       float64       `json:"price"`
	// 100 shs per Kilogram = 0,1 shs per gram
	UnitType   UnitType      `json:"unit"`
	Photo      string        `json:"photo"`
	IDCategory bson.ObjectId `json:"idcategory"`
}

// UnitType is the close list of different units that an Item can have
type UnitType string

// The options for UnitType
const (
	Meter    UnitType = "meter"
	Liter    UnitType = "liter"
	Kilogram UnitType = "kilogram"
	// Here refers that the Item is a box with 5.000 screws (5.000 PhisicalItem)
	Unit UnitType = "unit"
)

// Stock are the individual physical items, that can be or not trackable, delivered, buy or purchased.
// Video Course needs to have a IsOnlyEducative to don't be able to create stock out of a OnlyEducativeVideo
// and the stocks from the videos will all be trackable
// Restriction here, if it is trackable, it can not be split in the delivery.
// Restriction here, if tha category of the item have units, it can not be trackable.
// Related with the last one, if cat has
type Stock struct {
	ID bson.ObjectId `json:"id"`
	// If Stock comes from a Video Tutorial
	IDVideoCourseOrItem bson.ObjectId `json:"idvideocourseoritem"`
	Photo               string        `json:"photo"`
	// ServiceIDLocation where it stays physically in the moment
	IDServiceLocation bson.ObjectId `json:"idservicelocation"`
	TypeOfItem        TypeOfItem    `json:"typeofitem"`
	// // Only if applicable
	// Quantity float64 `json:"quantity"`
	// The User who created the stock
	IDUser    bson.ObjectId `json:"userid"`
	IDManager bson.ObjectId `json:"managerid"`
	Date      int64         `json:"date"`
	// This are the users that comes into the whole process of the product
	UsersRelated []bson.ObjectId `json:"usersrelated"`
	// ServiceCreated is where it was produced
	ServiceCreated bson.ObjectId `json:"servicecreated"`
	IsTrackable    bool          `json:"istrackable"`
}

/***********************  SYNC ***********************/

// Delivery is the amount of items send to a service from the Central, is THE WHOLE DELIVERY with all different stocks
type Delivery struct {
	ID                bson.ObjectId   `json:"id"`
	IDServiceEmitter  bson.ObjectId   `json:"idserviceemitter"`
	IDServiceReceiver bson.ObjectId   `json:"idservicereceiver"`
	Date              int64           `json:"date"`
	IsSent            bool            `json:"issent"`
	Stocks            []bson.ObjectId `json:"stocks"`
	// This parameter will be filled on the way back of the sync, in every local village
	IDManager  bson.ObjectId `json:"idmanager"`
	IsComplete bool          `json:"iscomplete"`
}

// DeliveryType is the close list of different types that an delivery can have
type DeliveryType string

// The options for DeliveryType
const (
	Sent     DeliveryType = "sent"
	Received DeliveryType = "received"
)

/***********************  WORKSHOP ***********************/

// WorkerOrder are the Orders for manufacturing  Products that the services get orders from the central to produce it
type WorkerOrder struct {
	ID          bson.ObjectId `json:"id"`
	IDProduct   bson.ObjectId `json:"idproduct"`
	IDWorker    bson.ObjectId `json:"idworker"`
	IDService   bson.ObjectId `json:"idservice"`
	Date        int64         `json:"date"`
	Quantity    int           `json:"quantity"`
	AlreadyMade int           `json:"alreadyMade"`
	IDManager   bson.ObjectId `json:"idmanager"`
}

// Payment are the payments that the services make to the Users
type Payment struct {
	ID        bson.ObjectId `json:"id"`
	IDWorker  bson.ObjectId `json:"idworker"`
	IDService bson.ObjectId `json:"idservice"`
	Quantity  float64       `json:"quantity"`
	Date      int64         `json:"date"`
	Photo     string        `json:"photo"`
	IDManager bson.ObjectId `json:"idmanager"`
}

// VideoCourse are the different courses in video
type VideoCourse struct {
	ID             bson.ObjectId   `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Price          float64         `json:"price"`
	Photo          string          `json:"photo"`
	ServicesAccess []bson.ObjectId `json:"servicesaccess"`
}

// VideoProblem are the Problems made in the manufacturing of the products
type VideoProblem struct {
	ID            bson.ObjectId `json:"id"`
	IDVideoCourse bson.ObjectId `json:"idvideo"`
	IDWorker      bson.ObjectId `json:"idworker"`
	TextProblem   string        `json:"textproblem"`
	PhotoProblem  string        `json:"photoproblem"`
	Audio         string        `json:"audio"`
	IDManager     bson.ObjectId `json:"idmanager"`
}

// Step are the different steps that we have in the tutorial to make a Product
type Step struct {
	ID              bson.ObjectId   `json:"id"`
	IDVideoCourse   bson.ObjectId   `json:"idvideocourse"`
	IndexOrder      int             `json:"indexorder"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Video           string          `json:"video"`
	Audio           string          `json:"audio"`
	Warnings        []string        `json:"warnings"`
	ToolsNeeded     []bson.ObjectId `json:"toolsneeded"`
	MaterialsNeeded []bson.ObjectId `json:"materialsneeded"`
}

// ServiceOrder are the Orders for production that the Central makes to every Service
type ServiceOrder struct {
	ID            bson.ObjectId    `json:"id"`
	IDVideoCourse bson.ObjectId    `json:"idvideocourse"`
	IDService     bson.ObjectId    `json:"idservice"`
	Quantity      int              `json:"quantity"`
	Assigned      int              `json:"assigned"`
	AlreadyMade   int              `json:"alreadymade"`
	Deadline      int64            `json:"deadline"`
	WindowPeriod  WindowPeriodType `json:"windowperiod"`
	IDManager     bson.ObjectId    `json:"idmanager"`
}

// WindowPeriodType is the close list of different times for putting a Window period for the workers before the deadline for the Workshop
type WindowPeriodType int

// The options for WindowPeriodType
const (
	ThreeDays WindowPeriodType = 3
	OneWeek   WindowPeriodType = 7
	TenDays   WindowPeriodType = 10
	TwoWeeks  WindowPeriodType = 14
	OneMonth  WindowPeriodType = 30
)

// VideoCheckList is the quality control list that the product should pass before being accepted by the manager
type VideoCheckList struct {
	ID            bson.ObjectId `json:"id"`
	IDVideoCourse bson.ObjectId `json:"idvideocourse"`
	Description   string        `json:"description"`
	Audio         string        `json:"audio"`
	Photo         string        `json:"photo"`
}

/***********************  	COMUNICATION    ***********************/

// Conversation is the conversation between two or more persons in the app
type Conversation struct {
	ID    bson.ObjectId   `json:"id"`
	Users []bson.ObjectId `json:"users"`
}

// Message are the messages send in a conversation between two or more persons in the app
type Message struct {
	ID             bson.ObjectId `json:"id"`
	IDUser         bson.ObjectId `json:"iduser"`
	IDConversation bson.ObjectId `json:"idconversation"`
	Text           string        `json:"text"`
	Photo          string        `json:"photo"`
	Audio          string        `json:"audio"`
	Date           int64         `json:"date"`
}

// Report are the reports that the users can make to the App to be read by the central
type Report struct {
	ID        bson.ObjectId `json:"id"`
	IDUser    bson.ObjectId `json:"iduser"`
	Text      string        `json:"text"`
	Photo     string        `json:"photo"`
	Audio     string        `json:"audio"`
	Date      int64         `json:"date"`
	Type      TypeOfReport  `json:"type"`
	IDVillage bson.ObjectId `json:"idvillage"`
	// If Close is true, the report is resolved
	IsClose bool `json:"isclose"`
}

// TypeOfReport is the close list of different types of reports
type TypeOfReport string

// The options for TypeOfReport
const (
	suggestion TypeOfReport = "suggestion"
	lost       TypeOfReport = "lost"
	mistake    TypeOfReport = "mistake"
	abuse      TypeOfReport = "abuse"
	another    TypeOfReport = "another"
)

/***********************  DEVELOPMENT ***********************/

// TimeChanges is the close list of different changes that we can make to the FakeTime in Development
type TimeChanges string

// The options for TimeChanges
const (
	current   TimeChanges = "current"
	reset     TimeChanges = "reset"
	dayPlus   TimeChanges = "dayplus"
	dayMinus  TimeChanges = "dayminus"
	weekPlus  TimeChanges = "weekplus"
	weekMinus TimeChanges = "weekminus"
)

/***********************  TOOLS ASSIGNMENTS ***********************/

// Assignment are tools assigned to the Workers
type Assignment struct {
	ID        bson.ObjectId `json:"id"`
	IDWorker  bson.ObjectId `json:"idworker"`
	IDService bson.ObjectId `json:"idservice"`
	IDManager bson.ObjectId `json:"idmanager"`
	IDStock   bson.ObjectId `json:"idstock"`
	Date      int64         `json:"date"`
	// This will be true when the user returns the Tool to the Manager
	IsBack bool `json:"isback"`
}

/***********************  PDF ***********************/

// PDF is the PDF generated related for print the QRs
type PDF struct {
	ID         bson.ObjectId `json:"id"`
	Date       int64         `json:"date"`
	TypeOfItem TypeOfItem    `json:"typeofitem"`
	Pages      int           `json:"pages"`
	TypeOfPDF  TypeOfPDF     `json:"typeofpdf"`
}

// TypeOfItem is the close list of different items that we can have in the Stock, the 3 main categorization
type TypeOfItem string

// The options for TypeOfItem
const (
	// ServiceProduct is any product that a service produces (Wallet, leather strips, leather bag...)
	ServiceProduct TypeOfItem = "serviceproduct"
	// PrimaryMaterial is the primary row material (wood, screws...)
	PrimaryMaterial TypeOfItem = "primarymaterial"
	// Tool is the tools used in the services
	Tool TypeOfItem = "tool"
	// Task is the Task app that will request to scan the QR code
	Task TypeOfItem = "task"
)

// TypeOfPDF is the close list of different items that we can have in the Stock
type TypeOfPDF string

// The options for TypeOfItem
const (
	PDFQR     TypeOfPDF = "qr"
	PDFReport TypeOfPDF = "report"
)

// QR are the QRs code still unasigned
type QR struct {
	ID         bson.ObjectId `json:"id"`
	PDF        bson.ObjectId `json:"pdf"`
	PageNumber int           `json:"pagenumber"`
	TypeOfItem TypeOfItem    `json:"typeofitem"`
}

// Animal in one individual animal that enter in the slaughterhouse
type Animal struct {
	ID            bson.ObjectId   `json:"id"`
	Front         string          `json:"front"`
	Left          string          `json:"left"`
	Right         string          `json:"right"`
	Back          string          `json:"back"`
	AnimalType    AnimalType      `json:"animaltype"`
	IsClosed      bool            `json:"isclosed"`
	StockObtained []bson.ObjectId `json:"productsobtained"`
}

// AnimalType is the close list of different Animals that we can have in the slaughterhouse
type AnimalType string

// The options for AnimalType
const (
	Goat  AnimalType = "goat"
	Camel AnimalType = "camel"
	Cow   AnimalType = "cow"
)

// HardcodedTypesCategory is the close list of different HardcodedTypesCategory that category can have
type HardcodedTypesCategory string

// The options for HardcodedTypesCategory
const (
	CategoryTypeTools HardcodedTypesCategory = "tools"
	CategoryTypePack  HardcodedTypesCategory = "pack"
	CategoryTypeKit   HardcodedTypesCategory = "kit"
)

// Translation are the changes made in the DataBase as an Event driven, is the backup is case of Sync problem
type Translation struct {
	ID         bson.ObjectId `json:"id"`
	IDInstance bson.ObjectId `json:"idinstance"`
	Instance   interface{}   `json:"instance"`
	Language   Languages     `json:"language"`
}

// Languages is the close list of different Languages that any datagory can have
type Languages string

// The options for Languages
const (
	English Languages = "english"
	Swahili Languages = "swahili"
	German  Languages = "german"
	Spanish Languages = "spanish"
	French  Languages = "french"
)

/****  CHECKING APP  ****/

// ToDo is the item to check
type ToDo struct {
	ID              bson.ObjectId      `json:"id"`
	TitleToDo       string             `json:"titletodo"`
	DescriptionToDo string             `json:"description"`
	IDVillage       bson.ObjectId      `json:"idvillage"`
	IDUser          bson.ObjectId      `json:"iduser"`
	IsTrackable     bool               `json:"istrackable"`
	TimeChecking    TimecheckingType   `json:"timecheking"`
	Descriptions    []FieldDescription `json:"descriptions"`
	Checkboxes      []FieldCheckbox    `json:"checkboxes"`
	Photos          []FieldPhoto       `json:"photos"`
	Numbers         []FieldNumber      `json:"numbers"`
}

// ToDoChecked is the different checks in the time from the ToDo
type ToDoChecked struct {
	ID bson.ObjectId `json:"id"`
	// When was checked
	Date         int64              `json:"date"`
	IDToDo       bson.ObjectId      `json:"idtodo"`
	Descriptions []FieldDescription `json:"descriptions"`
	Checkboxes   []FieldCheckbox    `json:"checkboxes"`
	Photos       []FieldPhoto       `json:"photos"`
	Numbers      []FieldNumber      `json:"numbers"`
}

// FieldDescription are the different descriptions fields that should be checked from the item
type FieldDescription struct {
	ID          bson.ObjectId `json:"id"`
	IDToDo      bson.ObjectId `json:"idtodo"`
	Description string        `json:"description"`
}

// FieldCheckbox are the different checkboxes fields that should be checked from the item
type FieldCheckbox struct {
	ID            bson.ObjectId `json:"id"`
	IDToDo        bson.ObjectId `json:"idtodo"`
	TitleCheckbox string        `json:"titlecheckbox"`
	Checkbox      bool          `json:"checkbox"`
}

// FieldPhoto are the different photo fields that should be checked from the item
type FieldPhoto struct {
	ID         bson.ObjectId `json:"id"`
	IDToDo     bson.ObjectId `json:"idtodo"`
	TitlePhoto string        `json:"titlephoto"`
	Photo      string        `json:"photo"`
}

// FieldNumber are the different number fields that should be checked from the item
type FieldNumber struct {
	ID          bson.ObjectId `json:"id"`
	IDToDo      bson.ObjectId `json:"idtodo"`
	TitleNumber string        `json:"titlenumber"`
	Number      int           `json:"number"`
}

// FieldQR is a QR that need to be scanned by the user
// type FieldQR struct {
// 	ID      bson.ObjectId `json:"id"`
// 	IDToDo  bson.ObjectId `json:"idtodo"`
// 	TitleQR string        `json:"titleqr"`
// 	QRvalue bson.ObjectId `json:"qr"`
// }

// CheckStock are the checked store for the stock in the services
type CheckStock struct {
	ID         bson.ObjectId `json:"id"`
	IDInstance bson.ObjectId `json:"idinstance"`
	IDService  bson.ObjectId `json:"idservice"`
	Date       int64         `json:"date"`
	Photo1     string        `json:"photo1"`
	Photo2     string        `json:"photo2"`
	Photo3     string        `json:"photo3"`
	Photo4     string        `json:"photo4"`
	LastPhoto  int           `json:"lastphoto"`
}

// Anything is to be able to sync the different DataBase
type Anything struct {
	ID       bson.ObjectId `json:"id"`
	Instance interface{}   `json:"interface"`
}
