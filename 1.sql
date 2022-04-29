
CREATE TABLE Address
(
    AddressId      int auto_increment primary key,
    Country        varchar(50)  NOT NULL,
    City           varchar(50)  NOT NULL,
    Address        varchar(255) NOT NULL
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
    CompanyName varchar(50) NOT NULL default '未命名',
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
    LandTransportCompanyId int Not Null default 0,
    SeaTransportCompanyId  int Not Null default 0,
    OrderStatus            varchar(30) Not Null
)charset = utf8mb4;


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
    Phone     varchar(30) not NULL default '-1',
    AddressId int not NULL default 1,
    Email     varchar(30) not NULL,
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

create table if not exists MessageQueue
(
    MessageId   int auto_increment primary key,
    MessageType int default 0 null,
    FromId      int default 1 not null,
    ToId        int           not null,
    isRead      int default 0 null,
    isDelete    int default 0 null,
    isReply     int default 0 null,
    SendTime    datetime      null
)charset  = utf8mb4,engine = InnoDB;

create table MessageInfo
(
    MessageId      int          not null,
    MessageContent varchar(255) null
)charset  = utf8mb4,engine = InnoDB;



Insert Into Address Set AddressId = 0,Country = '空',City = '空',Address='空';
Insert Into Address Set Country = 'China',City = 'Dalian',Address = 'Dalian HaiShi DaiXue';
Insert Into Company Set CompanyName = 'Manager',CompanyType='管理员';
Insert Into Account Set Account = 'admin',PassWord = 'admin',CompanyId = 1;
Insert Into CompanyInfo Set CompanyId = 1,Phone = '1008611',AddressId = 2,Email = 'admin@qq.com';

Select * From MessageQueue;