CREATE TABLE Address
(
    AddressId      int auto_increment primary key,
    Country        varchar(50)  NOT NULL,
    City           varchar(50)  NOT NULL,
    Address        varchar(255) NOT NULL,
    EnglishAddress varchar(255) NULL
)charset = utf8mb4;

CREATE TABLE Box
(
    BoxId   int auto_increment primary key,
    BoxName varchar(30) NOT NULL,
    BoxSize varchar(50) NOT NULL
)charset = utf8mb4,engine = MyISAM;

CREATE TABLE Cargo
(
    CargoId    int auto_increment,
    CargoName  varchar(50) NOT NULL,
    CargoModel varchar(30) NOT NULL,
    CargoNum   integer     NOT NULL,
    CategoryId int  NOT NULL,
    Price      float       NULL,
    Weight     float       NULL,
    Size       varchar(50) NULL,
    PRIMARY KEY (CargoId)
)charset = utf8mb4;

CREATE TABLE Categories
(
    CategoryId   int auto_increment,
    CategoryName varchar(50) NOT NULL,
    PRIMARY KEY (CategoryId)
)charset = utf8mb4,engine = MyISAM;

CREATE TABLE Company
(
    CompanyId   int auto_increment,
    CompanyName varchar(50) NOT NULL,
    CompanyType varchar(30) NOT NULL,
    PRIMARY KEY (CompanyId)
)charset = utf8mb4,engine = MyISAM;


CREATE TABLE Company_box
(
    CompanyId     int Not Null,
    BoxId         int NOT NULL,
    BoxTotalCount integer     NOT NULL,
    BoxUsed       integer     NOT NULL,
    PRIMARY KEY (CompanyId, BoxId)
)charset = utf8mb4;

CREATE TABLE Line
(
    LineId   int auto_increment,
    LineName  varchar(50)  NOT NULL,
    CompanyId int  NOT NULL,
    StartPort varchar(30)  NULL,
    EndPort   varchar(30)  NULL,
    LineDate  varchar(255) NULL,
    PRIMARY KEY (LineId)
)charset = utf8mb4;

Create Table Orders
(
    OrderId                int auto_increment primary key,
    ClientCompanyId        int NOT NULL,
    StartDate              datetime,
    LandTransportCompanyId int,
    SeaTransportCompanyId  int,
    OrderStatus            varchar(30) Not Null
)charset=utf8mb4;

CREATE TABLE Order_Cargo
(
    OrderId int Not Null,
    Cargo   varchar(30) NOT NULL,
    PRIMARY KEY (OrderId, Cargo)
)charset = utf8mb4;

CREATE TABLE Order_staff
(
    OrderId int NOT NULL,
    StaffId int NOT NULL,
    PRIMARY KEY (StaffId, OrderId)
)charset = utf8mb4;

CREATE TABLE Order_Status
(
    OrderId  int not NULL,
    Time     datetime    NULL,
    Status   varchar(50) NULL,
    Position varchar(30) NULL
)charset = utf8mb4;

CREATE TABLE Port
(
    PortId    int auto_increment primary key,
    PortName  varchar(50) NOT NULL,
    PortCode  varchar(30) NOT NULL,
    AddressId int NULL
)charset = utf8mb4,engine = MyISAM;

CREATE TABLE Record
(
    RecordId     int auto_increment primary key,
    RecordChain  varchar(64)  NOT NULL,
    RecordTime   datetime     NOT NULL,
    RecordSite   varchar(255) NOT NULL,
    RecordAction varchar(255) NOT NULL,
    stuffId      int NULL
)charset = utf8mb4;

CREATE TABLE Ship_Line
(
    ShipId     int Not Null,
    LineId     int Not NULL,
    StartDate  datetime    NOT NULL,
    ShipRemain integer     NOT NULL,
    PRIMARY KEY (LineId, ShipId)
)charset = utf8mb4;

CREATE TABLE Ships
(
    ShipId       int auto_increment NOT NULL,
    ShipName     varchar(50) NOT NULL,
    CompanyId    int NOT NULL,
    ShipCapacity float       NOT NULL,
    PRIMARY KEY (ShipId)
)charset = utf8mb4,engine = MyISAM;

CREATE TABLE Staff
(
    StaffId   int auto_increment,
    StaffName varchar(50) NOT NULL,
    StaffJob  varchar(30) NOT NULL,
    CompanyId int NOT NULL,
    isDelete varchar(1) default 'f',
    PRIMARY KEY (StaffId)
)charset = utf8mb4,engine = MyISAM;

CREATE TABLE OrderInfo
(
    OrderId        int  NOT NULL,
    LineId         int  NULL,
    TotalPrice     varchar(10)  NULL,
    StartAddressId int  NULL,
    EndAddressId   int  NULL,
    Phone          varchar(30)  NULL,
    Email          varchar(30)  NULL,
    Fax            varchar(30)  NULL,
    HopeReachDate  date         NULL,
    INCOTERMS      varchar(255) NULL,
    UnStackable    varchar(1)   NULL,
    Perishable     varchar(1)   NULL,
    Dangerous      varchar(1)   NULL,
    Clearance      varchar(1)   NULL,
    Other          varchar(255) NULL,
    CreateDate     datetime default now(),
    PRIMARY KEY (OrderId)
)charset = utf8mb4;

CREATE TABLE Bargain
(
    OrderId   int NOT NULL,
    CompanyId int Not NULL,
    Price     varchar(30) NULL,
    Time      datetime    NULL,
    PRIMARY KEY (OrderId, CompanyId)
)charset = utf8mb4;

CREATE TABLE CompanyInfo
(
    CompanyId int NOT NULL,
    Phone     varchar(30) NULL,
    AddressId int NULL,
    Email     varchar(30) NULL,
    PRIMARY KEY (CompanyId)
)charset = utf8mb4,engine = MyISAM;


CREATE TABLE Relation
(
    CompanyId      int NULL,
    TargetCompanyId int NULL
)charset = utf8mb4,engine = MyISAM;

CREATE TABLE Account
(
    CompanyId int NOT NULL,
    Account   varchar(30) NULL unique,
    PassWord  varchar(64) NULL,
    PRIMARY KEY (CompanyId)
)charset = utf8mb4,engine = MyISAM;


Insert Into Staff Set StaffName = '张三',StaffJob='鸭子',CompanyId = 1;
Insert Into Staff Set StaffName = '李四',StaffJob='鸡',CompanyId = 1;
Insert Into Staff Set StaffName = '王五',StaffJob='鸭子',CompanyId = 1;
Insert Into Staff Set StaffName = '张三22',StaffJob='鸭子',CompanyId = 2;
Insert Into Staff Set StaffName = '李四22',StaffJob='鸡',CompanyId = 2;
Insert Into Staff Set StaffName = '王五22',StaffJob='鸭子',CompanyId = 2;

Insert Into Orders Set ClientCompanyId = 2,StartDate = now(),LandTransportCompanyId = 5,SeaTransportCompanyId = 1,OrderStatus = '订舱';
Insert Into Orders Set ClientCompanyId = 2,StartDate = now(),LandTransportCompanyId = 3,SeaTransportCompanyId = 1,OrderStatus = '订舱';
Insert Into Orders Set ClientCompanyId = 2,StartDate = now(),LandTransportCompanyId = 3,SeaTransportCompanyId = 1,OrderStatus = '订舱';

Insert Into Relation SET CompanyId = 1,TargetCompanyId = 2;
Insert Into Relation SET CompanyId = 3 , TargetCompanyId = 2;
Insert Into Relation SET CompanyId = 2 , TargetCompanyId = 5;

show engines;